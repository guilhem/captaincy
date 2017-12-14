package main

import (
	"fmt"
	"net"
	"time"

	etcdclientset "github.com/coreos/etcd-operator/pkg/generated/clientset/versioned"
	kinky "github.com/guilhem/captaincy/pkg/apis/kinky/v1alpha1"
	apiextensionsclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/kubernetes"
	kubeadmconstants "k8s.io/kubernetes/cmd/kubeadm/app/constants"
	nodebootstraptokenphase "k8s.io/kubernetes/cmd/kubeadm/app/phases/bootstraptoken/node"

	extv1beta1 "k8s.io/api/extensions/v1beta1"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/wait"

	"github.com/golang/glog"
	"k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm"
	"k8s.io/kubernetes/cmd/kubeadm/app/phases/controlplane"
	"k8s.io/kubernetes/cmd/kubeadm/app/util/apiclient"
	"k8s.io/kubernetes/pkg/util/version"
)

func createCluster(k8sClient *kubernetes.Clientset, etcdClient *etcdclientset.Clientset, apiExtClient *apiextensionsclientset.Clientset, cluster kinky.Kinky, baseHost string) error {
	if err := createEtcdOperator(k8sClient, cluster.Namespace); err != nil {
		glog.Errorf("Error spawning ETCD operator: %v", err)
		return err
	}

	etcdName := cluster.Name + "-etcd"

	etcdCluster, err := createEtcdCluster(etcdClient, apiExtClient, etcdName, cluster.Namespace)
	if err != nil {
		glog.Errorf("Error spawning ETCD cluster: %v", err)
		return err
	}

	apiServiceName := cluster.Name + "-kube-apiserver"
	apiService := &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      apiServiceName,
			Namespace: cluster.Namespace,
			Labels: map[string]string{
				"component": apiServiceName,
				"tier":      "control-plane",
			},
		},
		Spec: apiv1.ServiceSpec{
			Selector: map[string]string{
				"component": apiServiceName,
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

	k8sClient.CoreV1().Services(apiService.ObjectMeta.Namespace).Create(apiService)
	// Wait for API service to have an IP
	if err := wait.Poll(5*time.Second, 30*time.Minute, func() (bool, error) {
		svc, err := k8sClient.CoreV1().Services(cluster.Namespace).Get(apiServiceName, metav1.GetOptions{})
		if err != nil {
			return false, err
		}
		if svc.Spec.ClusterIP == "" {
			return false, nil
		}
		return true, nil

	}); err != nil {
		glog.Errorf("error while checking pod status: %v", err)
		return err
	}
	svc, err := k8sClient.CoreV1().Services(cluster.Namespace).Get(apiServiceName, metav1.GetOptions{})
	if err != nil {
		glog.Errorf("Fail service: %v", err)
		return err
	}

	internalAPIIP := svc.Spec.ClusterIP
	// TODO better way to fix external IP
	externalAPIIP := "1.2.3.4"

	kubeadmCfg := &kubeadm.MasterConfiguration{
		Etcd: kubeadm.Etcd{
			Endpoints: []string{fmt.Sprintf("http://%s:%d", etcdCluster.Status.ServiceName, etcdCluster.Status.ClientPort)},
		},
		API: kubeadm.API{
			BindPort:         443,
			AdvertiseAddress: "0.0.0.0",
		},
		ControllerManagerExtraArgs: map[string]string{"address": "0.0.0.0"},
		SchedulerExtraArgs:         map[string]string{"address": "0.0.0.0"},
	}
	if cluster.Spec.Version != "" {
		kubeadmCfg.KubernetesVersion = cluster.Spec.Version
	}

	SetDefaults_MasterConfiguration(kubeadmCfg)

	internalKubeadmCfg := kubeadmCfg.DeepCopy()
	internalKubeadmCfg.API.AdvertiseAddress = internalAPIIP

	externalKubeadmCfg := kubeadmCfg.DeepCopy()
	externalKubeadmCfg.API.AdvertiseAddress = externalAPIIP

	clusterHostname := fmt.Sprintf("%s.%s.%s", cluster.Name, cluster.Namespace, baseHost)

	if err := certsPhase(k8sClient, internalKubeadmCfg, cluster.Namespace, []net.IP{net.ParseIP(internalAPIIP)}, clusterHostname); err != nil {
		glog.Errorf("Create certificates and configs fail: %v", err)
		return err
	}

	semK8sVersion, err := version.ParseSemantic(kubeadmCfg.KubernetesVersion)
	if err != nil {
		glog.Errorf("Fail to parse Version")
		return err
	}
	pods := controlplane.GetStaticPodSpecs(externalKubeadmCfg, semK8sVersion)
	for _, pod := range pods {
		// add exposed secured port to api-server
		if pod.Name == "kube-apiserver" {
			pod.Spec.Containers[0].Ports = []apiv1.ContainerPort{
				{
					ContainerPort: 443,
					Name:          "secure",
				},
			}
		}
		// We don't want to use host network
		pod.Spec.HostNetwork = false
		// patch pod names
		pod.ObjectMeta.Name = cluster.Name + "-" + pod.ObjectMeta.Name
		pod.ObjectMeta.Labels["component"] = cluster.Name + "-" + pod.ObjectMeta.Labels["component"]
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

	if err := createIngress(k8sClient, "ingress-"+cluster.Name, cluster.Namespace, clusterHostname, apiServiceName); err != nil {
		glog.Errorf("could not create ingress: %v", err)
		return err
	}
	tokenDescription := "The default bootstrap token generated."
	if err := nodebootstraptokenphase.UpdateOrCreateToken(k8sClient, kubeadmCfg.Token, false, kubeadmCfg.TokenTTL.Duration, kubeadmconstants.DefaultTokenUsages, []string{kubeadmconstants.V18NodeBootstrapTokenAuthGroup}, tokenDescription); err != nil {
		glog.Errorf("Creation default bootstrap: %v", err)
		return err
	}

	return nil
}
