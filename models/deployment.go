package models

import (
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	v1 "k8s.io/client-go/kubernetes/typed/apps/v1"

	"k8s.io/client-go/rest"
)

type DeploymentModel struct{}

type DeploymentOptions struct {
	namespace string
	image     string
	name      string
	ports     []int
}

// var clientset kubernetes.Clientset

// func (m DeploymentModel) init() {
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

// GetDeploymentsClient function for testing NS exporting
func (m DeploymentModel) GetDeploymentsClient(namespace string) (v1.DeploymentInterface, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	if len(namespace) == 0 {
		namespace = apiv1.NamespaceDefault
	}
	deploymentsClient := clientset.AppsV1().Deployments(namespace)

	return deploymentsClient, nil
}

// GetDeployments get all deployments
func (m DeploymentModel) GetDeployments(options metav1.GetOptions) (deployment *appsv1.Deployment, err error) {
	namespace := "default"
	deploymentsClient, err := m.GetDeploymentsClient(namespace)
	if err != nil {
		panic(err.Error())
	}
	result, err := deploymentsClient.Get(namespace, options)
	return result, err
}

// AddNewDeployment create new deployment by given options
func (m DeploymentModel) AddNewDeployment(namespace string, options *appsv1.Deployment) (*appsv1.Deployment, error) {
	deploymentsClient, err := m.GetDeploymentsClient(namespace)
	if err != nil {
		panic(err.Error())
	}
	result, err := deploymentsClient.Create(options.DeepCopy())
	return result, err
}

func (m DeploymentModel) GetDeploymentOption(plainOptions DeploymentOptions) *appsv1.Deployment {
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "deployment",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(2),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": plainOptions.name,
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": plainOptions.name,
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  plainOptions.name,
							Image: plainOptions.image,
							Ports: []apiv1.ContainerPort{
								{
									Name:          "http",
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: 80,
								},
							},
						},
					},
				},
			},
		},
	}
	return deployment
}

func int32Ptr(i int32) *int32 { return &i }
