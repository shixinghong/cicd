package controller

import (
	"context"
	"github.com/hongshixing/cicd/pkg/apis/task/v1alpha1"
	"github.com/hongshixing/cicd/pkg/builder"
	clientset "github.com/hongshixing/cicd/pkg/client/clientset/versioned"

	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	APIVersion      = "v1alpha1"
	Kind            = "task"
	Group           = "api.myit.fun"
	GroupAPIVersion = Group + "/" + APIVersion
)

type TaskController struct {
	E record.EventRecorder // 记录事件
	*clientset.Clientset
	client.Client
}

func NewTaskController(e record.EventRecorder, clientset *clientset.Clientset) *TaskController {
	return &TaskController{
		E:         e,
		Clientset: clientset,
	}
}
func (t *TaskController) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
	_ = logf.FromContext(ctx)
	task := &v1alpha1.Task{}
	err := t.Client.Get(ctx, req.NamespacedName, task)
	if err != nil {
		return reconcile.Result{}, nil
	}

	pb := builder.NewPodBuilder(task, t.Client)
	if err = pb.Build(ctx); err != nil {
		return reconcile.Result{}, err
	}
	return reconcile.Result{}, nil
}

func (t *TaskController) InjectClient(c client.Client) error {
	t.Client = c
	return nil
}
