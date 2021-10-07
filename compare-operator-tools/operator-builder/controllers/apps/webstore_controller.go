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

package apps

import (
	"context"
	"fmt"
	"time"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"

	appsv1alpha1 "github.com/scottd018/demos/apis/apps/v1alpha1"
	"github.com/scottd018/demos/apis/apps/v1alpha1/webstore"
	"github.com/scottd018/demos/apis/common"
	"github.com/scottd018/demos/internal/controllers/phases"
	"github.com/scottd018/demos/internal/controllers/utils"
	"github.com/scottd018/demos/internal/dependencies"
	"github.com/scottd018/demos/internal/mutate"
	"github.com/scottd018/demos/internal/resources"
	"github.com/scottd018/demos/internal/wait"
)

// WebStoreReconciler reconciles a WebStore object.
type WebStoreReconciler struct {
	client.Client
	Name       string
	Log        logr.Logger
	Scheme     *runtime.Scheme
	Context    context.Context
	Controller controller.Controller
	Watches    []client.Object
	Resources  []common.ComponentResource
	Component  *appsv1alpha1.WebStore
}

// +kubebuilder:rbac:groups=apps.acme.com,resources=webstores,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps.acme.com,resources=webstores/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=networking.k8s.io,resources=ingresses,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the WebApp object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.2/pkg/reconcile
func (r *WebStoreReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	r.Context = ctx
	log := r.Log.WithValues("webstore", req.NamespacedName)

	// get and store the component
	r.Component = &appsv1alpha1.WebStore{}
	if err := r.Get(r.Context, req.NamespacedName, r.Component); err != nil {
		log.V(0).Info("unable to fetch WebStore")

		return ctrl.Result{}, utils.IgnoreNotFound(err)
	}

	// get and store the resources
	if err := r.SetResources(); err != nil {
		return ctrl.Result{}, err
	}

	// execute the phases
	for _, phase := range utils.Phases(r.Component) {
		r.GetLogger().V(7).Info(fmt.Sprintf("enter phase: %T", phase))
		proceed, err := phase.Execute(r)
		result, err := phases.HandlePhaseExit(r, phase, proceed, err)

		// return only if we have an error or are told not to proceed
		if err != nil || !proceed {
			log.V(2).Info(fmt.Sprintf("not ready; requeuing phase: %T", phase))

			return result, err
		}

		r.GetLogger().V(5).Info(fmt.Sprintf("completed phase: %T", phase))
	}

	return phases.DefaultReconcileResult(), nil
}

// Construct resources runs the methods to properly construct the resources.
func (r *WebStoreReconciler) ConstructResources() ([]metav1.Object, error) {

	resourceObjects := make([]metav1.Object, len(webstore.CreateFuncs))

	// create resources in memory
	for i, f := range webstore.CreateFuncs {
		resource, err := f(r.Component)
		if err != nil {
			return nil, err
		}

		resourceObjects[i] = resource
	}

	return resourceObjects, nil
}

// GetResources will return the resources associated with the reconciler.
func (r *WebStoreReconciler) GetResources() []common.ComponentResource {
	return r.Resources
}

// SetResources will create and return the resources in memory.
func (r *WebStoreReconciler) SetResources() error {
	// create resources in memory
	baseResources, err := r.ConstructResources()
	if err != nil {
		return err
	}

	// loop through the in memory resources and store them on the reconciler
	for _, base := range baseResources {
		// run through the mutation functions to mutate the resources
		mutatedResources, skip, err := r.Mutate(&base)
		if err != nil {
			return err
		}
		if skip {
			continue
		}

		for _, mutated := range mutatedResources {
			resourceObject := resources.NewResourceFromClient(mutated.(client.Object))
			resourceObject.Reconciler = r

			r.SetResource(resourceObject)
		}
	}

	return nil
}

// SetResource will set a resource on the objects if the relevant object does not already exist.
func (r *WebStoreReconciler) SetResource(new common.ComponentResource) {

	// set and return immediately if nothing exists
	if len(r.Resources) == 0 {
		r.Resources = append(r.Resources, new)

		return
	}

	// loop through the resources and set or update when found
	for i, existing := range r.Resources {
		if new.EqualGVK(existing) && new.EqualNamespaceName(existing) {
			r.Resources[i] = new

			return
		}
	}

	// if we haven't returned yet, we have not found the resource and must add it
	r.Resources = append(r.Resources, new)
}

// CreateOrUpdate creates a resource if it does not already exist or updates a resource
// if it does already exist.
func (r *WebStoreReconciler) CreateOrUpdate(
	resource metav1.Object,
) error {
	// set ownership on the underlying resource being created or updated
	if err := ctrl.SetControllerReference(r.Component, resource, r.Scheme); err != nil {
		r.GetLogger().V(0).Info("unable to set owner reference on resource")

		return err
	}

	// create a stub object to store the current resource in the cluster so that we do not affect
	// the desired state of the resource object in memory
	newResource := resources.NewResourceFromClient(resource.(client.Object), r)
	resourceStub := &unstructured.Unstructured{}
	resourceStub.SetGroupVersionKind(newResource.Object.GetObjectKind().GroupVersionKind())
	oldResource := resources.NewResourceFromClient(resourceStub, r)

	if err := r.Get(
		r.Context,
		client.ObjectKeyFromObject(newResource.Object),
		oldResource.Object,
	); err != nil {
		// create the resource if we cannot find one
		if errors.IsNotFound(err) {
			if err := newResource.Create(); err != nil {
				return err
			}
		} else {
			return err
		}
	} else {
		// update the resource
		if err := newResource.Update(oldResource); err != nil {
			return err
		}
	}

	return utils.Watch(r, newResource.Object)
}

// GetLogger returns the logger from the reconciler.
func (r *WebStoreReconciler) GetLogger() logr.Logger {
	return r.Log
}

// GetClient returns the client from the reconciler.
func (r *WebStoreReconciler) GetClient() client.Client {
	return r.Client
}

// GetScheme returns the scheme from the reconciler.
func (r *WebStoreReconciler) GetScheme() *runtime.Scheme {
	return r.Scheme
}

// GetContext returns the context from the reconciler.
func (r *WebStoreReconciler) GetContext() context.Context {
	return r.Context
}

// GetName returns the name of the reconciler.
func (r *WebStoreReconciler) GetName() string {
	return r.Name
}

// GetComponent returns the component the reconciler is operating against.
func (r *WebStoreReconciler) GetComponent() common.Component {
	return r.Component
}

// GetController returns the controller object associated with the reconciler.
func (r *WebStoreReconciler) GetController() controller.Controller {
	return r.Controller
}

// GetWatches returns the objects which are current being watched by the reconciler.
func (r *WebStoreReconciler) GetWatches() []client.Object {
	return r.Watches
}

// SetWatch appends a watch to the list of currently watched objects.
func (r *WebStoreReconciler) SetWatch(watch client.Object) {
	r.Watches = append(r.Watches, watch)
}

// UpdateStatus updates the status for a component.
func (r *WebStoreReconciler) UpdateStatus() error {
	return r.Status().Update(r.Context, r.Component)
}

// CheckReady will return whether a component is ready.
func (r *WebStoreReconciler) CheckReady() (bool, error) {
	return dependencies.WebStoreCheckReady(r)
}

// Mutate will run the mutate phase of a resource.
func (r *WebStoreReconciler) Mutate(
	object *metav1.Object,
) ([]metav1.Object, bool, error) {
	return mutate.WebStoreMutate(r, object)
}

// Wait will run the wait phase of a resource.
func (r *WebStoreReconciler) Wait(
	object *metav1.Object,
) (bool, error) {
	return wait.WebStoreWait(r, object)
}

func (r *WebStoreReconciler) SetupWithManager(mgr ctrl.Manager) error {
	options := controller.Options{
		RateLimiter: utils.NewDefaultRateLimiter(5*time.Microsecond, 5*time.Minute),
	}

	baseController, err := ctrl.NewControllerManagedBy(mgr).
		WithOptions(options).
		WithEventFilter(utils.ComponentPredicates()).
		For(&appsv1alpha1.WebStore{}).
		Build(r)
	if err != nil {
		return err
	}

	r.Controller = baseController

	return nil
}
