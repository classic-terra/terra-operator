/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ValidatorSpec defines the desired state of Validator
type ValidatorSpec struct {
	ChainId                 string        `json:"chainId"`
	IsNewNetwork            bool          `json:"isNewNetwork,omitempty"`
	TerradNodeImage         string        `json:"terradNodeImage"`
	OracleNodeImages        []string      `json:"oracleNodeImages"`
	IndexerNodeImages       []string      `json:"indexerNodeImages"`
	Passphrase              string        `json:"passphrase"`
	Mnenomic                string        `json:"mnenomic"`
	Amount                  string        `json:"amount"`
	CommissionRate          string        `json:"commissionRate"`
	CommissionRateMax       string        `json:"commissionRateMax"`
	CommissionRateMaxChange string        `json:"commissionRateMaxChange"`
	MinimumSelfDelegation   string        `json:"minimumSelfDelegation"`
	AutoConfig              bool          `json:"autoConfig,omitempty"`
	IsPublic                bool          `json:"isPublic,omitempty"`
	Website                 string        `json:"website,omitempty"`
	Description             string        `json:"description,omitempty"`
	DataVolume              corev1.Volume `json:"dataVolume,omitempty"`
}

// ValidatorStatus defines the observed state of Validator
type ValidatorStatus struct {
	Nodes []string `json:"nodes"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Validator is the Schema for the validators API
type Validator struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ValidatorSpec   `json:"spec,omitempty"`
	Status ValidatorStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ValidatorList contains a list of Validator
type ValidatorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Validator `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Validator{}, &ValidatorList{})
}
