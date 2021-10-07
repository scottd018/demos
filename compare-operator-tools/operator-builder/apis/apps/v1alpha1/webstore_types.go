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

package v1alpha1

import (
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/scottd018/demos/apis/common"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// WebStoreSpec defines the desired state of WebStore.
type WebStoreSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// +kubebuilder:default="nginx:1.17"
	// +kubebuilder:validation:Optional
	// Defines the web store image
	WebstoreImage string `json:"webstoreImage"`

	// +kubebuilder:default="webstore-svc"
	// +kubebuilder:validation:Optional
	ServiceName string `json:"serviceName"`

	// +kubebuilder:default=2
	// +kubebuilder:validation:Optional
	WebStoreReplicas int `json:"webStoreReplicas"`
}

// WebStoreStatus defines the observed state of WebStore.
type WebStoreStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Created               bool                    `json:"created,omitempty"`
	DependenciesSatisfied bool                    `json:"dependenciesSatisfied,omitempty"`
	Conditions            []common.PhaseCondition `json:"conditions,omitempty"`
	Resources             []common.Resource       `json:"resources,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// WebStore is the Schema for the webstores API.
type WebStore struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              WebStoreSpec   `json:"spec,omitempty"`
	Status            WebStoreStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// WebStoreList contains a list of WebStore.
type WebStoreList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []WebStore `json:"items"`
}

// interface methods

// GetReadyStatus returns the ready status for a component.
func (component *WebStore) GetReadyStatus() bool {
	return component.Status.Created
}

// SetReadyStatus sets the ready status for a component.
func (component *WebStore) SetReadyStatus(status bool) {
	component.Status.Created = status
}

// GetDependencyStatus returns the dependency status for a component.
func (component *WebStore) GetDependencyStatus() bool {
	return component.Status.DependenciesSatisfied
}

// SetDependencyStatus sets the dependency status for a component.
func (component *WebStore) SetDependencyStatus(dependencyStatus bool) {
	component.Status.DependenciesSatisfied = dependencyStatus
}

// GetPhaseConditions returns the phase conditions for a component.
func (component WebStore) GetPhaseConditions() []common.PhaseCondition {
	return component.Status.Conditions
}

// SetPhaseCondition sets the phase conditions for a component.
func (component *WebStore) SetPhaseCondition(condition common.PhaseCondition) {
	if found := condition.GetPhaseConditionIndex(component); found >= 0 {
		if condition.LastModified == "" {
			condition.LastModified = time.Now().UTC().String()
		}
		component.Status.Conditions[found] = condition
	} else {
		component.Status.Conditions = append(component.Status.Conditions, condition)
	}
}

// GetResources returns the resources for a component.
func (component WebStore) GetResources() []common.Resource {
	return component.Status.Resources
}

// SetResources sets the phase conditions for a component.
func (component *WebStore) SetResource(resource common.Resource) {

	if found := resource.GetResourceIndex(component); found >= 0 {
		if resource.ResourceCondition.LastModified == "" {
			resource.ResourceCondition.LastModified = time.Now().UTC().String()
		}
		component.Status.Resources[found] = resource
	} else {
		component.Status.Resources = append(component.Status.Resources, resource)
	}
}

// GetDependencies returns the dependencies for a component.
func (*WebStore) GetDependencies() []common.Component {
	return []common.Component{}
}

// GetComponentGVK returns a GVK object for the component.
func (*WebStore) GetComponentGVK() schema.GroupVersionKind {
	return schema.GroupVersionKind{
		Group:   GroupVersion.Group,
		Version: GroupVersion.Version,
		Kind:    "WebStore",
	}
}

func init() {
	SchemeBuilder.Register(&WebStore{}, &WebStoreList{})
}
