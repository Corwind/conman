package clients

import (
	"fmt"

	helmclient "github.com/mittwald/go-helm-client"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func K8sClient(kubeconfig []byte) *kubernetes.Clientset {
	config, err := clientcmd.RESTConfigFromKubeConfig(kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)

	if err != nil {
		panic(err.Error())
	}
	return clientset
}

func GetHelmClient(kubeconfig []byte, namespace string) helmclient.Client {
	opt := &helmclient.KubeConfClientOptions{
		Options: &helmclient.Options{
			Namespace:        namespace,
			RepositoryCache:  fmt.Sprintf("/tmp/.%s-helmcache", namespace),
			RepositoryConfig: fmt.Sprintf("/tmp/.%s-helmrepo", namespace),
		},
		KubeContext: "",
		KubeConfig:  kubeconfig,
	}

	helm_client, err := helmclient.NewClientFromKubeConf(opt)
	if err != nil {
		panic(err.Error())
	}
	return helm_client
}
