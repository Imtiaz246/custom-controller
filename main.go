package main

import (
	customclientset "github.com/imtiaz246/custom-cntroller/pkg/client/clientset/versioned"
	custominformers "github.com/imtiaz246/custom-cntroller/pkg/client/informers/externalversions"
	"k8s.io/client-go/informers"
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

	exampleClientSet, err := customclientset.NewForConfig(config)
	if err != nil {
		log.Fatalf("erorr %s building clientset from config\n", err.Error())
	}

	kubeInformerFactory := informers.NewSharedInformerFactory(kubeClientSet, time.Second*30)
	exampleInformerFactory := custominformers.NewSharedInformerFactory(exampleClientSet, time.Second*30)

	Controller := NewController()

	stopChan := make(chan struct{})
	kubeInformerFactory.Start(stopChan)
	exampleInformerFactory.Start(stopChan)
	Controller.Run(5, stopChan)
	<-stopChan
}
