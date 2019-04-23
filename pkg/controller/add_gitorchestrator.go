package controller

import (
	"github.com/wiardvanrij/gitorchestrator/pkg/controller/gitorchestrator"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, gitorchestrator.Add)
}
