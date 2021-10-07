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
	DeploymentKind = "Deployment"
)

// DeploymentIsReady performs the logic to determine if a deployment is ready.
func DeploymentIsReady(resource common.ComponentResource) (bool, error) {
	var deployment appsv1.Deployment
	if err := getObject(resource, &deployment, true); err != nil {
		return false, err
	}

	// if we have a name that is empty, we know we did not find the object
	if deployment.Name == "" {
		return false, nil
	}

	// rely on observed generation to give us a proper status
	if deployment.Generation != deployment.Status.ObservedGeneration {
		return false, nil
	}

	// check the status for a ready deployment
	if deployment.Status.ReadyReplicas != deployment.Status.Replicas {
		return false, nil
	}

	return true, nil
}
