package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TerradNodeSpec defines the desired state of TerradNode
type TerradNodeSpec struct {
	NodeImage        string        `json:"nodeImage"`
	IsFullNode       bool          `json:"isFullNode,omitempty"`
	DataVolume       corev1.Volume `json:"dataVolume,omitempty"`
	PostStartCommand []string      `json:"postStartCommand,omitempty"`
}

// TerradNodeStatus defines the observed state of TerradNode
type TerradNodeStatus struct {
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// TerradNode is the Schema for the terradnodes API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=terradnodes,scope=Namespaced
type TerradNode struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	EnvFrom []corev1.EnvFromSource `json:"envFrom,omitempty"`
	Spec    TerradNodeSpec         `json:"spec,omitempty"`
	Status  TerradNodeStatus       `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// TerradNodeList contains a list of TerradNode
type TerradNodeList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []TerradNode `json:"items"`
}

func init() {
	SchemeBuilder.Register(&TerradNode{}, &TerradNodeList{})
}
