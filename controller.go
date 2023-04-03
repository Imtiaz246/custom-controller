package main

import (
	"fmt"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	corev1Informers "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/kubernetes"
	corev1 "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"log"
	"time"
)

type Controller struct {
	kubeClientSet   kubernetes.Interface
	sampleClientSet kubernetes.Interface

	podHasSynced cache.InformerSynced
	Lister       corev1.PodLister
	Queue        workqueue.RateLimitingInterface
}

func (c *Controller) HandleAdd(obj any) {

}
func (c *Controller) HandleDelete(obj any) {

}

func NewController(pod corev1Informers.PodInformer) *Controller {
	c := &Controller{
		podHasSynced: pod.Informer().HasSynced,
		// podHasSynced:
		Lister: pod.Lister(),
		Queue: workqueue.NewNamedRateLimitingQueue(
			workqueue.DefaultControllerRateLimiter(),
			"custom-controller"),
	}
	pod.Informer().AddEventHandler(
		&cache.ResourceEventHandlerFuncs{
			AddFunc:    c.HandleAdd,
			DeleteFunc: c.HandleDelete,
		})
	return c
}

func (c *Controller) runWorker() {
	for c.processNextItem() {

	}
}

func (c *Controller) processNextItem() bool {
	item, quit := c.Queue.Get()
	fmt.Print("item from queue is :", item, quit)
	return false
}

func (c *Controller) Run(threads int, stopChan <-chan struct{}) error {
	defer utilruntime.HandleCrash()
	defer c.Queue.ShutDown()

	log.Println("Waiting for informer caches to sync")
	if ok := cache.WaitForCacheSync(stopChan, c.podHasSynced); !ok {
		return fmt.Errorf("failed to wait for caches for sync")
	}

	for i := 0; i < threads; i++ {
		wait.Until(c.runWorker, time.Second, stopChan)
	}

	<-stopChan
	log.Println("Shutting down custom controller")
	return nil
}
