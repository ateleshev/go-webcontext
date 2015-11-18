package webcontext

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"sync"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

type Context struct {
	sync.Mutex

	server *http.Server // HttpServer (https://golang.org/pkg/net/http)
	router *mux.Router  // GorillaMux (http://www.gorillatoolkit.org/pkg/mux)
	db     *gorm.DB     // GORM (https://godoc.org/github.com/jinzhu/gorm)

	data map[string]interface{}
}

func NewContext() *Context { // {{{
	return &Context{data: make(map[string]interface{}, 0)}
} // }}}

func CreateContext(server *http.Server, router *mux.Router, db *gorm.DB) *Context { // {{{
	context := NewContext()
	context.SetServer(server)
	context.SetRouter(router)
	context.SetDB(db)

	return context
} // }}}

// [Server]

func (this *Context) HasServer() bool { // {{{
	return this.server != nil
} // }}}

func (this *Context) SetServer(server *http.Server) bool { // {{{
	this.Lock()
	defer this.Unlock()

	this.server = server
	return true
} // }}}

func (this *Context) GetServer() *http.Server { // {{{
	return this.server
} // }}}

// [Router]

func (this *Context) HasRouter() bool { // {{{
	return this.router != nil
} // }}}

func (this *Context) SetRouter(router *mux.Router) bool { // {{{
	this.Lock()
	defer this.Unlock()

	this.router = router
	return true
} // }}}

func (this *Context) GetRouter() *mux.Router { // {{{
	return this.router
} // }}}

// [DB]

func (this *Context) HasDB() bool { // {{{
	return this.db != nil
} // }}}

func (this *Context) SetDB(db *gorm.DB) bool { // {{{
	this.Lock()
	defer this.Unlock()

	this.db = db
	return true
} // }}}

func (this *Context) GetDB() *gorm.DB { // {{{
	return this.db
} // }}}

// [Controllers]

func (this *Context) AddController(controller iController) error { // {{{
	controller.Initialize(this)
	if router := this.GetRouter(); router != nil {
		controller.Register(router).HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			var err error
			if err = controller.Configure(request); err != nil {
				log.Printf("Cannot configure controller '%v'. Error: '%v'", reflect.TypeOf(controller), err)
				controller.Error(writer)
			} else {
				if err = controller.Render(writer); err != nil {
					log.Printf("Cannot execute controller '%v'. Error: '%v'", reflect.TypeOf(controller), err)
				}
			}
		})
		return nil
	}

	return fmt.Errorf("Cannot register controller, the router not intialized in this context")
} // }}}

// [Data]

func (this *Context) Has(key string) bool { // {{{
	_, ok := this.data[key]

	return ok
} // }}}

func (this *Context) Set(key string, value interface{}) bool { // {{{
	if this.Has(key) {
		return false
	}

	this.Lock()
	defer this.Unlock()

	this.data[key] = value

	return true
} // }}}

func (this *Context) Get(key string) interface{} { // {{{
	return this.data[key]
} // }}}

// [Handle]

func (this *Context) Handle() { // {{{
	server := this.GetServer()
	server.Handler = this.GetRouter()

	log.Printf("Start HTTP Server (http://%s)", server.Addr)
	log.Fatal(server.ListenAndServe())
} // }}}
