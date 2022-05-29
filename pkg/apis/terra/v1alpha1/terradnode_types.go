package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// TerradNodeSpec defines the desired state of TerradNode
type TerradNodeSpec struct {
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	PubKey                string `json:"pubKey"`
	Address               string `json:"address"`
	Name                  string `json:"name"`
	Website               string `json:"website,omitempty"`
	Description           string `json:"description,omitempty"`
	InitialCommissionRate string `json:"initialCommissionRate"`
	MaximumCommission     string `json:"maximumCommission"`
	CommissionChangeRate  string `json:"commissionChangeRate"`
	MinimumSelfBondAmount string `json:"minimumSelfBondAmount"`
	InitialSelfBondAmount string `json:"initialSelfBondAmount"`
}

// TerradNodeStatus defines the observed state of TerradNode
type TerradNodeStatus struct {
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
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
