package models

import (
	"k8s.io/client-go/kubernetes"
	// v1 "k8s.io/client-go/kubernetes/typed/core/v1"
	rv1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
)

type NamespaceModel struct{}

var clientset kubernetes.Clientset

// func (m NamespaceModel) init() {
// 	config, err := rest.InClusterConfig()
// 	if err != nil {
// 		panic(err.Error())
// 	}
// 	// creates the clientset
// 	localclientset, err := kubernetes.NewForConfig(config)
// 	if err != nil {
// 		panic(err.Error())
// 	}
// 	clientset = *localclientset
// }

// GetNamespaces function for testing NS exporting
func (m NamespaceModel) GetNamespaces(opts v1.ListOptions) (rv1.NamespaceList, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	list, err := clientset.CoreV1().Namespaces().List(opts)
	if err != nil {
		return *list, err
	}
	return *list, nil
}
