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

// TerradNetSpec defines the desired state of TerradNet
// TerradNet is only viable for simulating a tesnet. It cannot be used for mainnet.
type TerradNetSpec struct {
	Replica     int32         `json:"replica,omitempty"`
	Container   ContainerSpec `json:"container"`
	ChainId     string        `json:"chainId"`
	ServiceName string        `json:"serviceName"`
	DataSource  corev1.Volume `json:"dataSource"`
}

// TerradNetStatus defines the observed state of TerradNet
type TerradNetStatus struct {
	Nodes []string `json:"nodes"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// TerradNet is the Schema for the TerradNets API
type TerradNet struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TerradNetSpec   `json:"spec,omitempty"`
	Status TerradNetStatus `json:"status,omitempty"`
	Env    []corev1.EnvVar `json:"env,omitempty"`
}

//+kubebuilder:object:root=true

// TerradNetList contains a list of TerradNet
type TerradNetList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []TerradNet `json:"items"`
}

func init() {
	SchemeBuilder.Register(&TerradNet{}, &TerradNetList{})
}
