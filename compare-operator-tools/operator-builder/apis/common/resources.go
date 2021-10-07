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

package common

// ResourceCommon are the common fields used across multiple resource types.
type ResourceCommon struct {
	// Group defines the API Group of the resource.
	Group string `json:"group"`

	// Version defines the API Version of the resource.
	Version string `json:"version"`

	// Kind defines the kind of the resource.
	Kind string `json:"kind"`

	// Name defines the name of the resource from the metadata.name field.
	Name string `json:"name"`

	// Namespace defines the namespace in which this resource exists in.
	Namespace string `json:"namespace"`
}

// Resource is the resource and its condition as stored on the object status field.
type Resource struct {
	ResourceCommon `json:",omitempty"`

	// ResourceCondition defines the current condition of this resource.
	ResourceCondition `json:"condition,omitempty"`
}

// GetResourceIndex returns the index of a matching resource.  Any integer which is 0
// or greater indicates that the resource was found.  Anything lower indicates that an
// associated resource is not found.
func (resource *Resource) GetResourceIndex(component Component) int {
	for i, currentResource := range component.GetResources() {
		if currentResource.Group == resource.Group && currentResource.Version == resource.Version && currentResource.Kind == resource.Kind {
			if currentResource.Name == resource.Name && currentResource.Namespace == resource.Namespace {
				return i
			}
		}
	}

	return -1
}
