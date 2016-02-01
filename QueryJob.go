package webcontext

import (
	"log"
	"net/http"
	"time"

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
		startedAt := time.Now()
		this.Context.Router().ServeHTTP(this.ResponseWriter, this.Request)
		finishedAt := time.Now()

		serverType := SERVER_TYPE_HTTP
		if this.Context.HasConfig() && this.Context.Config().HasMain() {
			serverType = this.Context.Config().Main.ServerType
		}

		log.Printf("[Query:%s] Execute %s Request [%.4fs]\n", worker.Info(), serverType, finishedAt.Sub(startedAt).Seconds())
	} else {
		log.Printf("[QueryError:%s] Router not found\n", worker.Info())
	}
} // }}}
