package main

import (
	"github.com/golang/glog"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"

	etcdv1beta2 "github.com/coreos/etcd-operator/pkg/apis/etcd/v1beta2"
	appsv1beta1 "k8s.io/api/apps/v1beta1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	etcdclientset "github.com/coreos/etcd-operator/pkg/generated/clientset/versioned"
	apiextensionsclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
)

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

func createEtcdCluster(client *etcdclientset.Clientset, apiExtClient *apiextensionsclientset.Clientset, name string, ns string) (*etcdv1beta2.EtcdCluster, error) {
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

	if etcdCluster, err := client.EtcdV1beta2().EtcdClusters(ns).Get(name, metav1.GetOptions{}); err == nil {
		return etcdCluster, nil
	}

	etcdCl := etcdv1beta2.EtcdCluster{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			Labels: map[string]string{
				"captaincy": "kinky",
			},
		},
		Spec: etcdv1beta2.ClusterSpec{
			Size: 1,
		},
	}

	client.EtcdV1beta2().EtcdClusters(ns).Create(&etcdCl)
	return client.EtcdV1beta2().EtcdClusters(ns).Get(name, metav1.GetOptions{})
}
