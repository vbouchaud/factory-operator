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

type PathFeature struct {
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:=true
	// Whether to enable vault secret creation and, if bindings are enabled, policies and user/group mappings
	Vault bool `json:"vault,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default:=true
	// Whether to enable git repository creation and, if bindings are enabled, user and group access management
	Git bool `json:"git,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default:=true
	// Whether to enable a namespace creation in the docker registry and, if bindings are enabled, user and group access management
	Registry bool `json:"registry,omitempty"`
}

type ProjectPath struct {
	// +kubebuilder:validation:Required
	Name string `json:"name"`

	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern:=`^(?:[[:alnum:]-]+\/)+[[:alnum:]-]+$`
	// The Path attribute should match the pattern ^(?:[[:alnum:]-]+\/)+[[:alnum:]-]+$ where the last [[:alnum:]-]+ is the identifier (e.g. group/sub-group/name)
	Path string `json:"path"`

	// +kubebuilder:validation:Optional
	Features PathFeature `json:"features,omitempty"`

	// +kubebuilder:validation:Optional
	Description string `json:"description,omitempty"`
}

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

type EnvironmentBindings struct {
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

type EnvironmentTarget struct {
	// +kubebuilder:validation:Required
	Namespace string `json:"namespace"`

	// +kubebuilder:validation:Required
	Cluster string `json:"cluster"`
}

type ProjectEnvironment struct {
	// +kubebuilder:validation:Required
	// Each environment should contain a name. This name will be used when creating vault secrets, in generating RoleBindings names, etc.
	Name string `json:"name"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:MinItems:=1
	// Target defines a list of namespaces clusters related to this environment
	Target []EnvironmentTarget `json:"target,omitempty"`

	// +kubebuilder:validation:Optional
	ClusterName string `json:"clusterName,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:MinItems:=1
	Bindings []EnvironmentBindings `json:"bindings,omitempty"`
}

// ProjectSpec defines the desired state of Project
type ProjectSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinItems:=1
	// A project can contains multiple paths
	Paths []ProjectPath `json:"paths"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:MinItems:=1
	// A project can contains multiple environments
	Environments []ProjectEnvironment `json:"environments,omitempty"`
}

// ProjectStatus defines the observed state of Project
type ProjectStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Conditions []metav1.Condition `json:"conditions"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Cluster

// Project is the Schema for the projects API
type Project struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ProjectSpec   `json:"spec,omitempty"`
	Status ProjectStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ProjectList contains a list of Project
type ProjectList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Project `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Project{}, &ProjectList{})
}
