package kubeconfig

import (
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"testing"
)

func TestInitK8S(t *testing.T) {

	config := InitK8S()
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		t.Error(err)
	}
	pods, err := clientset.CoreV1().Pods("kube-system").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		t.Error(err)
	}
	for _, pod := range pods.Items {
		fmt.Println(pod.Name)
	}
}
