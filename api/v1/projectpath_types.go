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

// ProjectPathSpec defines the desired state of ProjectPath
type ProjectPathSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern:=`^(?:[[:alnum:]-]+\/)+[[:alnum:]-]+$`
	// The Path attribute should match the pattern ^(?:[[:alnum:]-]+\/)+[[:alnum:]-]+$ where the last [[:alnum:]-]+ is the identifier (e.g. group/sub-group/name)
	Path string `json:"path"`

	// +kubebuilder:validation:Optional
	Features PathFeature `json:"features,omitempty"`

	// +kubebuilder:validation:Optional
	Description string `json:"description,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default:=true
	// When service allows for it, archive created resources instead of deleting them
	ArchiveOnDelete bool `json:"archiveondelete,omitempty"`
}

// ProjectPathStatus defines the observed state of ProjectPath
type ProjectPathStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Conditions []metav1.Condition `json:"conditions"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Cluster

// ProjectPath is the Schema for the projectpaths API
type ProjectPath struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ProjectPathSpec   `json:"spec,omitempty"`
	Status ProjectPathStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ProjectPathList contains a list of ProjectPath
type ProjectPathList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ProjectPath `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ProjectPath{}, &ProjectPathList{})
}
