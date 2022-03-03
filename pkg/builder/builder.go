package builder

import (
	"context"

	"github.com/hongshixing/cicd/pkg/apis/task/v1alpha1"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type PodBuilder struct {
	task *v1alpha1.Task
	client.Client
}

func NewPodBuilder(task *v1alpha1.Task, client client.Client) *PodBuilder {
	return &PodBuilder{
		task:   task,
		Client: client,
	}
}

func (pb *PodBuilder) Build(ctx context.Context) error {
	newPod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "task-pod-" + pb.task.Name,
			Namespace: pb.task.Namespace,
		},
	}

	newPod.Spec.RestartPolicy = corev1.RestartPolicyNever // 用不重启
	var containers []corev1.Container
	for _, step := range pb.task.Spec.Steps {
		step.Container.ImagePullPolicy = corev1.PullIfNotPresent // 强制要求拉取策略
		containers = append(containers, step.Container)
	}

	newPod.Spec.Containers = containers
	newPod.OwnerReferences = append(newPod.OwnerReferences,
		metav1.OwnerReference{
			APIVersion: pb.task.APIVersion,
			Kind:       pb.task.Kind,
			Name:       pb.task.Name,
			UID:        pb.task.UID,
		},
	)
	newPod.Annotations = make(map[string]string)
	newPod.Annotations["task-order"] = "0"
	return pb.Client.Create(ctx, newPod)

}
