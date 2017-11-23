package main

import (
	"flag"
	"fmt"

	"github.com/golang/glog"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"

	captaincyclientset "github.com/guilhem/captaincy/pkg/client/clientset/versioned"

	etcdcluster "github.com/coreos/etcd-operator/pkg/apis/etcd/v1beta2"
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

	list, err := captaincyClient.KinkyV1alpha1().Kinkies(metav1.NamespaceAll).List(metav1.ListOptions{})
	if err != nil {
		glog.Fatalf("Error listing all kinkies: %v", err)
	}

	etcdClient, err := etcdclientset.NewForConfig(cfg)
	if err != nil {
		glog.Fatalf("Error building etcd clientset: %v", err)
	}

	fmt.Printf("cluster %v\n", list)

	for _, cluster := range list.Items {
		fmt.Printf("cluster %s\n", cluster)

		etcdCl := etcdcluster.EtcdCluster{
			Spec: etcdcluster.ClusterSpec{
				Size: 3,
			},
		}
		fmt.Println(cluster.Namespace)
		res, err := etcdClient.EtcdV1beta2().EtcdClusters(cluster.Namespace).Create(&etcdCl)
		if err != nil {
			glog.Errorf("Error spawning ETCD cluster: %v", err)
		}
		glog.Infof("Etcd created %v", res)
	}

}
