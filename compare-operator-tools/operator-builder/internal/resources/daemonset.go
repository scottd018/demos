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
	DaemonSetKind = "DaemonSet"
)

// DaemonSetIsReady checks to see if a daemonset is ready.
func DaemonSetIsReady(resource common.ComponentResource) (bool, error) {
	var daemonSet appsv1.DaemonSet
	if err := getObject(resource, &daemonSet, true); err != nil {
		return false, err
	}

	// if we have a name that is empty, we know we did not find the object
	if daemonSet.Name == "" {
		return false, nil
	}

	// ensure the desired number is scheduled and ready
	if daemonSet.Status.DesiredNumberScheduled == daemonSet.Status.NumberReady {
		if daemonSet.Status.NumberReady > 0 && daemonSet.Status.NumberUnavailable < 1 {
			return true, nil
		} else {
			return false, nil
		}
	}

	return false, nil
}
