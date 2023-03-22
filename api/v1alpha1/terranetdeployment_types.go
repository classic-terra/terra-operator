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

// TerraNetDeploymentSpec defines the desired state of TerraNetDeployment
// TerraNet is only viable for simulating a tesnet. It cannot be used for mainnet.
type TerraNetDeploymentSpec struct {
	Replica    int32         `json:"replica"`
	Container  ContainerSpec `json:"container"`
	ChainId    string        `json:"chainId"`
	DataVolume corev1.Volume `json:"dataVolume,omitempty"`
}

// TerraNetDeploymentStatus defines the observed state of TerraNetDeployment
type TerraNetDeploymentStatus struct {
	Nodes []string `json:"nodes"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// TerraNetDeployment is the Schema for the terranetdeployments API
type TerraNetDeployment struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TerraNetDeploymentSpec   `json:"spec,omitempty"`
	Status TerraNetDeploymentStatus `json:"status,omitempty"`
	Env    []corev1.EnvVar          `json:"env,omitempty"`
}

//+kubebuilder:object:root=true

// TerraNetDeploymentList contains a list of TerraNetDeployment
type TerraNetDeploymentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []TerraNetDeployment `json:"items"`
}

func init() {
	SchemeBuilder.Register(&TerraNetDeployment{}, &TerraNetDeploymentList{})
}
