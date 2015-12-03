package webcontext

import (
	"net/http"
)

type ViewInterface interface {
	Context() *Context
	Request() *http.Request

	// [Initialization]
	Initialize(*Context, *http.Request)

	// [Execution]
	Render(http.ResponseWriter) error
}
