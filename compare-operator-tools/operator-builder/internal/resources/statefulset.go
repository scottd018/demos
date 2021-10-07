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
	appsv1 "k8s.io/api/apps/v1"

	"github.com/scottd018/demos/apis/common"
)

const (
	StatefulSetKind = "StatefulSet"
)

// StatefulSetIsReady performs the logic to determine if a secret is ready.
func StatefulSetIsReady(resource common.ComponentResource, expectedKeys ...string) (bool, error) {
	var statefulSet appsv1.StatefulSet
	if err := getObject(resource, &statefulSet, true); err != nil {
		return false, err
	}

	// if we have a name that is empty, we know we did not find the object
	if statefulSet.Name == "" {
		return false, nil
	}

	// rely on observed generation to give us a proper status
	if statefulSet.Generation != statefulSet.Status.ObservedGeneration {
		return false, nil
	}

	// check for valid replicas
	replicas := statefulSet.Spec.Replicas
	if replicas == nil {
		return false, nil
	}

	// check to see if replicas have been updated
	var needsUpdate int32
	if statefulSet.Spec.UpdateStrategy.RollingUpdate != nil &&
		statefulSet.Spec.UpdateStrategy.RollingUpdate.Partition != nil &&
		*statefulSet.Spec.UpdateStrategy.RollingUpdate.Partition > 0 {

		needsUpdate -= *statefulSet.Spec.UpdateStrategy.RollingUpdate.Partition
	}
	notUpdated := needsUpdate - statefulSet.Status.UpdatedReplicas
	if notUpdated > 0 {
		return false, nil
	}

	// check to see if replicas are available
	notReady := *replicas - statefulSet.Status.ReadyReplicas
	if notReady > 0 {
		return false, nil
	}

	// check to see if a scale down operation is complete
	notDeleted := statefulSet.Status.Replicas - *replicas
	if notDeleted > 0 {
		return false, nil
	}

	return true, nil
}
