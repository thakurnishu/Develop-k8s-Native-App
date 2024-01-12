package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/viveksinghggits/kluster/pkg/apis/viveksingh.dev/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
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

func GetDynamicClientSet(kubeconfig *string) *dynamic.DynamicClient {
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
	/*
		// Typed ClientSet
		typedClientSet, err := kubernetes.NewForConfig(config)
		if err != nil {
			fmt.Printf("ERROR Creating Typed ClientSet from Config: \n%v\n", err.Error())
		}
	*/

	// Dynamic ClientSet
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		fmt.Printf("ERROR Creating Dynamic Client from Config: \n%v\n", err.Error())
	}

	return dynamicClient
}

func main() {

	kubeconfigHome := KubeconfigHome()
	// File location from user ( --kubeconfig ) or homeDir/.kube/config
	kubeconfig := flag.String("kubeconfig", kubeconfigHome, "Location of KubeConfig file")
	// NameSpace from user (--namespace) or default
	nameSpace := flag.String("namespace", "default", "List of pod in NameSpace")

	// Parse flags once
	flag.Parse()
	dynamicClient := GetDynamicClientSet(kubeconfig)

	ctx := context.Background()

	// List
	listResources, err := dynamicClient.Resource(schema.GroupVersionResource{
		Group:    "viveksingh.dev",
		Version:  "v1alpha1",
		Resource: "klusters",
	}).Namespace(*nameSpace).List(ctx, metav1.ListOptions{})
	if err != nil {
		fmt.Printf("ERROR Listing Resources from Dynamic Client: \n%v\n", err.Error())
	}
	fmt.Printf("Length of Items from .List : %v\n", len(listResources.Items))

	// Get
	unstructuredObject, err := dynamicClient.Resource(schema.GroupVersionResource{
		Group:    "viveksingh.dev",
		Version:  "v1alpha1",
		Resource: "klusters",
	}).Namespace(*nameSpace).Get(ctx, "kluster-1", metav1.GetOptions{})
	if err != nil {
		fmt.Printf("ERROR Getting Resources from Dynamic Client: \n%v\n", err.Error())
	}

	// Imported Object Structure
	structuedTypeOfUnStructureedObject := v1alpha1.Kluster{}

	// getting and setting fields on unstructredObject
	fmt.Printf("Got the object %s\n", unstructuredObject.GetName())

	// how to convert unstructredObject into a typed Object
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(unstructuredObject.UnstructuredContent(), &structuedTypeOfUnStructureedObject)
	if err != nil {
		fmt.Printf("ERROR While Converting UnStructObjectType to StructObjectType: \n%v\n", err.Error())
	}

	fmt.Printf("The Concrete type : \n%+v", structuedTypeOfUnStructureedObject)

}
