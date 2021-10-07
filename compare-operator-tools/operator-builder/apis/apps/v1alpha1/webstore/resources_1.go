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

package webstore

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	appsv1alpha1 "github.com/scottd018/demos/apis/apps/v1alpha1"
)

// CreateDeploymentWebstoreDeploy creates the webstore-deploy Deployment resource.
func CreateDeploymentWebstoreDeploy(
	parent *appsv1alpha1.WebStore) (metav1.Object, error) {
	var resourceObj = &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "apps/v1",
			"kind":       "Deployment",
			"metadata": map[string]interface{}{
				"name": "webstore-deploy",
			},
			"spec": map[string]interface{}{
				"replicas": parent.Spec.WebStoreReplicas,
				"selector": map[string]interface{}{
					"matchLabels": map[string]interface{}{
						"app": "webstore",
					},
				},
				"template": map[string]interface{}{
					"metadata": map[string]interface{}{
						"labels": map[string]interface{}{
							"app": "webstore",
						},
					},
					"spec": map[string]interface{}{
						"containers": []interface{}{
							map[string]interface{}{
								"name": "webstore-container",
								// Defines the web store image, controlled by webstoreImage
								"image": parent.Spec.WebstoreImage,
								"ports": []interface{}{
									map[string]interface{}{
										"containerPort": 8080,
									},
								},
							},
						},
					},
				},
			},
		},
	}

	resourceObj.SetNamespace(parent.Namespace)

	return resourceObj, nil
}

// CreateIngressWebstoreIng creates the webstore-ing Ingress resource.
func CreateIngressWebstoreIng(
	parent *appsv1alpha1.WebStore) (metav1.Object, error) {
	var resourceObj = &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "networking.k8s.io/v1beta1",
			"kind":       "Ingress",
			"metadata": map[string]interface{}{
				"name": "webstore-ing",
				"annotations": map[string]interface{}{
					"nginx.ingress.kubernetes.io/rewrite-target": "/",
				},
			},
			"spec": map[string]interface{}{
				"rules": []interface{}{
					map[string]interface{}{
						"host": "app.acme.com",
						"http": map[string]interface{}{
							"paths": []interface{}{
								map[string]interface{}{
									"path": "/",
									"backend": map[string]interface{}{
										"serviceName": "webstorep-svc",
										"servicePort": 80,
									},
								},
							},
						},
					},
				},
			},
		},
	}

	resourceObj.SetNamespace(parent.Namespace)

	return resourceObj, nil
}

// CreateServiceParentSpecServiceName creates the parent.Spec.ServiceName Service resource.
func CreateServiceParentSpecServiceName(
	parent *appsv1alpha1.WebStore) (metav1.Object, error) {
	var resourceObj = &unstructured.Unstructured{
		Object: map[string]interface{}{
			"kind":       "Service",
			"apiVersion": "v1",
			"metadata": map[string]interface{}{
				"name": parent.Spec.ServiceName,
			},
			"spec": map[string]interface{}{
				"selector": map[string]interface{}{
					"app": "webstore",
				},
				"ports": []interface{}{
					map[string]interface{}{
						"protocol":   "TCP",
						"port":       80,
						"targetPort": 8080,
					},
				},
			},
		},
	}

	resourceObj.SetNamespace(parent.Namespace)

	return resourceObj, nil
}
