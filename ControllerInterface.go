package webcontext

import (
	"net/http"

	"github.com/gorilla/mux" // http://www.gorillatoolkit.org/pkg/mux
)

type ControllerInterface interface {
	New() ControllerInterface
	Context() *Context

	// [Initialization]
	Initialize(*Context)
	Register(*mux.Route)

	// [Execution]
	Prepare() error
	Execute(*http.Request) ViewInterface
}
