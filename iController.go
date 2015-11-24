package webcontext

import (
	"net/http"

	"github.com/gorilla/mux"
)

type iController interface {
	// [Context]
	Context() *Context
	HasContext() bool

	// [Request]
	Request() *http.Request
	HasRequest() bool

	// [Controller]
	Initialize(context *Context)
	Register(router *mux.Router) *mux.Route
	Configure(request *http.Request)
	Prepare() error
	Render(writer http.ResponseWriter) error

	// [Error]
	Error(writer http.ResponseWriter)
}
