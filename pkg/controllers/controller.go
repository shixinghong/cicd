package controller

import (
	"context"
	//"github.com/hongshixing/cicd/pkg/builders"

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
	client.Client
}

func NewTaskController() *TaskController {
	return &TaskController{}
}

func (t *TaskController) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
	_ = logf.FromContext(ctx)

	return reconcile.Result{}, nil
}

func (t *TaskController) InjectClient(c client.Client) error {
	t.Client = c
	return nil
}
