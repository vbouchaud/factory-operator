/*
Copyright 2021.

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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

type EnvironmentTarget struct {
	// +kubebuilder:validation:Required
	Namespace string `json:"namespace"`

	// +kubebuilder:validation:Required
	Cluster string `json:"cluster"`
}

// ProjectEnvironmentSpec defines the desired state of ProjectEnvironment
type ProjectEnvironmentSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// +kubebuilder:validation:Required
	// Each environment should contain a name. This name will be used when creating vault secrets, in generating RoleBindings names, etc.
	Environment string `json:"environment"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:MinItems:=1
	// Target defines a list of namespaces clusters related to this environment
	Target []EnvironmentTarget `json:"target,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:MinItems:=1
	Bindings []BindingSpec `json:"bindings,omitempty"`
}

// ProjectEnvironmentStatus defines the observed state of ProjectEnvironment
type ProjectEnvironmentStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Cluster

// ProjectEnvironment is the Schema for the projectenvironments API
type ProjectEnvironment struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ProjectEnvironmentSpec   `json:"spec,omitempty"`
	Status ProjectEnvironmentStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ProjectEnvironmentList contains a list of ProjectEnvironment
type ProjectEnvironmentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ProjectEnvironment `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ProjectEnvironment{}, &ProjectEnvironmentList{})
}
