package models

import (
	"k8s.io/client-go/kubernetes"

	// v1 "k8s.io/client-go/kubernetes/typed/core/v1"

	"k8s.io/client-go/rest"
)

// ClientSet define k8s client set
type ClientSet struct {
	clientset *kubernetes.Clientset
}

var client *ClientSet

func (cs *ClientSet) init() {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset once
	if cs.clientset == nil {
		cs.clientset, err = kubernetes.NewForConfig(config)
		if err != nil {
			panic(err.Error())
		}
	}
}

// Getk8sClient Get kubernetes client set
func Getk8sClient() *ClientSet {
	if client == nil {
		client = new(ClientSet)
	}
	return client
}
