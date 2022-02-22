package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Task task_config
type Task struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec TaskSpec `json:"spec"`
	// +optional
	Status TaskStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// TaskList db_configs
type TaskList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []Task `json:"items"`
}

type TaskSpec struct {
	//Replicas int32  `json:"replicas,omitempty"`
	//Dsn      string `json:"dsn,omitempty"`
	Steps []TaskStep `json:"steps,omitempty"`
}

type TaskStep struct {
	corev1.Container `json:",inline"`
}

type TaskStatus struct {
	Replicas int32  `json:"replicas,omitempty"`
	Ready    string `json:"ready,omitempty"`
}
