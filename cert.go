package main

import (
	"crypto/x509"
	"fmt"
	"net"

	"github.com/golang/glog"
	"k8s.io/client-go/kubernetes"
	"k8s.io/kubernetes/cmd/kubeadm/app/phases/certs/pkiutil"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	certutil "k8s.io/client-go/util/cert"
	kubeadmconstants "k8s.io/kubernetes/cmd/kubeadm/app/constants"
	certsphase "k8s.io/kubernetes/cmd/kubeadm/app/phases/certs"
)

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
		return fmt.Errorf("failed to list bootstrap tokens [%v]", err)
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
