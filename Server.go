package webcontext

import (
	"log"
	"net"
	"net/http"
	"net/http/fcgi"
	"time"

	"github.com/ateleshev/go-queue"
)

func NewServer(context *Context, name string) *Server { // {{{
	return &Server{
		name:        name,
		context:     context,
		queueServer: queue.NewServer(name, context.PoolSize(), context.QueueSize()),
		// Timing
		startedAt: time.Now(),
	}
} // }}}

type Server struct {
	http.Server

	name        string
	context     *Context
	queueServer *queue.Server

	// Logging
	accessLog *log.Logger

	// Timing
	startedAt time.Time
}

// == fcgi/http ==

func (this *Server) fcgiListenAndServe(listener net.Listener) { // {{{
	log.Printf("[Server:%s] Start %s server ... (%s)", this.Name(), SERVER_TYPE_FCGI, this.context.Url("/", nil))
	log.Fatal(fcgi.Serve(listener, this))
} // }}}

func (this *Server) httpListenAndServe(listener net.Listener) { // {{{
	log.Printf("[Server:%s] Start %s server ... (%s)", this.Name(), SERVER_TYPE_HTTP, this.context.Url("/", nil))
	log.Fatal(http.Serve(listener, this))
} // }}}

// == Public ==

func (this *Server) Name() string { // {{{
	return this.name
} // }}}

func (this *Server) Run() { // {{{
	this.queueServer.Start()

	if this.Addr == "" {
		this.Addr = this.context.Addr()
	}

	listener, err := net.Listen("tcp", this.Addr)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	if this.context.IsServerType(SERVER_TYPE_FCGI) {
		this.fcgiListenAndServe(listener)
	} else {
		this.httpListenAndServe(listener)
	}
} // }}}

func (this *Server) Close() { // {{{
	this.queueServer.Stop()
} // }}}

func (this *Server) Start() { // {{{
	go this.Run()
} // }}}

func (this *Server) Stop() { // {{{
	go this.Close()
} // }}}

func (this *Server) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) { // {{{
	job := NewQueryJob(this.context, responseWriter, request)
	this.queueServer.Dispatch(job)
	job.Wait() // Waiting without timeout
	// job.WaitDeadline(time.Minute)
} // }}}
