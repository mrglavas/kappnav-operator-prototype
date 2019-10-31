package controller

import (
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	routev1 "github.com/openshift/api/route/v1"

	"sigs.k8s.io/controller-runtime/pkg/manager"
)

// AddToManagerFuncs is a list of functions to add all Controllers to the Manager
var AddToManagerFuncs []func(manager.Manager) error

// AddToManager adds all Controllers to the Manager
func AddToManager(m manager.Manager) error {
	if err := apiextensionsv1beta1.AddToScheme(m.GetScheme()); err != nil {
		return err
	}

	if err := routev1.AddToScheme(m.GetScheme()); err != nil {
		return err
	}

	for _, f := range AddToManagerFuncs {
		if err := f(m); err != nil {
			return err
		}
	}
	return nil
}
