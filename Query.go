package webcontext

import (
	"log"
	"net/http"

	"github.com/ArtemTeleshev/go-queue"
)

func NewQuery(context *Context, responseWriter http.ResponseWriter, request *http.Request) *Query { // {{{
	query := &Query{
		Context:        context,
		Request:        request,
		ResponseWriter: responseWriter,
	}

	query.Init()

	return query
} // }}}

type Query struct {
	queue.Job

	Context        *Context
	Request        *http.Request
	ResponseWriter http.ResponseWriter
}

func (this *Query) Execute() { // {{{
	if this.Context.HasRouter() {
		this.Context.Router().ServeHTTP(this.ResponseWriter, this.Request)
		// log.Printf("[Query] Execute Query\n")
	} else {
		log.Printf("[Query] Router not found\n")
	}
} // }}}
