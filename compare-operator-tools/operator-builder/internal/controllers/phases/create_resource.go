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
	"fmt"

	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/scottd018/demos/apis/common"
)

// CreateResourcesPhase.DefaultRequeue executes checking for a parent components readiness status.
func (phase *CreateResourcesPhase) DefaultRequeue() ctrl.Result {
	return Requeue()
}

// createResourcePhases defines the phases for resource creation and the order in which they run during the reconcile process.
func createResourcePhases() []ResourcePhase {
	return []ResourcePhase{
		// wait for other resources before attempting to create
		&WaitForResourcePhase{},

		// create the resource in the cluster
		&PersistResourcePhase{},
	}
}

// CreateResourcesPhase.Execute executes executes sub-phases which are required to create the resources.
func (phase *CreateResourcesPhase) Execute(
	r common.ComponentReconciler,
) (proceedToNextPhase bool, err error) {
	// execute the resource phases against each resource
	for _, resource := range r.GetResources() {
		resourceCommon := resource.ToCommonResource()
		resourceCondition := &common.ResourceCondition{}

		for _, resourcePhase := range createResourcePhases() {
			r.GetLogger().V(7).Info(fmt.Sprintf("enter resource phase: %T", resourcePhase))
			_, proceed, err := resourcePhase.Execute(resource, *resourceCondition)

			// set a message, return the error and result on error or when unable to proceed
			if err != nil || !proceed {
				return handleResourcePhaseExit(r, *resourceCommon, *resourceCondition, resourcePhase, proceed, err)
			}

			// set attributes on the resource condition before updating the status
			resourceCondition.LastResourcePhase = getResourcePhaseName(resourcePhase)

			r.GetLogger().V(5).Info(fmt.Sprintf("completed resource phase: %T", resourcePhase))
		}
	}

	return true, nil
}
