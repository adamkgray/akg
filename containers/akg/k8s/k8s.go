package k8s

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type K8s struct {
	config    *rest.Config
	client    *kubernetes.Clientset
	InCluster bool
}

func (k *K8s) Configure() {
	if k.InCluster {
		config, err := rest.InClusterConfig()
		if err != nil {
			panic(err.Error())
		}

		k.config = config
	} else {
		panic("not in cluster!")
	}
}

func (k *K8s) Connect() {
	clientset, err := kubernetes.NewForConfig(k.config)
	if err != nil {
		panic(err.Error())
	}

	k.client = clientset
}

func (k *K8s) Apps() []string {
	deployments, err := k.client.AppsV1().Deployments("app").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	var AppNames []string

	for _, deployment := range deployments.Items {
		name := deployment.Name
		AppNames = append(AppNames, name)
	}

	return AppNames
}
