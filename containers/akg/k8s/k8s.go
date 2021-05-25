package k8s

import (
	"context"
	"fmt"
	"log"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type K8s struct {
	config    *rest.Config
	client    *kubernetes.Clientset
	InCluster bool
}

type Instance struct {
	Name string
}

type App struct {
	Name      string
	Instances []Instance
}

func (k *K8s) Configure() {
	// TODO
	// need to use some digital ocean config here :/
	// https://stackoverflow.com/questions/65042279/python-kubernetes-client-requests-fail-with-unable-to-get-local-issuer-certific
	if k.InCluster {
		config, err := rest.InClusterConfig()
		if err != nil {
			panic(err.Error())
		}

		k.config = config
	} else {
		homedir := homedir.HomeDir()
		config, err := clientcmd.BuildConfigFromFlags("", fmt.Sprintf("%s/%s", homedir, ".kube/config"))
		if err != nil {
			panic(err.Error())
		}
		k.config = config
	}
}

func (k *K8s) Connect() {
	clientset, err := kubernetes.NewForConfig(k.config)
	if err != nil {
		panic(err.Error())
	}

	k.client = clientset
}

func (k *K8s) Instances() ([]Instance, error) {
	pods, err := k.client.CoreV1().Pods("app").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Printf("failed to list pods: %s", err)
		return []Instance{}, err
	}

	instances := []Instance{}

	for _, pod := range pods.Items {
		instance := Instance{
			Name: pod.Name,
		}
		instances = append(instances, instance)
	}

	return instances, nil
}

func (k *K8s) Apps() ([]App, error) {
	deployments, err := k.client.AppsV1().Deployments("app").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Printf("failed to list deployments: %s", err)
		return []App{}, err
	}

	apps := []App{}

	for _, deployment := range deployments.Items {
		name := deployment.Name
		app := App{
			Name: name,
		}

		pods, err := k.client.CoreV1().Pods("app").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			log.Printf("failed to list pods: %s", err)
			return []App{}, err
		}

		for _, pod := range pods.Items {
			if strings.Contains(pod.Name, fmt.Sprintf("%s-", name)) {
				instance := Instance{
					Name: pod.Name,
				}
				app.Instances = append(app.Instances, instance)
			}
		}

		apps = append(apps, app)
	}

	return apps, nil
}
