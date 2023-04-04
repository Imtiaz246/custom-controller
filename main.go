package main

import (
	"github.com/imtiaz246/custom-controller/controller"
	customclientset "github.com/imtiaz246/custom-controller/pkg/client/clientset/versioned"
	custominformers "github.com/imtiaz246/custom-controller/pkg/client/informers/externalversions"
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"os"
	"path/filepath"
	"time"
)

func main() {
	kubeConfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
	log.Println(kubeConfig)

	config, err := clientcmd.BuildConfigFromFlags("", kubeConfig)
	if err != nil {
		log.Printf("error %s building config from kubeconfig\n", err.Error())
	}

	kubeClientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("erorr %s building clientset from config\n", err.Error())
	}

	customClientSet, err := customclientset.NewForConfig(config)
	if err != nil {
		log.Fatalf("erorr %s building clientset from config\n", err.Error())
	}

	kubeInformerFactory := kubeinformers.NewSharedInformerFactory(kubeClientSet, time.Second*30)
	customInformerFactory := custominformers.NewSharedInformerFactory(customClientSet, time.Second*30)

	Controller := controller.NewController(
		kubeClientSet,
		customClientSet,
		kubeInformerFactory.Apps().V1().Deployments(),
		customInformerFactory.Cho().V1beta1().FooServers())

	stopChan := make(chan struct{})
	kubeInformerFactory.Start(stopChan)
	customInformerFactory.Start(stopChan)
	Controller.Run(5, stopChan)
	<-stopChan
}
