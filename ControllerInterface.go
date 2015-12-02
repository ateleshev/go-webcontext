package webcontext

import (
	"net/http"

	"github.com/gorilla/mux"
)

type ControllerInterface interface {
	// [Context]
	HasContext() bool
	Context() *Context

	// [Request]
	HasRequest() bool
	Request() *http.Request

	// [Error]
	HasError() bool
	Error() string

	// [Controller]
	Initialize(context *Context)
	Register(router *mux.Router) *mux.Route
	Configure(request *http.Request)
	Prepare() error

	// [Render]
	Render(writer http.ResponseWriter) error
	RenderError(writer http.ResponseWriter)
}
