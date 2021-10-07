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
	"k8s.io/apimachinery/pkg/api/errors"

	extensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"

	"github.com/scottd018/demos/apis/common"
)

const (
	CustomResourceDefinitionKind = "CustomResourceDefinition"
)

// CustomResourceDefinitionIsReady performs the logic to determine if a custom resource definition is ready.
func CustomResourceDefinitionIsReady(resource common.ComponentResource) (bool, error) {
	var crd extensionsv1.CustomResourceDefinition
	if err := getObject(resource, &crd, false); err != nil {
		if errors.IsNotFound(err) {
			return false, nil
		}
	}

	return true, nil
}
