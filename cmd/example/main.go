package main

import (
	"flag"
	"fmt"

	"github.com/golang/glog"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"

	captaincyclientset "github.com/guilhem/captaincy/pkg/client/clientset/versioned"

	etcdclientset "github.com/coreos/etcd-operator/pkg/generated/clientset/versioned"
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

	list, err := captaincyClient.KinkyV1alpha1().Kinkies().List(metav1.ListOptions{})
	if err != nil {
		glog.Fatalf("Error listing all kinkies: %v", err)
	}

	for _, cluster := range list.Items {
		fmt.Printf("cluster %s\n", cluster)
	}

	etcdClient, err := etcdclientset.NewForConfig(cfg)
	if err != nil {
		glog.Fatalf("Error building captaincy clientset: %v", err)
	}
	res, err := etcdClient.EtcdV1beta2().EtcdClusters("lol").Create()
	if err != nil {
		glog.Fatalf("Error building captaincy clientset: %v", err)
	}
	glog.Fatalf("Error building captaincy clientset: %v", res)
}
