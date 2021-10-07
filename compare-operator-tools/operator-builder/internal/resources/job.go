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
	"fmt"

	batchv1 "k8s.io/api/batch/v1"

	"github.com/scottd018/demos/apis/common"
)

const (
	JobKind = "Job"
)

// JobIsReady checks to see if a job is ready.
func JobIsReady(resource common.ComponentResource) (bool, error) {
	var job batchv1.Job
	if err := getObject(resource, &job, true); err != nil {
		return false, err
	}

	// if we have a name that is empty, we know we did not find the object
	if job.Name == "" {
		return false, nil
	}

	// return immediately if the job is active or has no completion time
	if job.Status.Active == 1 || job.Status.CompletionTime == nil {
		return false, nil
	}

	// ensure the completion is actually successful
	if job.Status.Succeeded != 1 {
		return false, fmt.Errorf("job " + job.GetName() + " was not successful")
	}

	return true, nil
}
