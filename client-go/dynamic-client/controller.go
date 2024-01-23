package main

import (
	"log"
	"time"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/tools/cache"
)

type controller struct {
	dynamicClient dynamic.Interface
	informer      cache.SharedIndexInformer
}

func newController(dyn dynamic.Interface, dynInformer dynamicinformer.DynamicSharedInformerFactory) *controller {

	sharedInformer := dynInformer.ForResource(schema.GroupVersionResource{
		Group:    "viveksingh.dev",
		Version:  "v1alpha1",
		Resource: "klusters",
	}).Informer()

	sharedInformer.AddEventHandler(
		cache.ResourceEventHandlerDetailedFuncs{
			AddFunc: func(obj interface{}, isInInitialList bool) {
				log.Println("resource was added")
			},
		},
	)

	return &controller{
		dynamicClient: dyn,
		informer:      sharedInformer,
	}
}

func (c *controller) run(ch <-chan struct{}) {

	log.Println("starting controller...")
	if !cache.WaitForCacheSync(ch, c.informer.HasSynced) {
		log.Println("waiting for cache to be synced")
	}

	go wait.Until(c.worker, 1*time.Second, ch)

	<-ch

}

func (c *controller) worker() {
	for c.processItem() {

	}
}

func (c *controller) processItem() bool {

	return true
}
