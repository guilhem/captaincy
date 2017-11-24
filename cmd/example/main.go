package main

import (
	"crypto/x509"
	"flag"
	"fmt"
	"net"

	"github.com/golang/glog"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/clientcmd"

	appsv1beta1 "k8s.io/api/apps/v1beta1"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"

	apiextensionsclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"

	captaincyclientset "github.com/guilhem/captaincy/pkg/client/clientset/versioned"

	etcdcluster "github.com/coreos/etcd-operator/pkg/apis/etcd/v1beta2"
	etcdclientset "github.com/coreos/etcd-operator/pkg/generated/clientset/versioned"

	certutil "k8s.io/client-go/util/cert"
	"k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm"
	kubeadmconstants "k8s.io/kubernetes/cmd/kubeadm/app/constants"
	certsphase "k8s.io/kubernetes/cmd/kubeadm/app/phases/certs"
	"k8s.io/kubernetes/cmd/kubeadm/app/phases/certs/pkiutil"
	"k8s.io/kubernetes/cmd/kubeadm/app/phases/controlplane"

	"k8s.io/kubernetes/pkg/util/version"
)

var (
	kuberconfig = flag.String("kubeconfig", "", "Path to a kubeconfig. Only required if out-of-cluster.")
	master      = flag.String("master", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")
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

		if err := createEtcdCluster(etcdClient, apiExtClient, etcdName, cluster.Namespace); err != nil {
			glog.Errorf("Error spawning ETCD cluster: %v", err)
		}
		glog.Infof("Etcd created")
		test(k8sClient, cluster.Namespace)

		kubeadmCfg := &kubeadm.MasterConfiguration{
			Etcd: kubeadm.Etcd{
				Endpoints: []string{"http://" + etcdName + ":2"},
			},
			Networking: kubeadm.Networking{
				ServiceSubnet: "10.1.0.0/10",
				PodSubnet:     "192.168.0.0/16",
			},
		}

		k8sVersion, err := version.ParseSemantic(cluster.Spec.Version)
		if err != nil {
			glog.Errorf("Fail to parse Version")
		}
		pods := controlplane.GetStaticPodSpecs(kubeadmCfg, k8sVersion)
		for _, pod := range pods {
			deploy := &appsv1beta1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name: "",
				},
			}
			k8sClient.AppsV1beta1().Deployments(cluster.Namespace).Create(deploy)
		}
		fmt.Printf("pods %#v", pods)
	}

}

func test(k8sClient *kubernetes.Clientset, ns string) error {
	caCert, caKey, _ := certsphase.NewCACertAndKey()
	// fmt.Printf("ca: %v - %v\n", caCert, caKey)

	altNames := &certutil.AltNames{
		DNSNames: []string{
			"Default",
			"kubernetes",
			"kubernetes.default",
			"kubernetes.default.svc",
			fmt.Sprintf("kubernetes.default.svc.%s", "apiserver"),
		},
		IPs: []net.IP{
			[]byte{10, 0, 0, 1},
			[]byte{10, 0, 0, 2},
		},
	}
	config := certutil.Config{
		CommonName: kubeadmconstants.APIServerCertCommonName,
		AltNames:   *altNames,
		Usages:     []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
	}

	apiCert, apiKey, err := pkiutil.NewCertAndKey(caCert, caKey, config)
	if err != nil {
		glog.Fatalf("failure while creating API server key and certificate: %v", err)
	}

	// fmt.Printf("\napicert: %v, %v\n", apiKey, apiCert)

	config = certutil.Config{
		CommonName:   kubeadmconstants.APIServerKubeletClientCertCommonName,
		Organization: []string{kubeadmconstants.MastersGroup},
		Usages:       []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
	}
	apiClientCert, apiClientKey, err := pkiutil.NewCertAndKey(caCert, caKey, config)
	if err != nil {
		glog.Fatalf("failure while creating API server kubelet client key and certificate: %v", err)
	}

	// fmt.Printf("\napicliencert: %v, %v\n", apiClientCert, apiClientKey)

	saSigningKey, err := certutil.NewPrivateKey()
	if err != nil {
		glog.Fatalf("failure while creating service account token signing key: %v", err)
	}
	fmt.Printf("\nsaSigningKey: %v\n", saSigningKey)

	frontProxyCACert, frontProxyCAKey, err := pkiutil.NewCertificateAuthority()
	if err != nil {
		glog.Fatalf("failure while generating front-proxy CA certificate and key: %v", err)
	}
	// fmt.Printf("\nfrontProxyCACert: %v, %v\n", frontProxyCACert, frontProxyCAKey)

	config = certutil.Config{
		CommonName: kubeadmconstants.FrontProxyClientCertCommonName,
		Usages:     []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
	}
	frontProxyClientCert, frontProxyClientKey, err := pkiutil.NewCertAndKey(frontProxyCACert, frontProxyCAKey, config)
	if err != nil {
		glog.Fatalf("failure while creating front-proxy client key and certificate: %v", err)
	}
	// fmt.Printf("\nfrontProxyClientCert: %v, %v\n", frontProxyClientCert, frontProxyClientKey)
	// // PHASE 1: Generate certificates
	// if err := certsphase.CreatePKIAssets(i.cfg); err != nil {
	// 	return err
	// }
	//
	// // PHASE 2: Generate kubeconfig files for the admin and the kubelet
	// if err := kubeconfigphase.CreateInitKubeConfigFiles(kubeConfigDir, i.cfg); err != nil {
	// 	return err
	// }

	pub, err := certutil.EncodePublicKeyPEM(&saSigningKey.PublicKey)
	if err != nil {
		glog.Fatalf("failure while creating public key: %v", err)
	}

	secret := &apiv1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: ns,
			Name:      "pki-k8s",
		},
		Data: map[string][]byte{
			kubeadmconstants.CACertName:                     certutil.EncodeCertPEM(caCert),
			kubeadmconstants.CAKeyName:                      certutil.EncodePrivateKeyPEM(caKey),
			kubeadmconstants.APIServerCertName:              certutil.EncodeCertPEM(apiCert),
			kubeadmconstants.APIServerKeyName:               certutil.EncodePrivateKeyPEM(apiKey),
			kubeadmconstants.APIServerKubeletClientCertName: certutil.EncodeCertPEM(apiClientCert),
			kubeadmconstants.APIServerKubeletClientKeyName:  certutil.EncodePrivateKeyPEM(apiClientKey),
			kubeadmconstants.ServiceAccountPublicKeyName:    pub,
			kubeadmconstants.ServiceAccountPrivateKeyName:   certutil.EncodePrivateKeyPEM(saSigningKey),
			kubeadmconstants.FrontProxyCAKeyName:            certutil.EncodePrivateKeyPEM(frontProxyCAKey),
			kubeadmconstants.FrontProxyCACertName:           certutil.EncodeCertPEM(frontProxyCACert),
			kubeadmconstants.FrontProxyClientKeyName:        certutil.EncodePrivateKeyPEM(frontProxyClientKey),
			kubeadmconstants.FrontProxyClientCertName:       certutil.EncodeCertPEM(frontProxyClientCert),
		},
	}

	if _, err := k8sClient.CoreV1().Secrets(ns).Create(secret); err != nil {
		return fmt.Errorf("failed to list bootstrap tokens [%v]", err)
	}

	return nil
}

func getSecretString(secret *apiv1.Secret, key string) string {
	if secret.Data == nil {
		return ""
	}
	if val, ok := secret.Data[key]; ok {
		return string(val)
	}
	return ""
}

func createEtcdOperator(client *kubernetes.Clientset, ns string) error {
	if _, err := client.AppsV1beta1().Deployments(ns).Get("etcd-operator", metav1.GetOptions{}); err == nil {
		return nil
	}

	deployment := &appsv1beta1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "etcd-operator",
		},
		Spec: appsv1beta1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"name": "etcd-operator",
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:    "etcd-operator",
							Image:   "quay.io/coreos/etcd-operator:v0.7.0",
							Command: []string{"etcd-operator"},
							Env: []apiv1.EnvVar{
								{
									Name: "MY_POD_NAMESPACE",
									ValueFrom: &apiv1.EnvVarSource{
										FieldRef: &apiv1.ObjectFieldSelector{
											FieldPath: "metadata.namespace",
										},
									},
								},
								{
									Name: "MY_POD_NAME",
									ValueFrom: &apiv1.EnvVarSource{
										FieldRef: &apiv1.ObjectFieldSelector{
											FieldPath: "metadata.name",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	_, err := client.AppsV1beta1().Deployments(ns).Create(deployment)
	return err
}

func createEtcdCluster(client *etcdclientset.Clientset, apiExtClient *apiextensionsclientset.Clientset, name string, ns string) error {
	if _, err := apiExtClient.ApiextensionsV1beta1().CustomResourceDefinitions().Get("etcdclusters.etcd.database.coreos.com", metav1.GetOptions{}); err != nil {

		wi, err := apiExtClient.ApiextensionsV1beta1().CustomResourceDefinitions().Watch(metav1.ListOptions{
			TimeoutSeconds: int64Ptr(30),
			FieldSelector:  fields.OneTermEqualSelector("metadata.name", "etcdclusters.etcd.database.coreos.com").String(),
		})
		if err != nil {
			glog.Errorf("Error spawning ETCD cluster: %v", err)
		}
		defer wi.Stop()

		select {
		case watchEvent := <-wi.ResultChan():
			if watch.Added == watchEvent.Type {
				glog.Info("etcd operator register")
				wi.Stop()
			} else {
				glog.Errorf("expected add, but got %#v", watchEvent)
			}
		}
	} else {
		glog.Info("etcdclusters.etcd.database.coreos.com exist")
	}

	if _, err := client.EtcdV1beta2().EtcdClusters(ns).Get(name, metav1.GetOptions{}); err == nil {
		return nil
	}

	etcdCl := etcdcluster.EtcdCluster{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			Labels: map[string]string{
				"captaincy": "kinky",
			},
		},
		Spec: etcdcluster.ClusterSpec{
			Size: 3,
		},
	}

	_, err := client.EtcdV1beta2().EtcdClusters(ns).Create(&etcdCl)
	return err
}

func int32Ptr(i int32) *int32 { return &i }
func int64Ptr(i int64) *int64 { return &i }
