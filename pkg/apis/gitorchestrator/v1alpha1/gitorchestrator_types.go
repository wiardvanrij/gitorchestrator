package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// GitOrchestratorSpec defines the desired state of GitOrchestrator
// +k8s:openapi-gen=true
type GitOrchestratorSpec struct {
	GitType             string `json:"gitType"`
	RepositoryName      string `json:"repositoryName"`
	RepositoryNamespace string `json:"repositoryNamespace"`
	Organisation        string `json:"organisation"`
	Description         string `json:"description"`
	Visibility          string `json:"visibility"`
}

// GitOrchestratorStatus defines the observed state of GitOrchestrator
// +k8s:openapi-gen=true
type GitOrchestratorStatus struct {
	Repositories []string `json:"repositories"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// GitOrchestrator is the Schema for the gitorchestrators API
// +k8s:openapi-gen=true
type GitOrchestrator struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GitOrchestratorSpec   `json:"spec,omitempty"`
	Status GitOrchestratorStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// GitOrchestratorList contains a list of GitOrchestrator
type GitOrchestratorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []GitOrchestrator `json:"items"`
}

func init() {
	SchemeBuilder.Register(&GitOrchestrator{}, &GitOrchestratorList{})
}
