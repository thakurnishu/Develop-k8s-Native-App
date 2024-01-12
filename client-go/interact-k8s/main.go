package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func KubeconfigHome() string {
	// To get kubeconfig file location
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("ERROR getting UserHome dir: \n%v\n", err.Error())
	}
	kubeconfigPath := filepath.Join(homeDir, ".kube", "config")

	return kubeconfigPath
}

func GetTypedClientSet(kubeconfig *string) *kubernetes.Clientset {
	var config *rest.Config

	if _, err := os.Stat(*kubeconfig); err == nil {
		config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			fmt.Printf("ERROR Building Config from Flag: \n%v\n", err.Error())
		}
	} else {
		// Get config from kubernetes pod using default serviceaccount
		//   which is attached to pod
		config, err = rest.InClusterConfig()
		if err != nil {
			fmt.Printf("ERROR getting Config from K8S Cluster: \n%v\n", err.Error())
		}
	}

	// Typed ClientSet
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Printf("ERROR Creating Typed ClientSet from Config: \n%v\n", err.Error())
	}
	return clientSet
}

func main() {

	kubeconfigHome := KubeconfigHome()
	// File location from user ( --kubeconfig ) or homeDir/.kube/config
	kubeconfig := flag.String("kubeconfig", kubeconfigHome, "Location of KubeConfig file")
	// NameSpace from user (--namespace) or default
	nameSpace := flag.String("namespace", "default", "List of pod in NameSpace")

	// Parse flags once
	flag.Parse()

	clientset := GetTypedClientSet(kubeconfig)
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
