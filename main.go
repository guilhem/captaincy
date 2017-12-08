package main

import (
	"flag"
	"fmt"
	"net"
	"time"

	"github.com/golang/glog"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/clientcmd"

	apiv1 "k8s.io/api/core/v1"
	extv1beta1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/client-go/kubernetes"

	etcdclientset "github.com/coreos/etcd-operator/pkg/generated/clientset/versioned"
	captaincyclientset "github.com/guilhem/captaincy/pkg/client/clientset/versioned"
	apiextensionsclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"

	"k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm"
	kubeadmconstants "k8s.io/kubernetes/cmd/kubeadm/app/constants"
	nodebootstraptokenphase "k8s.io/kubernetes/cmd/kubeadm/app/phases/bootstraptoken/node"
	"k8s.io/kubernetes/cmd/kubeadm/app/phases/controlplane"

	"k8s.io/kubernetes/cmd/kubeadm/app/util/apiclient"
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

		apiService := &apiv1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "kube-apiserver",
				Namespace: cluster.Namespace,
				Labels: map[string]string{
					"component": "kube-apiserver",
					"tier":      "control-plane",
				},
			},
			Spec: apiv1.ServiceSpec{
				Selector: map[string]string{
					"component": "kube-apiserver",
					"tier":      "control-plane",
				},
				Ports: []apiv1.ServicePort{
					{
						Name:       "https",
						Port:       443,
						TargetPort: intstr.Parse("443"),
						Protocol:   "TCP",
					},
				},
			},
		}

		k8sClient.CoreV1().Services(apiService.Namespace).Create(apiService)

		// Wait for API service to have an IP
		if err := wait.Poll(5*time.Second, 30*time.Minute, func() (bool, error) {
			svc, err := k8sClient.CoreV1().Services(cluster.Namespace).Get("kube-apiserver", metav1.GetOptions{})
			if err != nil {
				return false, err
			}
			if svc.Spec.ClusterIP == "" {
				return false, nil
			}
			return true, nil

		}); err != nil {
			glog.Errorf("error while checking pod status: %v", err)
		}
		svc, err := k8sClient.CoreV1().Services(cluster.Namespace).Get("kube-apiserver", metav1.GetOptions{})
		if err != nil {
			glog.Errorf("Fail service: %v", err)
		}
		internalApiIP := svc.Spec.ClusterIP
		// TODO better way to fix external IP
		externalApiIP := "1.2.3.4"

		kubeadmCfg := &kubeadm.MasterConfiguration{
			Etcd: kubeadm.Etcd{
				Endpoints: []string{fmt.Sprintf("http://%s:%d", etcdCluster.Status.ServiceName, etcdCluster.Status.ClientPort)},
			},
			API: kubeadm.API{
				BindPort: 443,
			},
		}
		if cluster.Spec.Version != "" {
			kubeadmCfg.KubernetesVersion = cluster.Spec.Version
		}

		SetDefaults_MasterConfiguration(kubeadmCfg)

		internalKubeadmCfg := kubeadmCfg.DeepCopy()
		internalKubeadmCfg.API.AdvertiseAddress = internalApiIP

		externalKubeadmCfg := kubeadmCfg.DeepCopy()
		externalKubeadmCfg.API.AdvertiseAddress = externalApiIP

		if err := certsPhase(k8sClient, internalKubeadmCfg, cluster.Namespace, []net.IP{net.ParseIP(internalApiIP), net.ParseIP(externalApiIP)}); err != nil {
			glog.Errorf("Create certificates and configs fail: %v", err)
		}

		semK8sVersion, err := version.ParseSemantic(kubeadmCfg.KubernetesVersion)
		if err != nil {
			glog.Errorf("Fail to parse Version")
		}
		pods := controlplane.GetStaticPodSpecs(externalKubeadmCfg, semK8sVersion)
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
								pod.Spec.Containers[iC].VolumeMounts[iVM].ReadOnly = false
							}
						}
					}
				}
			}
			// add exposed secured port to api-server
			if pod.Name == "kube-apiserver" {
				pod.Spec.Containers[0].Ports = []apiv1.ContainerPort{
					{
						ContainerPort: 443,
						Name:          "secure",
					},
				}
			}
			deploy := &extv1beta1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      pod.Name,
					Namespace: cluster.Namespace,
				},
				Spec: extv1beta1.DeploymentSpec{
					Replicas: int32Ptr(1),
					Template: apiv1.PodTemplateSpec{
						ObjectMeta: pod.ObjectMeta,
						Spec:       pod.Spec,
					},
				},
			}
			if err := apiclient.CreateOrUpdateDeployment(k8sClient, deploy); err != nil {
				glog.Errorf("Pod deployment fail: %v", err)
			}
		}

		tokenDescription := "The default bootstrap token generated."
		if err := nodebootstraptokenphase.UpdateOrCreateToken(k8sClient, kubeadmCfg.Token, false, kubeadmCfg.TokenTTL.Duration, kubeadmconstants.DefaultTokenUsages, []string{kubeadmconstants.V18NodeBootstrapTokenAuthGroup}, tokenDescription); err != nil {
			glog.Errorf("Creation default bootstrap: %v", err)
		}
	}
}

func int32Ptr(i int32) *int32 { return &i }
func int64Ptr(i int64) *int64 { return &i }
