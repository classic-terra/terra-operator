package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ValidatorSpec defines the desired state of Validator
type ValidatorSpec struct {
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Important: Run "operator-sdk generate crds" to regenerate crds after modifying this file
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

// ValidatorStatus defines the observed state of Validator
type ValidatorStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Important: Run "operator-sdk generate crds" to regenerate crds after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Validator is the Schema for the validators API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=validators,scope=Namespaced
type Validator struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	EnvFrom []corev1.EnvFromSource `json:"envFrom,omitempty"`
	Spec    ValidatorSpec          `json:"spec,omitempty"`
	Status  ValidatorStatus        `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ValidatorList contains a list of Validator
type ValidatorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Validator `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Validator{}, &ValidatorList{})
}
