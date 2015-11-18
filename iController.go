package webcontext

import (
	"net/http"

	"github.com/gorilla/mux"
)

type iController interface {
	Initialize(context *Context)
	Register(router *mux.Router) *mux.Route
	Configure(request *http.Request) error
	Render(writer http.ResponseWriter) error
	Error(writer http.ResponseWriter)
}
