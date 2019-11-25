/*
Copyright 2019 IBM Corporation
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

package kappnav

import (
	"bytes"
	"context"
	"io/ioutil"
	"os"
	"strings"
	"text/template"

	kappnavv1 "github.com/kappnav/operator/pkg/apis/kappnav/v1"
	kappnavutils "github.com/kappnav/operator/pkg/utils"
	routev1 "github.com/openshift/api/route/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	rbacv1 "k8s.io/api/rbac/v1"
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
	"sigs.k8s.io/yaml"
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
	reconciler := &ReconcileKappnav{ReconcilerBase: kappnavutils.NewReconcilerBase(mgr.GetClient(),
		mgr.GetScheme(), mgr.GetConfig(), mgr.GetRecorder("kappnav-operator"))}

	// Create CRDs if they do not already exist.
	files, err := ioutil.ReadDir("crds")
	if err != nil {
		log.Error(err, "Failed to read directory: crds")
		os.Exit(1)
	}
	for _, file := range files {
		if !file.IsDir() {
			fileName := "crds/" + file.Name()
			if strings.HasSuffix(fileName, ".yaml") || strings.HasSuffix(fileName, ".yml") {
				// Read the file from the image.
				fData, err := ioutil.ReadFile(fileName)
				if err != nil {
					log.Error(err, "Failed to read file: "+fileName)
					os.Exit(1)
				}
				crd := &apiextensionsv1beta1.CustomResourceDefinition{}
				// Unmarshal the YAML into an object.
				err = yaml.Unmarshal(fData, crd)
				if err != nil {
					log.Error(err, "Failed to unmarshal YAML file: "+fileName)
					os.Exit(1)
				}
				// Create the CRD if it does not already exist.
				err = reconciler.GetClient().Create(context.TODO(), crd)
				if err != nil && !errors.IsAlreadyExists(err) {
					log.Error(err, "Failed to create CRD: "+crd.GetName())
					os.Exit(1)
				}
			}
		}
	}

	return reconciler
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

	// Watch for changes to secondary resource Deployment and requeue the owner Kappnav
	err = c.Watch(&source.Kind{Type: &appsv1.Deployment{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &kappnavv1.Kappnav{},
	})
	if err != nil {
		return err
	}

	// Watch for changes to secondary resource ConfigMap and requeue the owner Kappnav
	err = c.Watch(&source.Kind{Type: &corev1.ConfigMap{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &kappnavv1.Kappnav{},
	})
	if err != nil {
		return err
	}

	// Watch for changes to secondary resource Route (when available) and requeue the owner Kappnav
	_ = c.Watch(&source.Kind{Type: &routev1.Route{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &kappnavv1.Kappnav{},
	})
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

	// Call factory method to create new KappnavExtension
	extension := kappnavutils.NewKappnavExtension()

	// Apply defaults to the Kappnav instance
	err = kappnavutils.SetKappnavDefaults(instance)
	if err != nil {
		reqLogger.Error(err, "Failed to process default values file")
		return reconcile.Result{}, err
	}

	uiServiceAndRouteName := &metav1.ObjectMeta{
		Name:      instance.GetName() + "-ui-service",
		Namespace: instance.GetNamespace(),
	}

	// Create or update service account
	serviceAccount := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.GetName() + "-" + kappnavutils.ServiceAccountNameSuffix,
			Namespace: instance.GetNamespace(),
		},
	}
	err = r.CreateOrUpdate(serviceAccount, instance, func() error {
		kappnavutils.CustomizeServiceAccount(serviceAccount, uiServiceAndRouteName, instance)
		return nil
	})
	if err != nil {
		reqLogger.Error(err, "Failed to reconcile the ServiceAccount")
		return r.ManageError(err, kappnavv1.StatusConditionTypeReconciled, instance)
	}

	// Create or update cluster role binding
	crb := &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.GetName() + "-" + instance.GetNamespace() + "-crb",
			Namespace: instance.GetNamespace(),
		},
	}
	err = r.CreateOrUpdate(crb, instance, func() error {
		kappnavutils.CustomizeClusterRoleBinding(crb, serviceAccount, instance)
		return nil
	})
	if err != nil && !errors.IsAlreadyExists(err) {
		reqLogger.Error(err, "Failed to reconcile the ClusterRoleBinding")
		return r.ManageError(err, kappnavv1.StatusConditionTypeReconciled, instance)
	}

	// Dummy secret for Minikube support
	dummySecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.GetName() + "-" + kappnavutils.OAuthVolumeName,
			Namespace: instance.GetNamespace(),
		},
	}

	// The UI service
	uiService := &corev1.Service{
		ObjectMeta: *uiServiceAndRouteName,
	}
	uiServiceAnnotations := map[string]string{
		"service.alpha.openshift.io/serving-cert-secret-name": dummySecret.Name,
	}

	// Kappnav URL is computed from the route
	kappnavURL := ""

	isMinikube := kappnavutils.IsMinikubeEnv(instance.Spec.Env.KubeEnv)
	if isMinikube {
		// Create or update dummy secret
		err = r.CreateOrUpdate(dummySecret, instance, func() error {
			kappnavutils.CustomizeSecret(dummySecret, instance)
			return nil
		})
		if err != nil {
			reqLogger.Error(err, "Failed to reconcile the dummy secret")
			return r.ManageError(err, kappnavv1.StatusConditionTypeReconciled, instance)
		}
		// Create or update the UI service
		err = r.CreateOrUpdate(uiService, instance, func() error {
			kappnavutils.CustomizeService(uiService, instance, uiServiceAnnotations)
			kappnavutils.CustomizeUIServiceSpec(&uiService.Spec, instance)
			return nil
		})
		if err != nil {
			reqLogger.Error(err, "Failed to reconcile the UI service")
			return r.ManageError(err, kappnavv1.StatusConditionTypeReconciled, instance)
		}
		// Create or update UI ingress
		uiIngress := &extensionsv1beta1.Ingress{
			ObjectMeta: metav1.ObjectMeta{
				Name:      instance.GetName() + "-ui-ingress",
				Namespace: instance.GetNamespace(),
			},
		}
		err = r.CreateOrUpdate(uiIngress, instance, func() error {
			kappnavutils.CustomizeIngress(uiIngress, instance)
			kappnavutils.CustomizeUIIngressSpec(&uiIngress.Spec, uiService, instance)
			return nil
		})
		if err != nil {
			reqLogger.Error(err, "Failed to reconcile the UI ingress")
			return r.ManageError(err, kappnavv1.StatusConditionTypeReconciled, instance)
		}
	} else {
		// Create or update the UI service
		err = r.CreateOrUpdate(uiService, instance, func() error {
			kappnavutils.CustomizeService(uiService, instance, uiServiceAnnotations)
			kappnavutils.CustomizeUIServiceSpec(&uiService.Spec, instance)
			return nil
		})
		if err != nil {
			reqLogger.Error(err, "Failed to reconcile the UI service")
			return r.ManageError(err, kappnavv1.StatusConditionTypeReconciled, instance)
		}
		// Create or update UI route
		uiRoute := &routev1.Route{
			ObjectMeta: *uiServiceAndRouteName,
		}
		err = r.CreateOrUpdate(uiRoute, instance, func() error {
			kappnavutils.CustomizeRoute(uiRoute, instance)
			kappnavutils.CustomizeUIRouteSpec(&uiRoute.Spec, uiServiceAndRouteName, instance)
			return nil
		})
		if err != nil {
			reqLogger.Error(err, "Failed to reconcile the UI route")
			return r.ManageError(err, kappnavv1.StatusConditionTypeReconciled, instance)
		}
		// Compute Kappnav URL from route.
		routeHost := uiRoute.Spec.Host
		routePath := uiRoute.Spec.Path
		if len(routeHost) > 0 && len(routePath) > 0 {
			kappnavURL = "https://" + routeHost + routePath
		}
	}

	// Create or update action, section and status config maps.
	mapDirs := []string{"maps/action", "maps/sections", "maps/status"}
	for _, dir := range mapDirs {
		files, err := ioutil.ReadDir(dir)
		if err != nil {
			reqLogger.Error(err, "Failed to read directory: "+dir)
			return r.ManageError(err, kappnavv1.StatusConditionTypeReconciled, instance)
		}
		for _, file := range files {
			if !file.IsDir() {
				fileName := dir + "/" + file.Name()
				if strings.HasSuffix(fileName, ".yaml") || strings.HasSuffix(fileName, ".yml") {
					// Read the file from the image.
					fData, err := ioutil.ReadFile(fileName)
					if err != nil {
						reqLogger.Error(err, "Failed to read file: "+fileName)
						return r.ManageError(err, kappnavv1.StatusConditionTypeReconciled, instance)
					}
					// Parse the file into a template.
					t, err := template.New(fileName).Parse(string(fData))
					if err != nil {
						reqLogger.Error(err, "Failed to parse template file: "+fileName)
						return r.ManageError(err, kappnavv1.StatusConditionTypeReconciled, instance)
					}
					// Execute the template against the Kappnav CR instance.
					var buf bytes.Buffer
					err = t.Execute(&buf, instance)
					if err != nil {
						reqLogger.Error(err, "Failed to execute template: "+fileName)
						return r.ManageError(err, kappnavv1.StatusConditionTypeReconciled, instance)
					}
					configMap := &corev1.ConfigMap{}
					// Unmarshal the YAML into an object.
					err = yaml.Unmarshal(buf.Bytes(), configMap)
					if err != nil {
						reqLogger.Error(err, "Failed to unmarshal YAML file: "+fileName)
						return r.ManageError(err, kappnavv1.StatusConditionTypeReconciled, instance)
					}
					clusterMap := &corev1.ConfigMap{
						ObjectMeta: metav1.ObjectMeta{
							Name:      configMap.GetName(),
							Namespace: instance.GetNamespace(),
						},
					}
					// Write the data to the map in the cluster.
					err = r.CreateOrUpdate(clusterMap, instance, func() error {
						kappnavutils.CustomizeConfigMap(clusterMap, instance)
						// Write the data section if it doesn't exist or is empty.
						if clusterMap.Data == nil || len(clusterMap.Data) == 0 {
							clusterMap.Data = configMap.Data
						}
						return nil
					})
					if err != nil {
						reqLogger.Error(err, "Failed to reconcile the "+configMap.GetName()+" ConfigMap")
						return r.ManageError(err, kappnavv1.StatusConditionTypeReconciled, instance)
					}
				}
			}
		}
	}

	// Create or update builtin config
	builtinConfig := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "builtin",
			Namespace: instance.GetNamespace(),
		},
	}
	err = r.CreateOrUpdate(builtinConfig, instance, func() error {
		kappnavutils.CustomizeConfigMap(builtinConfig, instance)
		kappnavutils.CustomizeBuiltinConfigMap(builtinConfig, &r.ReconcilerBase, instance)
		return nil
	})
	if err != nil {
		reqLogger.Error(err, "Failed to reconcile the kappnav-config ConfigMap")
		return r.ManageError(err, kappnavv1.StatusConditionTypeReconciled, instance)
	}

	// Create or update kappnav-config
	kappnavConfig := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "kappnav-config",
			Namespace: instance.GetNamespace(),
		},
	}
	err = r.CreateOrUpdate(kappnavConfig, instance, func() error {
		kappnavutils.CustomizeConfigMap(kappnavConfig, instance)
		kappnavutils.CustomizeKappnavConfigMap(kappnavConfig, kappnavURL, instance)
		return nil
	})
	if err != nil {
		reqLogger.Error(err, "Failed to reconcile the kappnav-config ConfigMap")
		return r.ManageError(err, kappnavv1.StatusConditionTypeReconciled, instance)
	}

	// Create or update the UI deployment
	uiDeployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.GetName() + "-ui",
			Namespace: instance.GetNamespace(),
		},
	}
	err = r.CreateOrUpdate(uiDeployment, instance, func() error {
		pts := &uiDeployment.Spec.Template
		kappnavutils.CustomizeDeployment(uiDeployment, instance)
		kappnavutils.CustomizePodSpec(pts, &uiDeployment.ObjectMeta,
			kappnavutils.CreateUIDeploymentContainers(pts.Spec.Containers, instance),
			kappnavutils.CreateUIVolumes(instance), instance)
		return nil
	})
	if err != nil {
		reqLogger.Error(err, "Failed to reconcile the UI Deployment")
		return r.ManageError(err, kappnavv1.StatusConditionTypeReconciled, instance)
	}

	// Create or update the Controller deployment
	controllerDeployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.GetName() + "-controller",
			Namespace: instance.GetNamespace(),
		},
	}
	err = r.CreateOrUpdate(controllerDeployment, instance, func() error {
		pts := &controllerDeployment.Spec.Template
		kappnavutils.CustomizeDeployment(controllerDeployment, instance)
		kappnavutils.CustomizePodSpec(pts, &controllerDeployment.ObjectMeta,
			kappnavutils.CreateControllerDeploymentContainers(pts.Spec.Containers, instance), nil, instance)
		return nil
	})
	if err != nil {
		reqLogger.Error(err, "Failed to reconcile the Controller Deployment")
		return r.ManageError(err, kappnavv1.StatusConditionTypeReconciled, instance)
	}

	// If an extension exists call its reconcile function, otherwise return success.
	if extension != nil {
		return extension.ReconcileAdditionalResources(request, &r.ReconcilerBase, instance)
	}
	return r.ManageSuccess(kappnavv1.StatusConditionTypeReconciled, instance)
}
