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

type BindingSubject struct {
	// +kubebuilder:validation:Required
	Name string `json:"name"`

	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Enum:=Group;User
	Kind string `json:"kind"`
}

type BindingRole struct {
	// +kubebuilder:validation:Required
	Name string `json:"name"`

	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Enum:=ClusterRole;Role
	Kind string `json:"kind"`
}

type BindingAccess struct {
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:=true
	Logging bool `json:"logging,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default:=true
	Monitoring bool `json:"monitoring,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default:=true
	Registry bool `json:"registry,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default:=true
	Git bool `json:"git,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default:=true
	Vault bool `json:"vault,omitempty"`
}

// BindingSpec defines the desired state of Binding
type BindingSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// +kubebuilder:validation:Required
	Name string `json:"name"`

	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Enum:=Owner;Manager;Developer
	Kind string `json:"kind"`

	// +kubebuilder:validation:Optional
	Access BindingAccess `json:"access,omitempty"`

	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinItems:=1
	Subjects []BindingSubject `json:"subjects"`

	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinItems:=1
	Roles []BindingRole `json:"roles"`
}

// BindingStatus defines the observed state of Binding
type BindingStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Conditions []metav1.Condition `json:"conditions"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Cluster

// Binding is the Schema for the bindings API
type Binding struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BindingSpec   `json:"spec,omitempty"`
	Status BindingStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// BindingList contains a list of Binding
type BindingList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Binding `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Binding{}, &BindingList{})
}
