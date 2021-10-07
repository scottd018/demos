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
	ConfigMapKind = "ConfigMap"
)

// ConfigMapIsReady performs the logic to determine if a secret is ready.
func ConfigMapIsReady(resource common.ComponentResource, expectedKeys ...string) (bool, error) {
	var configMap v1.ConfigMap
	if err := getObject(resource, &configMap, true); err != nil {
		return false, err
	}

	// if we have a name that is empty, we know we did not find the object
	if configMap.Name == "" {
		return false, nil
	}

	// check the status for a ready ca keypair secret
	for _, key := range expectedKeys {
		if string(configMap.Data[key]) == "" {
			return false, nil
		}
	}

	return true, nil
}