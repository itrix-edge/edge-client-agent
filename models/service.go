package models

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	// v1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

type ServiceModel struct{}

var clientsets *kubernetes.Clientset

func (sm ServiceModel) getServiceList(namespace string, options metav1.ListOptions) (servicesList *v1.ServiceList) {
	k8s := Getk8sClient()
	servicesList, err := k8s.clientset.CoreV1().Services(namespace).List(options)
	if err != nil {
		panic(err.Error())
	}
	return servicesList
}

func (sm ServiceModel) getService(namespace string, name string, options metav1.GetOptions) (services *v1.Service) {
	k8s := Getk8sClient()
	services, err := k8s.clientset.CoreV1().Services(namespace).Get(name, options)
	if err != nil {
		panic(err.Error())
	}
	return services
}
