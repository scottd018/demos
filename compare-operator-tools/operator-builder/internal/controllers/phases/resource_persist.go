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
	"time"

	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/scottd018/demos/apis/common"
)

// PersistResourcePhase.Execute executes persisting resources to the Kubernetes database.
func (phase *PersistResourcePhase) Execute(
	resource common.ComponentResource,
	resourceCondition common.ResourceCondition,
) (ctrl.Result, bool, error) {
	// persist the resource
	if err := persistResource(
		resource,
		resourceCondition,
		phase,
	); err != nil {
		return ctrl.Result{}, false, err
	}

	return ctrl.Result{}, true, nil
}

// persistResource persists a single resource to the Kubernetes database.
func persistResource(
	resource common.ComponentResource,
	condition common.ResourceCondition,
	phase *PersistResourcePhase,
) error {
	// persist resource
	r := resource.GetReconciler()
	if err := r.CreateOrUpdate(resource.GetObject()); err != nil {
		if IsOptimisticLockError(err) {
			return nil
		} else {
			r.GetLogger().V(0).Info(err.Error())

			return err
		}
	}

	// set attributes related to the persistence of this child resource
	condition.LastResourcePhase = getResourcePhaseName(phase)
	condition.LastModified = time.Now().UTC().String()
	condition.Message = "resource created successfully"
	condition.Created = true

	// update the condition to notify that we have created a child resource
	return updateResourceConditions(r, *resource.ToCommonResource(), &condition)
}
