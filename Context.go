package webcontext

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/fcgi"
	"reflect"
	"runtime"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"

	"github.com/gorilla/mux" // http://www.gorillatoolkit.org/pkg/mux
	"github.com/jinzhu/gorm" // https://godoc.org/github.com/jinzhu/gorm

	"github.com/ArtemTeleshev/go-repository"
)

const (
	DEFAULT_SERVER_ADDR = "0.0.0.0:8090"

	NAMESPACE            = "common"
	NAMESPACE_REPOSITORY = "repository"
)

type Context struct {
	sync.Mutex

	config *Config      // App configuration
	server *http.Server // HttpServer (https://golang.org/pkg/net/http)
	router *mux.Router  // GorillaMux (http://www.gorillatoolkit.org/pkg/mux)
	db     *gorm.DB     // GORM (https://godoc.org/github.com/jinzhu/gorm)

	data map[string]interface{}
}

// == Protected ==

func buildKey(NAMESPACE string, name string) string { // {{{
	return fmt.Sprintf("%s::%s", NAMESPACE, name)
} // }}}

func (this *Context) has(key string) bool { // {{{
	_, ok := this.data[key]

	return ok
} // }}}

func (this *Context) set(key string, value interface{}) bool { // {{{
	if this.Has(key) {
		return false
	}

	this.Lock()
	defer this.Unlock()

	this.data[key] = value

	return true
} // }}}

func (this *Context) get(key string) interface{} { // {{{
	return this.data[key]
} // }}}

// == Static ==

func NewContext() *Context { // {{{
	return &Context{data: make(map[string]interface{}, 0)}
} // }}}

// == Public ==

func CreateContext(configPath string) (*Context, error) { // {{{
	context := NewContext()

	config, errs := LoadConfig(configPath)
	if errs != nil && errs.Len() > 0 {
		for name, err := range *errs {
			log.Printf("Config '%s' is not loaded. Details: %v", name, err)
		}
	}

	return context, context.Init(config)
} // }}}

func (this *Context) Init(config *Config) error { // {{{
	var err error
	this.SetConfig(config)

	// [Main]

	if config.HasMain() {
		if config.Main.UseThread {
			runtime.UnlockOSThread()
		}

		switch {
		case config.Main.MaxProcs > 0:
			runtime.GOMAXPROCS(config.Main.MaxProcs)
			break
		case config.Main.MaxProcs < 0:
			runtime.GOMAXPROCS(runtime.NumCPU())
			break
		}
	}

	// [Router]
	this.SetRouter(mux.NewRouter())

	// [Server]

	server := &http.Server{Addr: DEFAULT_SERVER_ADDR}

	if config.HasServer() {
		server.Addr = config.Server.Addr()
	}

	this.SetServer(server)

	// [Database]

	if config.HasDatabase() {
		var db gorm.DB
		if db, err = gorm.Open(config.Database.Driver, config.Database.DSN); err != nil {
			return err
		}

		db.DB().Ping()

		if config.Database.MaxIdleConns > 0 {
			db.DB().SetMaxIdleConns(config.Database.MaxIdleConns)
		}

		if config.Database.MaxOpenConns > 0 {
			db.DB().SetMaxOpenConns(config.Database.MaxOpenConns)
		}

		db.SingularTable(config.Database.SingularTable)

		// Logger [Enable|Disable]
		db.LogMode(config.Database.LogMode)

		this.SetDB(&db)
	}

	return nil
} // }}}

// [Config]

func (this *Context) HasConfig() bool { // {{{
	return this.config != nil
} // }}}

func (this *Context) SetConfig(config *Config) bool { // {{{
	this.Lock()
	defer this.Unlock()

	this.config = config
	return true
} // }}}

func (this *Context) Config() *Config { // {{{
	return this.config
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

func (this *Context) Server() *http.Server { // {{{
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

func (this *Context) Router() *mux.Router { // {{{
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

func (this *Context) DB() *gorm.DB { // {{{
	return this.db
} // }}}

// [Controllers]

func (this *Context) AddController(controller IController) error { // {{{
	controller.Initialize(this)
	if router := this.Router(); router != nil {
		controller.Register(router).HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			var err error
			controller.Configure(request)
			if err = controller.Prepare(); err != nil {
				log.Printf("Cannot prepare controller '%v'. Error: '%v'", reflect.TypeOf(controller), err)
				controller.Error(writer)
			} else {
				if err = controller.Render(writer); err != nil {
					log.Printf("Cannot render controller '%v'. Error: '%v'", reflect.TypeOf(controller), err)
				}
			}
		})
		return nil
	}

	return fmt.Errorf("Cannot register controller, the router not intialized in this context")
} // }}}

// [Data]

func (this *Context) Has(name string) bool { // {{{
	return this.has(buildKey(NAMESPACE, name))
} // }}}

func (this *Context) Set(name string, value interface{}) bool { // {{{
	return this.set(buildKey(NAMESPACE, name), value)
} // }}}

func (this *Context) Get(name string) interface{} { // {{{
	return this.get(buildKey(NAMESPACE, name))
} // }}}

// [Repository]

func (this *Context) HasRepository(name string) bool { // {{{
	return this.has(buildKey(NAMESPACE_REPOSITORY, name))
} // }}}

func (this *Context) SetRepository(name string, repository *repository.Repository) bool { // {{{
	return this.set(buildKey(NAMESPACE_REPOSITORY, name), repository)
} // }}}

func (this *Context) Repository(name string) *repository.Repository { // {{{
	if !this.has(buildKey(NAMESPACE_REPOSITORY, name)) {
		return nil
	}

	return this.get(buildKey(NAMESPACE_REPOSITORY, name)).(*repository.Repository)
} // }}}

// [Handle]

func (this *Context) Handle() { // {{{
	server := this.Server()
	server.Handler = this.Router()
} // }}}

func (this *Context) ListenAndServe() { // {{{
	config := this.Config()
	if config.HasMain() && config.Main.IsServerType(SERVER_TYPE_FCGI) {
		this.fcgiListenAndServe()
	} else {
		this.httpListenAndServe()
	}
} // }}}

func (this *Context) fcgiListenAndServe() { // {{{
	server := this.Server()

	log.Printf("Start FCGI Server (fcgi://%s)", server.Addr)
	listener, err := net.Listen("tcp", server.Addr)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	log.Fatal(fcgi.Serve(listener, server.Handler))
} // }}}

func (this *Context) httpListenAndServe() { // {{{
	server := this.Server()

	log.Printf("Start HTTP Server (http://%s)", server.Addr)
	log.Fatal(server.ListenAndServe())
} // }}}

func (this *Context) Dispatch() { // {{{
	this.Handle()
	this.ListenAndServe()
} // }}}
