package controller

import (
	"github.com/terra-rebels/terra-operator/pkg/controller/terradnode"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, terradnode.Add)
}
