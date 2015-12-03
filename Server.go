package webcontext

import (
	"log"
	"net"
	"net/http"
	"net/http/fcgi"
	//	"time"

	"github.com/ArtemTeleshev/go-queue"
)

func NewServer(context *Context, name string) *Server { // {{{
	return &Server{
		name:     name,
		workers:  make(map[int]*queue.Worker),
		shutdown: make(chan bool),
		// [Public]
		Context:     context,
		QueryQueue:  make(QueryQueue, context.DepthTasksQueue()),
		WorkerQueue: make(queue.WorkerQueue, context.DepthWorkersQueue()),
	}
} // }}}

type Server struct {
	http.Server
	// [Protected]
	name     string
	workers  map[int]*queue.Worker
	shutdown chan bool
	// [Public]
	Context     *Context
	QueryQueue  QueryQueue
	WorkerQueue queue.WorkerQueue
	AccessLog   *log.Logger
}

// == Protected ==

func (this *Server) startWorkers() { // {{{
	numServerWorkers := this.Context.NumServerWorkers()
	for i := 0; i < numServerWorkers; i++ {
		worker := queue.NewWorker(this.Name(), i+1, this.WorkerQueue)
		worker.Start()
		this.workers[worker.Id()] = worker
	}
	log.Printf("[Server:%s] Run %d workers\n", this.Name(), numServerWorkers)
} // }}}

func (this *Server) stopWorkers() { // {{{
	for id, worker := range this.workers {
		worker.Stop()
		delete(this.workers, id)
	}
} // }}}

func (this *Server) execute(query *Query) { // {{{
	worker := <-this.WorkerQueue
	// Dispatching work request
	worker <- query
} // }}}

func (this *Server) dispatch() { // {{{
	for {
		select {
		case query := <-this.QueryQueue:
			// Received work requeust
			go this.execute(query)
		case <-this.shutdown:
			this.stopWorkers()
			log.Printf("[Server:%s] Shutdown\n", this.Name())
			return
		}
	}
} // }}}

// == Public ===

func (this *Server) Name() string { // {{{
	return this.name
} // }}}

func (this *Server) Run() { // {{{
	this.startWorkers()
	go this.dispatch()

	if this.Addr == "" {
		this.Addr = this.Context.Addr()
	}

	listener, err := net.Listen("tcp", this.Addr)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	if this.Context.IsServerType(SERVER_TYPE_FCGI) {
		this.fcgiListenAndServe(listener)
	} else {
		this.httpListenAndServe(listener)
	}
} // }}}

func (this *Server) Close() { // {{{
	this.stopWorkers()
	this.shutdown <- true
} // }}}

func (this *Server) Start() { // {{{
	go this.Run()
} // }}}

func (this *Server) Stop() { // {{{
	go this.Close()
} // }}}

func (this *Server) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) { // {{{
	query := NewQuery(this.Context, responseWriter, request)
	this.QueryQueue <- query
	query.Wait() // Waiting without timeout
	// query.WaitDeadline(time.Minute)
} // }}}

func (this *Server) fcgiListenAndServe(listener net.Listener) { // {{{
	log.Printf("[Server:%s] Start (fcgi://%s)", this.Name(), this.Addr)
	log.Fatal(fcgi.Serve(listener, this))
} // }}}

func (this *Server) httpListenAndServe(listener net.Listener) { // {{{
	log.Printf("[Server:%s] Start (http://%s)", this.Name(), this.Addr)
	log.Fatal(http.Serve(listener, this))
} // }}}
