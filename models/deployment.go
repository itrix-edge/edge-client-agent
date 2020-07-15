package models

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/client-go/kubernetes/typed/apps/v1"
)

type DeploymentModel struct{}

// GetDeploymentsClient function for testing NS exporting
func (m DeploymentModel) GetDeploymentsClient(namespace string) (v1.DeploymentInterface, error) {
	clientset, err := GetClientSet()
	if err != nil {
		panic(err.Error())
	}
	if len(namespace) == 0 {
		namespace = corev1.NamespaceDefault
	}
	deploymentsClient := clientset.AppsV1().Deployments(namespace)

	return deploymentsClient, nil
}

// GetDeployments get all deployments
func (m DeploymentModel) GetDeployments(namespace string, options metav1.ListOptions) (deploymentList *appsv1.DeploymentList, err error) {
	deploymentsClient, err := m.GetDeploymentsClient(namespace)
	if err != nil {
		panic(err.Error())
	}
	result, err := deploymentsClient.List(options)
	return result, err
}

// CreateDeployment create new deployment by given options
func (m DeploymentModel) CreateDeployment(namespace string, options *appsv1.Deployment) (*appsv1.Deployment, error) {
	deploymentsClient, err := m.GetDeploymentsClient(namespace)
	if err != nil {
		panic(err.Error())
	}
	result, err := deploymentsClient.Create(options.DeepCopy())
	return result, err
}

// ReadDeployment create new deployment by given options
func (m DeploymentModel) ReadDeployment(namespace string, name string, getOptions metav1.GetOptions) (*appsv1.Deployment, error) {
	deploymentsClient, err := m.GetDeploymentsClient(namespace)
	if err != nil {
		panic(err.Error())
	}
	result, err := deploymentsClient.Get(name, getOptions)
	return result, err
}

// UpdateDeplyment create new deployment by given options
func (m DeploymentModel) UpdateDeplyment(namespace string, options *appsv1.Deployment) (*appsv1.Deployment, error) {
	deploymentsClient, err := m.GetDeploymentsClient(namespace)
	if err != nil {
		panic(err.Error())
	}
	result, err := deploymentsClient.Update(options)
	return result, err
}

// DeleteDeployment delete deployment by given name
func (m DeploymentModel) DeleteDeployment(namespace string, name string, deleteOptions *metav1.DeleteOptions) bool {
	deploymentsClient, err := m.GetDeploymentsClient(namespace)
	if err != nil {
		panic(err.Error())
	}
	err = deploymentsClient.Delete(name, deleteOptions)
	if err != nil {
		return false
	} else {
		return true
	}
}
