package controller

import (
	"fmt"
	customclientset "github.com/imtiaz246/custom-controller/pkg/client/clientset/versioned"
	customv1beta1informers "github.com/imtiaz246/custom-controller/pkg/client/informers/externalversions/cho.me/v1beta1"
	"github.com/imtiaz246/custom-controller/pkg/client/listers/cho.me/v1beta1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	kubeappsv1informers "k8s.io/client-go/informers/apps/v1"
	"k8s.io/client-go/kubernetes"
	appsv1 "k8s.io/client-go/listers/apps/v1"
	corev1 "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"log"
	"time"
)

const controllerAgentName = "custom-controller"
const fooServerQueue = "foo-server-queue"

type Controller struct {
	kubeClientSet   kubernetes.Interface
	customClientSet customclientset.Interface

	deploymentsSynced cache.InformerSynced
	deploymentsLister appsv1.DeploymentLister

	serviceSynced cache.InformerSynced
	serviceLister corev1.ServiceLister

	secretSynced cache.InformerSynced
	secretLister corev1.SecretLister

	fooServerSynced cache.InformerSynced
	fooServerLister v1beta1.FooServerLister

	workQueue workqueue.RateLimitingInterface
}

func NewController(
	kubeClientSet kubernetes.Interface,
	customClientSet customclientset.Interface,
	deploymentsInformer kubeappsv1informers.DeploymentInformer,
	fooServerInformer customv1beta1informers.FooServerInformer) *Controller {

	c := &Controller{
		kubeClientSet:   kubeClientSet,
		customClientSet: customClientSet,

		deploymentsSynced: deploymentsInformer.Informer().HasSynced,
		deploymentsLister: deploymentsInformer.Lister(),

		fooServerSynced: fooServerInformer.Informer().HasSynced,
		fooServerLister: fooServerInformer.Lister(),

		workQueue: workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), fooServerQueue),
	}

	fooServerInformer.Informer().AddEventHandler(&cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj any) {
			c.enqObjKeyToWorkQueue(obj)
		},
		UpdateFunc: func(oldObj, newObj any) {
			c.enqObjKeyToWorkQueue(newObj)
		},
	})

	return c
}

func (c *Controller) enqObjKeyToWorkQueue(obj any) {
	key, err := cache.MetaNamespaceKeyFunc(obj)
	if err != nil {
		utilruntime.HandleError(err)
		return
	}
	c.workQueue.Add(key)
}

func (c *Controller) Run(workers int, stopChan <-chan struct{}) error {
	defer utilruntime.HandleCrash()
	defer c.workQueue.ShutDown()

	if ok := cache.WaitForCacheSync(stopChan, c.deploymentsSynced, c.fooServerSynced); !ok {
		return fmt.Errorf("failed to wait for cache sync")
	}

	// spawns all the workers to process the fooServer resource
	for i := 0; i < workers; i++ {
		go wait.Until(c.runWorker, time.Second, stopChan)
	}
	log.Println("Started all workers to process fooServer resource")
	<-stopChan
	log.Println("Shutting down all workers of the fooServer resource")

	return nil
}

// runWorker is a long-running function that will continually call the
// processNextQueueItem function in order to read and process a message on the
// workQueue.
func (c *Controller) runWorker() {
	for c.processNextQueueItem() {
	}
}

func (c *Controller) processNextQueueItem() bool {
	obj, shutdown := c.workQueue.Get()
	if shutdown {
		return false
	}

	err := func(obj any) error {
		var key string
		var ok bool
		if key, ok = obj.(string); !ok {
			// As the item in the workQueue is actually invalid, we call
			// Forget here else we'd go into a loop of attempting to
			// process a work item that is invalid.
			c.workQueue.Forget(obj)
			utilruntime.HandleError(fmt.Errorf("expected string in workqueue but found %#v", obj))
			return nil
		}

		// Run the syncHandler, passing it the namespace/name string of the
		// fooServer resource to be synced.
		if err := c.fooServerSyncHandler(key); err != nil {
			// Put the item back on the workQueue to handle any transient errors.
			c.workQueue.AddRateLimited(obj)
			return fmt.Errorf("error syncing '%s': %s, requeuing", key, err.Error())
		}

		// Finally, if no error occurs we Forget this item, so it does not
		// get queued again until another change happens.
		c.workQueue.Forget(obj)
		log.Println("Successfully synced", "resourceName", key)

		return nil
	}(obj)

	if err != nil {
		utilruntime.HandleError(err)
		return true
	}

	return true
}

// fooServerSyncHandler compares the actual state with the desired, and attempts to
// converge the two. It then updates the Status block of the fooServer resource
// with the current status of the resource. If any error occurs we will requeue
// the item for later processing.
// implement the business logic here.
func (c *Controller) fooServerSyncHandler(key string) error {
	/*
		step1: convert the namespaceName to namespace & name
		step2: get the fooServer resource from the local cache with this namespace/name
				if not found or the resource is no longer exists:
					handleError and stop processing
		step3: create a deployment, secret and service from fooServer spec if not exists previously.
				if any error occurs, return error for future processing of the obj
		step4:
	*/

	fmt.Println("changed yay:", key)
	return nil
}
