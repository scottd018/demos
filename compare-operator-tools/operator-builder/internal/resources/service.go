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
	corev1 "k8s.io/api/core/v1"

	"github.com/scottd018/demos/apis/common"
)

const (
	ServiceKind = "Service"
)

// ServiceIsReady checks to see if a job is ready.
func ServiceIsReady(resource common.ComponentResource) (bool, error) {
	var service corev1.Service
	if err := getObject(resource, &service, true); err != nil {
		return false, err
	}

	// if we have a name that is empty, we know we did not find the object
	if service.Name == "" {
		return false, nil
	}

	// return if we have an external service type
	if service.Spec.Type == corev1.ServiceTypeExternalName {
		return true, nil
	}

	// ensure a cluster ip address exists for cluster ip types
	if service.Spec.ClusterIP != corev1.ClusterIPNone && len(service.Spec.ClusterIP) == 0 {
		return false, nil
	}

	// ensure a load balancer ip or hostname is present
	if service.Spec.Type == corev1.ServiceTypeLoadBalancer {
		if len(service.Status.LoadBalancer.Ingress) == 0 {
			return false, nil
		}
	}

	return true, nil
}
