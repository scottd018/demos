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

package mutate

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/scottd018/demos/apis/common"
)

// WebStoreMutate performs the logic to mutate resources that belong to the parent.
func WebStoreMutate(reconciler common.ComponentReconciler,
	object *metav1.Object,
) (replacedObjects []metav1.Object, skip bool, err error) {
	return []metav1.Object{*object}, false, nil
}
