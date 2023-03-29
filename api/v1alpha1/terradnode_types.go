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

// TerradNodeSpec defines the desired state of TerradNode
type TerradNodeSpec struct {
	Container    ContainerSpec `json:"container"`
	ChainId      string        `json:"chainId"`
	IsFullNode   bool          `json:"isFullNode,omitempty"`
	IsNewNetwork bool          `json:"isNewNetwork,omitempty"`
	HasPeers     bool          `json:"hasPeers,omitempty"`
	DataVolume   corev1.Volume `json:"dataVolume,omitempty"`
}

// TerradNodeStatus defines the observed state of TerradNode
type TerradNodeStatus struct {
	Nodes []string `json:"nodes"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// TerradNode is the Schema for the terradnodes API
type TerradNode struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	//TODO: Change to corev1.SecretEnvSource
	Spec   TerradNodeSpec   `json:"spec"`
	Status TerradNodeStatus `json:"status,omitempty"`
	Env    []corev1.EnvVar  `json:"env,omitempty"`
}

//+kubebuilder:object:root=true

// TerradNodeList contains a list of TerradNode
type TerradNodeList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []TerradNode `json:"items"`
}

func init() {
	SchemeBuilder.Register(&TerradNode{}, &TerradNodeList{})
}
