package models

import (
	rv1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type NamespaceModel struct{}

// GetNamespaces function for testing NS exporting
func (m NamespaceModel) GetNamespaces(opts v1.ListOptions) (rv1.NamespaceList, error) {
	// creates the clientset
	clientset, err := GetClientSet()
	if err != nil {
		panic(err.Error())
	}
	list, err := clientset.CoreV1().Namespaces().List(opts)
	if err != nil {
		return *list, err
	}
	return *list, nil
}

// GetNamespace function for testing NS exporting
func (m NamespaceModel) GetNamespace(namespace string, opts v1.GetOptions) (*rv1.Namespace, error) {
	// creates the clientset
	clientset, err := GetClientSet()
	if err != nil {
		panic(err.Error())
	}
	ns, err := clientset.CoreV1().Namespaces().Get(namespace, opts)
	if err != nil {
		return ns, err
	}
	return ns, nil
}

// CreateNamespace function for testing NS exporting
func (m NamespaceModel) CreateNamespace(opts *rv1.Namespace) (*rv1.Namespace, error) {
	// creates the clientset
	clientset, err := GetClientSet()
	if err != nil {
		panic(err.Error())
	}
	ns, err := clientset.CoreV1().Namespaces().Create(opts)
	return ns, nil
}
