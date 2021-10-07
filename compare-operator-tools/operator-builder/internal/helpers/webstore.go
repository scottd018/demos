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

package helpers

import (
	"fmt"

	appsv1alpha1 "github.com/scottd018/demos/apis/apps/v1alpha1"
	common "github.com/scottd018/demos/apis/common"
)

// WebStoreUnique returns only one WebStore and returns an error if more than one are found.
func WebStoreUnique(
	reconciler common.ComponentReconciler,
) (
	*appsv1alpha1.WebStore,
	error,
) {
	components, err := WebStoreList(reconciler)
	if err != nil {
		return nil, err
	}

	if len(components.Items) != 1 {
		return nil, fmt.Errorf("expected only 1 WebStore; found %v\n", len(components.Items))
	}

	component := components.Items[0]

	return &component, nil
}

// WebStoreList gets a WebStoreList from the cluster.
func WebStoreList(
	reconciler common.ComponentReconciler,
) (
	*appsv1alpha1.WebStoreList,
	error,
) {
	components := &appsv1alpha1.WebStoreList{}
	if err := reconciler.List(reconciler.GetContext(), components); err != nil {
		reconciler.GetLogger().V(0).Info("unable to retrieve WebStoreList from cluster")

		return nil, err
	}

	return components, nil
}
