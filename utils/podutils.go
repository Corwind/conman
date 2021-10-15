package utils

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func PodList(k8s_client *kubernetes.Clientset, namespace string) (*corev1.PodList, error) {
	return k8s_client.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
}

func PodGet(k8s_client *kubernetes.Clientset, namespace string, name string) (*corev1.Pod, error) {
	return k8s_client.CoreV1().Pods(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}
