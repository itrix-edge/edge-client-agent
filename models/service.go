package models

import (
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	// v1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

type ServiceModel struct{}

var clientsets *kubernetes.Clientset

func (m ServiceModel) GetServiceList(namespace string, options metav1.ListOptions) (*core.ServiceList, error) {
	k8s := Getk8sClient()
	servicesList, err := k8s.clientset.CoreV1().Services(namespace).List(options)
	return servicesList, err
}

func (m ServiceModel) GetService(namespace string, name string, options metav1.GetOptions) (*core.Service, error) {
	k8s := Getk8sClient()
	services, err := k8s.clientset.CoreV1().Services(namespace).Get(name, options)
	return services, err
}

func (m ServiceModel) CreateService(namespace string, options *core.Service) (*core.Service, error) {
	k8s := Getk8sClient()
	services, err := k8s.clientset.CoreV1().Services(namespace).Create(options)
	return services, err
}
