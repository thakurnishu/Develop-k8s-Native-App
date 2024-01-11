package main

import (
	"context"
	"flag"
	"fmt"

	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func GetClientSet() *kubernetes.Clientset {
	// Get config from kubernetes pod using default serviceaccount
	//   which is attached to pod
	config, err := rest.InClusterConfig()
	if err != nil {
		fmt.Printf("ERROR getting Config from K8S Cluster: \n%v\n", err.Error())
	}

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Printf("ERROR Creating ClientSet from Config: \n%v\n", err.Error())
	}
	return clientSet
}

func main() {

	// kubeconce from user (--namespace) or default
	nameSpace := flag.String("namespace", "default", "List of pod in NameSpace")

	// Parse flags once
	flag.Parse()

	clientset := GetClientSet()
	ctx := context.Background()

	podList, err := clientset.CoreV1().Pods(*nameSpace).List(ctx, metaV1.ListOptions{})
	if err != nil {
		fmt.Printf("ERROR in listing Pod: \n%v\n", err.Error())
	}
	fmt.Printf("Pod in %v namespace\n", *nameSpace)
	for _, pod := range podList.Items {
		fmt.Println(pod.Name)
	}

	fmt.Println()

	deploymentList, err := clientset.AppsV1().Deployments(*nameSpace).List(ctx, metaV1.ListOptions{})
	if err != nil {
		fmt.Printf("ERROR in listing Deployments: \n%v\n", err.Error())
	}
	fmt.Printf("Deployments in %v namespace\n", *nameSpace)
	for _, deployment := range deploymentList.Items {
		fmt.Println(deployment.Name)
	}
}
