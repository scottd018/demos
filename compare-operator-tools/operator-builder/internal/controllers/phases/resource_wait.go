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

package phases

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/scottd018/demos/apis/common"
	"github.com/scottd018/demos/internal/resources"
)

// WaitForResourcePhase.Execute executes waiting for a resource to be ready before continuing.
func (phase *WaitForResourcePhase) Execute(
	resource common.ComponentResource,
	resourceCondition common.ResourceCondition,
) (ctrl.Result, bool, error) {
	// TODO: loop through functions instead of repeating logic
	// common wait logic for a resource
	ready, err := commonWait(resource.GetReconciler(), resource)
	if err != nil {
		return ctrl.Result{}, false, err
	}

	// return the result if the object is not ready
	if !ready {
		return Requeue(), false, nil
	}

	// specific wait logic for a resource
	meta := resource.GetObject().(metav1.Object)
	ready, err = resource.GetReconciler().Wait(&meta)
	if err != nil {
		return ctrl.Result{}, false, err
	}

	// return the result if the object is not ready
	if !ready {
		return Requeue(), false, nil
	}

	return ctrl.Result{}, true, nil
}

// commonWait applies all common waiting functions for known resources.
func commonWait(
	r common.ComponentReconciler,
	resource common.ComponentResource,
) (bool, error) {
	// Namespace
	if resource.GetObject().GetNamespace() != "" {
		return resources.NamespaceForResourceIsReady(resource)
	}

	return true, nil
}
