package controller

import (
	"context"
	"fmt"
	"github.com/imtiaz246/custom-controller/pkg/apis/cho.me/v1beta1"
	customclientset "github.com/imtiaz246/custom-controller/pkg/client/clientset/versioned"
	customv1beta1informers "github.com/imtiaz246/custom-controller/pkg/client/informers/externalversions/cho.me/v1beta1"
	customlisters "github.com/imtiaz246/custom-controller/pkg/client/listers/cho.me/v1beta1"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	kubeappsv1informers "k8s.io/client-go/informers/apps/v1"
	kubecorev1informers "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/kubernetes"
	appsv1listers "k8s.io/client-go/listers/apps/v1"
	corev1listers "k8s.io/client-go/listers/core/v1"
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
	deploymentsLister appsv1listers.DeploymentLister

	serviceSynced cache.InformerSynced
	serviceLister corev1listers.ServiceLister

	secretSynced cache.InformerSynced
	secretLister corev1listers.SecretLister

	fooServerSynced cache.InformerSynced
	fooServerLister customlisters.FooServerLister

	workQueue workqueue.RateLimitingInterface
}

func NewController(
	kubeClientSet kubernetes.Interface,
	customClientSet customclientset.Interface,
	deploymentsInformer kubeappsv1informers.DeploymentInformer,
	servicesInformer kubecorev1informers.ServiceInformer,
	secretsInformer kubecorev1informers.SecretInformer,
	fooServerInformer customv1beta1informers.FooServerInformer) *Controller {

	c := &Controller{
		kubeClientSet:   kubeClientSet,
		customClientSet: customClientSet,

		deploymentsSynced: deploymentsInformer.Informer().HasSynced,
		deploymentsLister: deploymentsInformer.Lister(),

		serviceSynced: servicesInformer.Informer().HasSynced,
		serviceLister: servicesInformer.Lister(),

		secretSynced: secretsInformer.Informer().HasSynced,
		secretLister: secretsInformer.Lister(),

		fooServerSynced: fooServerInformer.Informer().HasSynced,
		fooServerLister: fooServerInformer.Lister(),

		workQueue: workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), fooServerQueue),
	}
	// event handler for fooServer resource events
	fooServerInformer.Informer().AddEventHandler(&cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj any) {
			c.enqFooServerKeyToWorkQueue(obj)
		},
		UpdateFunc: func(oldObj, newObj any) {
			c.enqFooServerKeyToWorkQueue(newObj)
		},
	})
	// event handler for deployment resource events
	deploymentsInformer.Informer().AddEventHandler(&cache.ResourceEventHandlerFuncs{
		AddFunc: c.handleInformedObject,
		UpdateFunc: func(old, new any) {
			// check the resource version, if not same call handleInformedObject
			oldDeploy := old.(*appsv1.Deployment)
			newDeploy := new.(*appsv1.Deployment)

			if oldDeploy.ResourceVersion == newDeploy.ResourceVersion {
				return
			}

			c.handleInformedObject(new)
		},
		DeleteFunc: c.handleInformedObject,
	})
	// todo: event handler for service resource events
	// todo: event handler for secret resource events

	return c
}

func (c *Controller) handleInformedObject(obj any) {
	var object metav1.Object // the object implements metav1.Object interface
	var ok bool

	//todo: understand
	if object, ok = obj.(metav1.Object); !ok {
		tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
		if !ok {
			utilruntime.HandleError(fmt.Errorf("error decoding object, invalid type"))
			return
		}
		object, ok = tombstone.Obj.(metav1.Object)
		if !ok {
			utilruntime.HandleError(fmt.Errorf("error decoding object tombstone, invalid type"))
			return
		}
		log.Println("Recovered deleted object", "resourceName", object.GetName())
	}

	/*
		- Check the owner is one of FooServer kind
		- Get the owner resource with the name
		- If the owner resource is still exits, enqueue to workQueue
	*/
	if ownerRef := metav1.GetControllerOf(object); ownerRef != nil {
		if ownerRef.Kind != "FooServer" {
			return
		}
		fooServerResource, err := c.fooServerLister.FooServers("default").Get(ownerRef.Name)
		if err != nil {
			log.Println("Ignore orphaned object", "object: ", object.GetName(), "ownerName: ", ownerRef.Name)
			return
		}
		c.enqFooServerKeyToWorkQueue(fooServerResource)
	}
}

func (c *Controller) enqFooServerKeyToWorkQueue(obj any) {
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
		- convert the namespaceName to namespace & name
		- get the fooServer resource and check the spec fields
		- get the resources owned by the fooServer resource
		- if the owned resources not found create them in order
		- secret -> deployment -> service
		- compare the stat & sync the state
		- update the status sub-resource
	*/
	// convert the namespaceName to namespace & nam
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		utilruntime.HandleError(fmt.Errorf("invalid resource key %v", key))
		return nil
	}
	// get the fooServer resource and check the spec fields
	fooServer, err := c.fooServerLister.FooServers(namespace).Get(name)
	if err != nil {
		// The fooServer resource may no longer exist, in which case we stop processing.
		if errors.IsNotFound(err) {
			utilruntime.HandleError(fmt.Errorf("foo '%s' in work queue no longer exists", key))
			return nil
		}
	}
	deploymentName := fooServer.Spec.DeploymentSpec.Name
	serviceName := fooServer.Spec.ServiceSpec.Name
	secretName := fooServer.Spec.SecretSpec.Name

	// todo: use webhook for these kind of validation or use client-gen validation
	if deploymentName == "" {
		utilruntime.HandleError(fmt.Errorf("%s: deployment name must be specified", key))
		return nil
	}
	if serviceName == "" {
		utilruntime.HandleError(fmt.Errorf("%s: service name must be specified", key))
		return nil
	}
	if secretName == "" {
		utilruntime.HandleError(fmt.Errorf("%s: secret name must be specified", key))
	}

	// create secret resource if not exists
	secret, err := c.secretLister.Secrets(fooServer.Namespace).Get(secretName)
	if errors.IsNotFound(err) {
		secret, err = c.kubeClientSet.CoreV1().Secrets(fooServer.Namespace).Create(context.TODO(), newSecret(fooServer), metav1.CreateOptions{})
	}
	if err != nil {
		return err
	}

	// create deployment resource if not exists
	deployment, err := c.deploymentsLister.Deployments(fooServer.Namespace).Get(deploymentName)
	if errors.IsNotFound(err) {
		deployment, err = c.kubeClientSet.AppsV1().Deployments(fooServer.Namespace).Create(context.TODO(), newDeployment(fooServer), metav1.CreateOptions{})
	}
	if err != nil {
		return err
	}

	// create service resource if not exists
	service, err := c.serviceLister.Services(fooServer.Namespace).Get(serviceName)
	if errors.IsNotFound(err) {
		service, err = c.kubeClientSet.CoreV1().Services(fooServer.Namespace).Create(context.TODO(), newService(fooServer), metav1.CreateOptions{})
	}
	if err != nil {
		return err
	}

	// reconcile resources
	fooServerSecSpec := fooServer.Spec.SecretSpec
	fooServerSerSpec := fooServer.Spec.ServiceSpec
	fooServerDepSpec := fooServer.Spec.DeploymentSpec

	if fooServerSecSpec.Username != string(secret.Data["ADMIN_USERNAME"]) || fooServerSecSpec.Password != string(secret.Data["ADMIN_PASSWORD"]) {
		secret.Data["ADMIN_USERNAME"] = []byte(fooServerSecSpec.Username)
		secret.Data["ADMIN_PASSWORD"] = []byte(fooServerSecSpec.Username)
		secret, err = c.kubeClientSet.CoreV1().Secrets(fooServer.Namespace).Update(context.TODO(), secret, metav1.UpdateOptions{})
		if err != nil {
			return err
		}
	}

	if fooServerDepSpec.PodReplicas != nil && fooServerDepSpec.PodReplicas != deployment.Spec.Replicas {
		deployment.Spec.Replicas = fooServerDepSpec.PodReplicas
		deployment, err = c.kubeClientSet.AppsV1().Deployments(fooServer.Namespace).Update(context.TODO(), deployment, metav1.UpdateOptions{})
		if err != nil {
			return err
		}
	}
	if fooServerDepSpec.PodContainerPort != nil && *fooServerDepSpec.PodContainerPort != deployment.Spec.Template.Spec.Containers[0].Ports[0].ContainerPort {
		deployment.Spec.Template.Spec.Containers[0].Ports[0].ContainerPort = *fooServerDepSpec.PodContainerPort
		deployment, err = c.kubeClientSet.AppsV1().Deployments(fooServer.Namespace).Update(context.TODO(), deployment, metav1.UpdateOptions{})
		if err != nil {
			return err
		}
	}

	if fooServerSerSpec.NodePort != nil && *fooServerSerSpec.NodePort != service.Spec.Ports[0].NodePort {
		service.Spec.Ports[0].NodePort = *fooServerSerSpec.NodePort
		service, err = c.kubeClientSet.CoreV1().Services(fooServer.Namespace).Update(context.TODO(), service, metav1.UpdateOptions{})
		if err != nil {
			return err
		}
	}

	// todo: reconcile other service spec

	err = c.updateFooServerStatus(deployment, fooServer)
	if err != nil {
		return err
	}

	return nil
}

func (c *Controller) updateFooServerStatus(deployment *appsv1.Deployment, fooServer *v1beta1.FooServer) error {
	fooServerCopy := fooServer.DeepCopy()
	fooServerCopy.Status.AvailableReplicas = *deployment.Spec.Replicas

	_, err := c.customClientSet.ChoV1beta1().FooServers(fooServer.Namespace).Update(context.TODO(), fooServerCopy, metav1.UpdateOptions{})
	return err
}
