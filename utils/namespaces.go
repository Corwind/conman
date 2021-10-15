package utils

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func CreateNamespace(client kubernetes.Clientset, name string) {
	namespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}
	client.CoreV1().Namespaces().Create(context.Background(), namespace, metav1.CreateOptions{})
}

func GetNamespace(client kubernetes.Clientset, name string) (*corev1.Namespace, error) {
	return client.CoreV1().Namespaces().Get(context.Background(), name, metav1.GetOptions{})
}
