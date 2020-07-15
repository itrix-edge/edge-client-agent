package models

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	typev1 "k8s.io/client-go/kubernetes/typed/apps/v1"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type KuberClient struct{}

var config *rest.Config
var clientset *kubernetes.Clientset

// var deploymentsClient *typev1.DeploymentInterface
// var serviceClient *v1.ServiceInterface

func InitKuberClient(kubeconfig *string) {
	if kubeconfig == nil {
		config = GetInClusterConfig()
	} else {
		config = GetOutOfClusterConfig(kubeconfig)
	}
}

func GetInClusterConfig() *rest.Config {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	return config
}

func GetOutOfClusterConfig(kubeconfig *string) *rest.Config {
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	// config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	return config
}

// func (m KuberClient) GetOutOfClusterConfig(kubeconfig *string) *rest.Config {
// 	// var kubeconfig *string
// 	if home := m.homeDir(); home != "" {
// 		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
// 	} else {
// 		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
// 	}
// 	flag.Parse()

// 	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
// 	if err != nil {
// 		panic(err.Error())
// 	}
// 	return config
// }

// func (m KuberClient) homeDir() string {
// 	if h := os.Getenv("HOME"); h != "" {
// 		return h
// 	}
// 	return os.Getenv("USERPROFILE") // windows
// }

// GetClientSet get kubernetes client set object
func GetClientSet() (*kubernetes.Clientset, error) {
	if clientset == nil {
		var err error
		clientset, err = kubernetes.NewForConfig(config)
		return clientset, err
	}
	return clientset, nil
}

func (m KuberClient) GetDeploymentClient(namespace string) (typev1.DeploymentInterface, error) {
	clientset, err := GetClientSet()
	namespace = m.FilterNamespace(namespace)
	deploymentsClient := clientset.AppsV1().Deployments(namespace)
	return deploymentsClient, err
}

func (m KuberClient) GetServiceClient(namespace string) (v1.ServiceInterface, error) {
	clientset, err := GetClientSet()
	namespace = m.FilterNamespace(namespace)
	serviceClient := clientset.CoreV1().Services(namespace)
	return serviceClient, err
}

func (m KuberClient) FilterNamespace(namespace string) string {
	if len(namespace) == 0 {
		namespace = corev1.NamespaceDefault
	}
	return namespace
}
