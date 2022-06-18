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

// OracleFeederSpec defines the desired state of OracleFeeder
type OracleFeederSpec struct {
	NodeImage string `json:"nodeImage"`
}

// OracleFeederStatus defines the observed state of OracleFeeder
type OracleFeederStatus struct {
	Oracles []string `json:"oracles"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// OracleFeeder is the Schema for the oraclefeeders API
type OracleFeeder struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	//TODO: Change to corev1.SecretEnvSource
	Env    []corev1.EnvVar    `json:"env"`
	Spec   OracleFeederSpec   `json:"spec"`
	Status OracleFeederStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// OracleFeederList contains a list of OracleFeeder
type OracleFeederList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []OracleFeeder `json:"items"`
}

func init() {
	SchemeBuilder.Register(&OracleFeeder{}, &OracleFeederList{})
}
