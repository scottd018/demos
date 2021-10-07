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

package resources

import (
	v1 "k8s.io/api/core/v1"

	"github.com/scottd018/demos/apis/common"
)

const (
	NamespaceKind = "Namespace"
)

// NamespaceIsReady defines the criteria for a namespace to be condsidered ready.
func NamespaceIsReady(resource common.ComponentResource) (bool, error) {
	var namespace v1.Namespace
	if err := getObject(resource, &namespace, true); err != nil {
		return false, err
	}

	// if we have a name that is empty, we know we did not find the object
	if namespace.Name == "" {
		return false, nil
	}

	// if the namespace is terminating, it is not considered ready
	if namespace.Status.Phase == v1.NamespaceTerminating {
		return false, nil
	}

	// finally, rely on the active field to determine if this namespace is ready
	return namespace.Status.Phase == v1.NamespaceActive, nil
}

// NamespaceForResourceIsReady checks to see if the namespace of a resource is ready.
func NamespaceForResourceIsReady(resource common.ComponentResource) (bool, error) {
	// create a stub namespace resource to pass to the NamespaceIsReady method
	namespace := &Resource{
		Reconciler: resource.GetReconciler(),
	}

	// insert the inherited fields
	namespace.Name = resource.GetNamespace()
	namespace.Group = ""
	namespace.Version = "v1"
	namespace.Kind = NamespaceKind

	return NamespaceIsReady(namespace)
}
