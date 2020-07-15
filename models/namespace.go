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
