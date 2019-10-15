package kappnav

import (
	"context"

	kappnavv1 "github.com/kappnav/operator/pkg/apis/kappnav/v1"
	kappnavutils "github.com/kappnav/operator/pkg/utils"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_kappnav")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new Kappnav Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileKappnav{ReconcilerBase: kappnavutils.NewReconcilerBase(mgr.GetClient(), 
		mgr.GetScheme(), mgr.GetConfig(), mgr.GetRecorder("kappnav-operator"))}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("kappnav-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource Kappnav
	err = c.Watch(&source.Kind{Type: &kappnavv1.Kappnav{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Deployment and requeue the owner Kappnav
	err = c.Watch(&source.Kind{Type: &appsv1.Deployment{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &kappnavv1.Kappnav{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileKappnav implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileKappnav{}

// ReconcileKappnav reconciles a Kappnav object
type ReconcileKappnav struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	kappnavutils.ReconcilerBase
}

// Reconcile reads that state of the cluster for a Kappnav object and makes changes based on the state read
// and what is in the Kappnav.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileKappnav) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling Kappnav")

	// Fetch the Kappnav instance
	instance := &kappnavv1.Kappnav{}
	err := r.GetClient().Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}
	// Apply defaults to the Kappnav instance
	kappnavutils.SetKappnavDefaults(instance)

	// Create or update the UI deployment
	uiDeployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: instance.GetName() + "-ui",
			Namespace: instance.GetNamespace(),
		},
	}
	err = r.CreateOrUpdate(uiDeployment, instance, func() error {
		kappnavutils.CustomizeDeployment(uiDeployment, instance)
		kappnavutils.CustomizePodSpec(&uiDeployment.Spec.Template, 
			kappnavutils.CreateUIDeploymentContainers(instance),
			kappnavutils.CreateUIVolumes(), instance)
		return nil
	})
	if err != nil {
		reqLogger.Error(err, "Failed to reconcile the UI Deployment")
		return r.ManageError(err, kappnavv1.StatusConditionTypeReconciled, instance)
	}

	// Create or update the Controller deployment
	controllerDeployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: instance.GetName() + "-controller",
			Namespace: instance.GetNamespace(),
		},
	}
	err = r.CreateOrUpdate(controllerDeployment, instance, func() error {
		kappnavutils.CustomizeDeployment(controllerDeployment, instance)
		kappnavutils.CustomizePodSpec(&controllerDeployment.Spec.Template, 
			kappnavutils.CreateControllerDeploymentContainers(instance), nil, instance)
		return nil
	})
	if err != nil {
		reqLogger.Error(err, "Failed to reconcile the Controller Deployment")
		return r.ManageError(err, kappnavv1.StatusConditionTypeReconciled, instance)
	}

	return r.ManageSuccess(kappnavv1.StatusConditionTypeReconciled, instance)
}
