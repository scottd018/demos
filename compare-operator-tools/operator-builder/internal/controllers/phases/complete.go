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
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/scottd018/demos/apis/common"
)

// CompletePhase.DefaultRequeue executes checking for a parent components readiness status.
func (phase *CompletePhase) DefaultRequeue() ctrl.Result {
	return Requeue()
}

// CompletePhase.Execute executes the completion of a reconciliation loop.
func (phase *CompletePhase) Execute(
	r common.ComponentReconciler,
) (proceedToNextPhase bool, err error) {
	r.GetComponent().SetReadyStatus(true)
	r.GetLogger().V(0).Info("successfully reconciled")

	return true, nil
}
