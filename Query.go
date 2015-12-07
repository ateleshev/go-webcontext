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

	query.Initialize()

	return query
} // }}}

type Query struct {
	queue.Job

	Context        *Context
	Request        *http.Request
	ResponseWriter http.ResponseWriter
}

func (this *Query) Execute(w interface{}) { // {{{
	worker := w.(*queue.Worker)
	if this.Context.HasRouter() {
		this.Context.Router().ServeHTTP(this.ResponseWriter, this.Request)
		log.Printf("[Query:%s#%d] Execute HTTP Request\n", worker.Name(), worker.Id())
	} else {
		log.Printf("[Query:%s#%d] Router not found\n", worker.Name(), worker.Id())
	}
} // }}}
