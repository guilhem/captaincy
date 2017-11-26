package main

import (
	"flag"
	"fmt"
	"net"

	"github.com/golang/glog"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/tools/clientcmd"

	appsv1beta1 "k8s.io/api/apps/v1beta1"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"

	apiextensionsclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"

	captaincyclientset "github.com/guilhem/captaincy/pkg/client/clientset/versioned"

	etcdclientset "github.com/coreos/etcd-operator/pkg/generated/clientset/versioned"

	"k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm"
	kubeadmconstants "k8s.io/kubernetes/cmd/kubeadm/app/constants"
	"k8s.io/kubernetes/cmd/kubeadm/app/phases/controlplane"

	"k8s.io/kubernetes/pkg/util/version"
)

var (
	kuberconfig = flag.String("kubeconfig", "", "Path to a kubeconfig. Only required if out-of-cluster.")
	master      = flag.String("master", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")
)

const (
	kubeconfigSecret = "kubeconfig"
)

func main() {
	flag.Parse()

	cfg, err := clientcmd.BuildConfigFromFlags(*master, *kuberconfig)
	if err != nil {
		glog.Fatalf("Error building kubeconfig: %v", err)
	}

	captaincyClient, err := captaincyclientset.NewForConfig(cfg)
	if err != nil {
		glog.Fatalf("Error building captaincy clientset: %v", err)
	}

	etcdClient, err := etcdclientset.NewForConfig(cfg)
	if err != nil {
		glog.Fatalf("Error building etcd clientset: %v", err)
	}

	k8sClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		glog.Fatalf("Error building etcd clientset: %v", err)
	}

	apiExtClient, err := apiextensionsclientset.NewForConfig(cfg)
	if err != nil {
		glog.Fatalf("Error building etcd clientset: %v", err)
	}

	list, err := captaincyClient.KinkyV1alpha1().Kinkies(metav1.NamespaceAll).List(metav1.ListOptions{})
	if err != nil {
		glog.Fatalf("Error listing all kinkies: %v", err)
	}

	for _, cluster := range list.Items {

		if err := createEtcdOperator(k8sClient, cluster.Namespace); err != nil {
			glog.Errorf("Error spawning ETCD operator: %v", err)
		}

		etcdName := "etcd-" + cluster.Name

		etcdCluster, err := createEtcdCluster(etcdClient, apiExtClient, etcdName, cluster.Namespace)

		if err != nil {
			glog.Errorf("Error spawning ETCD cluster: %v", err)
		}

		glog.Infof("Etcd created")

		apiService := &apiv1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name: "kube-apiserver",
				Labels: map[string]string{
					"component": "kube-apiserver",
					"tier":      "control-plane",
				},
			},
			Spec: apiv1.ServiceSpec{
				Ports: []apiv1.ServicePort{
					{
						Name:       "https",
						Port:       443,
						TargetPort: intstr.Parse("secure"),
						Protocol:   "TCP",
					},
				},
			},
		}

		if _, err := k8sClient.CoreV1().Services(cluster.Namespace).Create(apiService); err != nil {
			glog.Errorf("Fail service: %v", err)
		}

		svc, err := k8sClient.CoreV1().Services(cluster.Namespace).Get("kube-apiserver", metav1.GetOptions{})
		if err != nil {
			glog.Errorf("Fail service: %v", err)
		}
		apiIP := svc.Spec.ClusterIP

		kubeadmCfg := &kubeadm.MasterConfiguration{
			Etcd: kubeadm.Etcd{
				Endpoints: []string{fmt.Sprintf("http://%s:%d", etcdCluster.Status.ServiceName, etcdCluster.Status.ClientPort)},
			},
			API: kubeadm.API{
				AdvertiseAddress: apiIP,
				BindPort:         443,
			},
		}
		if cluster.Spec.Version != "" {
			kubeadmCfg.KubernetesVersion = cluster.Spec.Version
		}

		SetDefaults_MasterConfiguration(kubeadmCfg)

		if err := test(k8sClient, kubeadmCfg, cluster.Namespace, []net.IP{}); err != nil {
			glog.Errorf("Create certificates and configs fail: %v", err)
		}

		semK8sVersion, err := version.ParseSemantic(kubeadmCfg.KubernetesVersion)
		if err != nil {
			glog.Errorf("Fail to parse Version")
		}
		pods := controlplane.GetStaticPodSpecs(kubeadmCfg, semK8sVersion)
		for _, pod := range pods {
			// We don't want to use host network
			pod.Spec.HostNetwork = false
			// Use secret instead of hostPath
			for i, volume := range pod.Spec.Volumes {
				if volume.Name == kubeadmconstants.KubeCertificatesVolumeName {
					pod.Spec.Volumes[i].VolumeSource = apiv1.VolumeSource{
						Secret: &apiv1.SecretVolumeSource{
							SecretName: kubeadmconstants.KubeCertificatesVolumeName,
						},
					}
				}
				if volume.Name == kubeadmconstants.KubeConfigVolumeName {
					pod.Spec.Volumes[i].VolumeSource = apiv1.VolumeSource{
						Secret: &apiv1.SecretVolumeSource{
							SecretName: kubeconfigSecret,
						},
					}
					for iC, container := range pod.Spec.Containers {
						for iVM, volumeMount := range container.VolumeMounts {
							if volumeMount.Name == kubeadmconstants.KubeConfigVolumeName {
								pod.Spec.Containers[iC].VolumeMounts[iVM].MountPath = kubeadmconstants.KubernetesDir
							}
						}
					}
				}
			}
			// add exposed secured port to api-server
			if pod.Name == "kube-apiserver" {
				pod.Spec.Containers[0].Ports = []apiv1.ContainerPort{
					{
						ContainerPort: 6443,
						Name:          "secure",
					},
				}
			}
			deploy := &appsv1beta1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name: pod.Name,
				},
				Spec: appsv1beta1.DeploymentSpec{
					Replicas: int32Ptr(1),
					Template: apiv1.PodTemplateSpec{
						ObjectMeta: pod.ObjectMeta,
						Spec:       pod.Spec,
					},
				},
			}
			if _, err := k8sClient.AppsV1beta1().Deployments(cluster.Namespace).Create(deploy); err != nil {
				if _, err := k8sClient.AppsV1beta1().Deployments(cluster.Namespace).Update(deploy); err != nil {
					glog.Errorf("Pod deployment fail: %v", err)
				}
			}
		}
	}
}

func int32Ptr(i int32) *int32 { return &i }
func int64Ptr(i int64) *int64 { return &i }
