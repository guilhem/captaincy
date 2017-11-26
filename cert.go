package main

import (
	"crypto/rsa"
	"crypto/x509"
	"errors"
	"fmt"
	"net"
	"strconv"

	"github.com/golang/glog"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"

	"k8s.io/kubernetes/cmd/kubeadm/app/phases/certs/pkiutil"
	kubeconfigutil "k8s.io/kubernetes/cmd/kubeadm/app/util/kubeconfig"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubeadmapi "k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm"

	certutil "k8s.io/client-go/util/cert"
	kubeadmconstants "k8s.io/kubernetes/cmd/kubeadm/app/constants"
	certsphase "k8s.io/kubernetes/cmd/kubeadm/app/phases/certs"
)

func test(k8sClient *kubernetes.Clientset, cfg *kubeadmapi.MasterConfiguration, ns string, ips []net.IP) error {
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
		IPs: ips,
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

	config = certutil.Config{
		CommonName:   kubeadmconstants.APIServerKubeletClientCertCommonName,
		Organization: []string{kubeadmconstants.MastersGroup},
		Usages:       []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
	}
	apiClientCert, apiClientKey, err := pkiutil.NewCertAndKey(caCert, caKey, config)
	if err != nil {
		glog.Fatalf("failure while creating API server kubelet client key and certificate: %v", err)
	}

	saSigningKey, err := certutil.NewPrivateKey()
	if err != nil {
		glog.Fatalf("failure while creating service account token signing key: %v", err)
	}

	frontProxyCACert, frontProxyCAKey, err := pkiutil.NewCertificateAuthority()
	if err != nil {
		glog.Fatalf("failure while generating front-proxy CA certificate and key: %v", err)
	}

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

	secret := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: ns,
			Name:      kubeadmconstants.KubeCertificatesVolumeName,
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
		if _, err := k8sClient.CoreV1().Secrets(ns).Update(secret); err != nil {
			return err
		}
	}

	kubeConfigs, err := createKubeConfigFiles(cfg, caCert, caKey)
	if err != nil {
		return err
	}

	schedulerConfig, ok := kubeConfigs[kubeadmconstants.SchedulerKubeConfigFileName]
	if !ok {
		return errors.New("No Scheduler Kubeconfig found")
	}
	controllerConfig, ok := kubeConfigs[kubeadmconstants.ControllerManagerKubeConfigFileName]
	if !ok {
		return errors.New("No Controller Kubeconfig found")
	}

	schedulerFile, err := clientcmd.Write(*schedulerConfig)
	if err != nil {
		return err
	}
	controllerFile, err := clientcmd.Write(*controllerConfig)
	if err != nil {
		return err
	}

	secret = &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: ns,
			Name:      kubeconfigSecret,
		},
		Data: map[string][]byte{
			kubeadmconstants.SchedulerKubeConfigFileName:         schedulerFile,
			kubeadmconstants.ControllerManagerKubeConfigFileName: controllerFile,
		},
	}
	if _, err := k8sClient.CoreV1().Secrets(ns).Create(secret); err != nil {
		if _, err := k8sClient.CoreV1().Secrets(ns).Update(secret); err != nil {
			return err
		}
	}

	return nil
}

func getSecretString(secret *v1.Secret, key string) string {
	if secret.Data == nil {
		return ""
	}
	if val, ok := secret.Data[key]; ok {
		return string(val)
	}
	return ""
}

func createKubeConfigFiles(cfg *kubeadmapi.MasterConfiguration, caCert *x509.Certificate, caKey *rsa.PrivateKey) (map[string]*clientcmdapi.Config, error) {
	configs := make(map[string]*clientcmdapi.Config)
	// gets the KubeConfigSpecs, actualized for the current MasterConfiguration
	specs, err := getKubeConfigSpecs(cfg, caCert, caKey)
	if err != nil {
		return configs, err
	}

	for key, spec := range specs {
		// builds the KubeConfig object
		config, err := buildKubeConfigFromSpec(spec)
		if err != nil {
			return configs, err
		}
		configs[key] = config
	}

	return configs, nil
}

/// Copy of kubeadm code

// clientCertAuth struct holds info required to build a client certificate to provide authentication info in a kubeconfig object
type clientCertAuth struct {
	CAKey         *rsa.PrivateKey
	Organizations []string
}

// tokenAuth struct holds info required to use a token to provide authentication info in a kubeconfig object
type tokenAuth struct {
	Token string
}

// kubeConfigSpec struct holds info required to build a KubeConfig object
type kubeConfigSpec struct {
	CACert         *x509.Certificate
	APIServer      string
	ClientName     string
	TokenAuth      *tokenAuth
	ClientCertAuth *clientCertAuth
}

// buildKubeConfigFromSpec creates a kubeconfig object for the given kubeConfigSpec
func buildKubeConfigFromSpec(spec *kubeConfigSpec) (*clientcmdapi.Config, error) {

	// If this kubeconfig should use token
	if spec.TokenAuth != nil {
		// create a kubeconfig with a token
		return kubeconfigutil.CreateWithToken(
			spec.APIServer,
			"kubernetes",
			spec.ClientName,
			certutil.EncodeCertPEM(spec.CACert),
			spec.TokenAuth.Token,
		), nil
	}

	// otherwise, create a client certs
	clientCertConfig := certutil.Config{
		CommonName:   spec.ClientName,
		Organization: spec.ClientCertAuth.Organizations,
		Usages:       []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
	}
	clientCert, clientKey, err := pkiutil.NewCertAndKey(spec.CACert, spec.ClientCertAuth.CAKey, clientCertConfig)
	if err != nil {
		return nil, fmt.Errorf("failure while creating %s client certificate: %v", spec.ClientName, err)
	}

	// create a kubeconfig with the client certs
	return kubeconfigutil.CreateWithCerts(
		spec.APIServer,
		"kubernetes",
		spec.ClientName,
		certutil.EncodeCertPEM(spec.CACert),
		certutil.EncodePrivateKeyPEM(clientKey),
		certutil.EncodeCertPEM(clientCert),
	), nil
}

// getKubeConfigSpecs returns all KubeConfigSpecs actualized to the context of the current MasterConfiguration
// NB. this methods holds the information about how kubeadm creates kubeconfig files.
func getKubeConfigSpecs(cfg *kubeadmapi.MasterConfiguration, caCert *x509.Certificate, caKey *rsa.PrivateKey) (map[string]*kubeConfigSpec, error) {

	masterEndpoint := "https://" + net.JoinHostPort(cfg.API.AdvertiseAddress, strconv.Itoa(int(cfg.API.BindPort)))

	var kubeConfigSpec = map[string]*kubeConfigSpec{
		kubeadmconstants.AdminKubeConfigFileName: {
			CACert:     caCert,
			APIServer:  masterEndpoint,
			ClientName: "kubernetes-admin",
			ClientCertAuth: &clientCertAuth{
				CAKey:         caKey,
				Organizations: []string{kubeadmconstants.MastersGroup},
			},
		},
		kubeadmconstants.KubeletKubeConfigFileName: {
			CACert:     caCert,
			APIServer:  masterEndpoint,
			ClientName: fmt.Sprintf("system:node:%s", cfg.NodeName),
			ClientCertAuth: &clientCertAuth{
				CAKey:         caKey,
				Organizations: []string{kubeadmconstants.NodesGroup},
			},
		},
		kubeadmconstants.ControllerManagerKubeConfigFileName: {
			CACert:     caCert,
			APIServer:  masterEndpoint,
			ClientName: kubeadmconstants.ControllerManagerUser,
			ClientCertAuth: &clientCertAuth{
				CAKey: caKey,
			},
		},
		kubeadmconstants.SchedulerKubeConfigFileName: {
			CACert:     caCert,
			APIServer:  masterEndpoint,
			ClientName: kubeadmconstants.SchedulerUser,
			ClientCertAuth: &clientCertAuth{
				CAKey: caKey,
			},
		},
	}

	return kubeConfigSpec, nil
}
