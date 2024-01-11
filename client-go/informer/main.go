package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	// metav1 "k8s.io/apimachinery/pkg/apis/meta/v1" for List Options
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
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

func GetClientSet(kubeconfig *string) *kubernetes.Clientset {
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

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Printf("ERROR Creating ClientSet from Config: \n%v\n", err.Error())
	}
	return clientSet
}

func main() {

	kubeconfigHome := KubeconfigHome()
	// File location from user ( --kubeconfig ) or homeDir/.kube/config
	kubeconfig := flag.String("kubeconfig", kubeconfigHome, "Location of KubeConfig file")
	namespace := flag.String("namespace", "default", "NameSpace of resource")
	// Parse flags once
	flag.Parse()
	clientset := GetClientSet(kubeconfig)

	informerfactory := informers.NewSharedInformerFactory(clientset, 30*time.Second) // defaultResync is about 10-20 mins
	/*
		informers.NewFilteredSharedInformerFactory(clientset, 10*time.Minute, "default", func(lo *metav1.ListOptions) {
			lo.LabelSelector = ""
			lo.APIVersion = ""
		})
	*/

	podinformer := informerfactory.Core().V1().Pods()
	podinformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(new interface{}) {
			// fmt.Printf("%v\n WAS ADDED", new)
			fmt.Println("Add was Called")
		},
		DeleteFunc: func(del interface{}) {
			// fmt.Printf("%v\n WAS DELETED", del)
			fmt.Println("Delete was Called")
		},
		UpdateFunc: func(old, new interface{}) {
			// fmt.Printf("%v\n WAS UPDATED BY\n %v", old, new)
			fmt.Println("Update was Called")
		},
	})

	informerfactory.Start(wait.NeverStop)
	informerfactory.WaitForCacheSync(wait.NeverStop)
	pod, err := podinformer.Lister().Pods(*namespace).Get("kube-scheduler-kind-control-plane")
	if err != nil {
		fmt.Printf("ERROR Listing podinformer: \n%v\n", err.Error())
	}
	fmt.Println(pod)

}
