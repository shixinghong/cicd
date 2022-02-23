package main

import (
	"os"

	"github.com/hongshixing/cicd/pkg/apis/task/v1alpha1"
	"github.com/hongshixing/cicd/pkg/client/clientset/versioned"
	controller "github.com/hongshixing/cicd/pkg/controllers"
	"github.com/hongshixing/cicd/pkg/kubeconfig"

	"sigs.k8s.io/controller-runtime/pkg/builder"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
)

func main() {
	logf.SetLogger(zap.New())
	var (
		log        = logf.Log.WithName("ci")
		kubeConfig = kubeconfig.InitK8S()
		taskClient = versioned.NewForConfigOrDie(kubeConfig)
	)
	mgr, err := manager.New(kubeConfig, manager.Options{
		Logger: log,
	})
	if err != nil {
		log.Error(err, "could not create manager")
		os.Exit(1)
	}

	if err = v1alpha1.SchemeBuilder.AddToScheme(mgr.GetScheme()); err != nil {
		log.Error(err, "could not add manager")
		os.Exit(1)
	}

	taskController := controller.NewTaskController(
		mgr.GetEventRecorderFor("ci"),
		taskClient,
	)
	if err = builder.ControllerManagedBy(mgr).For(&v1alpha1.Task{}).Complete(taskController); err != nil {
		log.Error(err, "could not create controller")
		os.Exit(1)
	}

	if err = mgr.Start(signals.SetupSignalHandler()); err != nil {
		log.Error(err, "could not start manager")
		os.Exit(1)
	}
}
