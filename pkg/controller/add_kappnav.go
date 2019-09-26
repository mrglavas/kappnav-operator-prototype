package controller

import (
	"github.com/kappnav/operator/pkg/controller/kappnav"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, kappnav.Add)
}
