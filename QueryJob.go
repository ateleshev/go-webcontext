package webcontext

import (
	"log"
	"net"
	"net/http"
	"time"

	"github.com/ateleshev/go-queue"
)

func NewQueryJob(context *Context, responseWriter http.ResponseWriter, request *http.Request) *QueryJob { // {{{
	job := &QueryJob{
		Context:        context,
		Request:        request,
		ResponseWriter: responseWriter,
	}

	job.RemoteIp, _, _ = net.SplitHostPort(request.RemoteAddr)
	job.Initialize()

	return job
} // }}}

type QueryJob struct {
	queue.Job

	Context        *Context
	Request        *http.Request
	ResponseWriter http.ResponseWriter

	RemoteIp string
}

func (this *QueryJob) Execute(w interface{}) { // {{{
	worker := w.(*queue.Worker)

	if this.Context.HasRouter() {
		startedAt := time.Now()
		this.Context.Router().ServeHTTP(this.ResponseWriter, this.Request)
		finishedAt := time.Now()

		go log.Printf("%s [%s] %s %v [%.4fs]\n", this.RemoteIp, worker.Info(), this.Request.Method, this.Request.URL, finishedAt.Sub(startedAt).Seconds())
	} else {
		go log.Printf("%s [%s:Error] Router not configured\n", this.RemoteIp, worker.Info())
	}
} // }}}
