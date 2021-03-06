package webcontext

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux" // http://www.gorillatoolkit.org/pkg/mux
	"github.com/jinzhu/gorm" // https://godoc.org/github.com/jinzhu/gorm

	"github.com/ateleshev/go-webconfig"
)

const (
	DATE_FORMAT      = "2006-01-02"
	TIME_FORMAT      = "15:04:05"
	DATE_TIME_FORMAT = DATE_FORMAT + " " + TIME_FORMAT

	DEFAULT_SERVER_HOST = "localhost:8090"
	DEFAULT_SERVER_ADDR = "0.0.0.0:8090"
	DEFAULT_POOL_SIZE   = 4
	DEFAULT_QUEUE_SIZE  = 100

	NAMESPACE = "common"

	MEMORY_TEMPLATE = "%.3f%s"

	SERVER_TYPE_HTTP = "http"
	SERVER_TYPE_FCGI = "fcgi"
)

func NewContext(router *mux.Router) *Context { // {{{
	return &Context{
		startedAt: time.Now(),
		router:    router,
		data:      make(ContextNamespaceData, 0),
	}
} // }}}

func CreateContext(router *mux.Router, appConfig webconfig.ConfigInterface, configPath string) (*Context, error) { // {{{
	context := NewContext(router)

	config, errs := webconfig.LoadConfig(appConfig, configPath)
	if errs != nil && errs.Len() > 0 {
		for name, err := range *errs {
			go log.Printf("[Error:Context] Config '%s' is not loaded: %v\n", name, err)
		}
	}

	return context, context.Initialize(config)
} // }}}

type Context struct {
	sync.Mutex

	startedAt time.Time
	config    *webconfig.Config    // App configuration
	router    *mux.Router          // GorillaMux (http://www.gorillatoolkit.org/pkg/mux)
	db        *gorm.DB             // GORM (https://godoc.org/github.com/jinzhu/gorm)
	data      ContextNamespaceData // Namespaced data
}

// == Protected ==

func (this *Context) has(ns string, key string) bool { // {{{
	var ok bool
	if _, ok = this.data[ns]; ok {
		_, ok = this.data[ns][key]
	}

	return ok
} // }}}

func (this *Context) set(ns string, key string, value interface{}) bool { // {{{
	this.Lock()
	defer this.Unlock()

	if _, ok := this.data[ns]; !ok {
		this.data[ns] = make(ContextData)
	}

	this.data[ns][key] = value

	return true
} // }}}

func (this *Context) get(ns string, key string) interface{} { // {{{
	if !this.has(ns, key) {
		return nil
	}

	return this.data[ns][key]
} // }}}

// == Public ==

func (this *Context) Initialize(config *webconfig.Config) error { // {{{
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

	// [Database]

	if config.HasDatabase() {
		db, err := gorm.Open(config.Database.Driver, config.Database.DSN)
		if err != nil {
			return err
		}

		con := db.DB()

		con.Ping()
		if config.Database.MaxIdleConns > 0 {
			con.SetMaxIdleConns(config.Database.MaxIdleConns)
		}

		if config.Database.MaxOpenConns > 0 {
			con.SetMaxOpenConns(config.Database.MaxOpenConns)
		}

		db.SingularTable(config.Database.SingularTable)

		// Logger [Enable|Disable]
		db.LogMode(config.Database.LogMode)

		this.SetDB(db)
	}

	return nil
} // }}}

func (this *Context) StartedAt() time.Time { // {{{
	return this.startedAt
} // }}}

func (this *Context) Executed() time.Duration { // {{{
	return time.Since(this.startedAt)
} // }}}

func (this *Context) DateTimeFormat() string { // {{{
	return DATE_TIME_FORMAT
} // }}}

// [Config]

func (this *Context) HasConfig() bool { // {{{
	return this.config != nil
} // }}}

func (this *Context) SetConfig(config *webconfig.Config) bool { // {{{
	this.Lock()
	defer this.Unlock()

	this.config = config
	return true
} // }}}

func (this *Context) Config() *webconfig.Config { // {{{
	return this.config
} // }}}

// [Server]

func (this *Context) IsServerType(serverType string) bool { // {{{
	config := this.Config()
	if config.HasMain() {
		return config.Main.IsServerType(serverType)
	}

	return false
} // }}}

/**
 * Pool size
 */
func (this *Context) PoolSize() int { // {{{
	config := this.Config()
	if config.HasServer() && config.Server.PoolSize > 0 {
		return config.Server.PoolSize
	}

	return DEFAULT_POOL_SIZE
} // }}}

/**
 * Queue size
 */
func (this *Context) QueueSize() int { // {{{
	config := this.Config()
	if config.HasServer() && config.Server.QueueSize > 0 {
		return config.Server.QueueSize
	}

	return DEFAULT_QUEUE_SIZE
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

func (this *Context) Addr() string { // {{{
	if this.Config().HasServer() {
		return this.Config().Server.Addr
	}

	return DEFAULT_SERVER_ADDR
} // }}}

func (this *Context) Host() string { // {{{
	if this.Config().HasServer() {
		return this.Config().Server.Host
	}

	return DEFAULT_SERVER_HOST
} // }}}

func (this *Context) IsSecure() bool { // {{{
	if this.Config().HasServer() {
		return this.Config().Server.Secure
	}

	return false
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

func (this *Context) AddController(controller ControllerInterface) error { // {{{
	controller.Initialize(this)
	if router := this.Router(); router != nil {
		route := router.NewRoute()
		controller.Register(route)
		route.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
			c := controller.New()
			if err := c.Prepare(); err != nil {
				go log.Printf("[Error:Context] Cannot prepare controller '%v': %v", reflect.TypeOf(c), err)
				http.Error(responseWriter, "Forbidden", http.StatusForbidden) // @TODO: Add error
			} else {
				view := c.Execute(request)
				if err := view.Render(responseWriter); err != nil {
					go log.Printf("[Error:Context] Cannot render view '%v' %v", reflect.TypeOf(view), err)
				}
			}
		})
		return nil
	}

	return fmt.Errorf("[Error:Context] Cannot register controller, the router not intialized in this context")
} // }}}

// [URL]

func (this *Context) Url(p string, params map[string]string) string { // {{{
	var url string
	var scheme string = "http"

	if this.IsSecure() {
		scheme = "https"
	}

	url = fmt.Sprintf("%s://%s%s", scheme, this.Host(), p)
	if params != nil {
		for k, v := range params {
			url = strings.Replace(url, "{{"+k+"}}", v, 1)
		}
	}

	return url
} // }}}

// [Data]

func (this *Context) Has(name string) bool { // {{{
	return this.has(NAMESPACE, name)
} // }}}

func (this *Context) Set(name string, value interface{}) bool { // {{{
	return this.set(NAMESPACE, name, value)
} // }}}

func (this *Context) Get(name string) interface{} { // {{{
	return this.get(NAMESPACE, name)
} // }}}

// [Math]

func (this *Context) Log(v, b float64) float64 { // {{{
	return math.Log(v) / math.Log(b)
} // }}}

func (this *Context) Round(v float64, p int) float64 { // {{{
	var rounder float64
	intermed := v * math.Pow(10, float64(p))

	if v >= 0.5 {
		rounder = math.Ceil(intermed)
	} else {
		rounder = math.Floor(intermed)
	}

	return rounder / math.Pow(10, float64(p))
} // }}}

// [Memory]

func (this *Context) MemoryFormat(size uint64) string { // {{{
	s := float64(size)
	b := float64(1024)

	units := []string{"B", "Kb", "Mb", "Gb", "Tb", "Pb"}
	if size == 0 {
		return fmt.Sprintf(MEMORY_TEMPLATE, s, units[size])
	} else {
		i := math.Floor(this.Log(s, b))

		return fmt.Sprintf(MEMORY_TEMPLATE, (s / math.Pow(b, i)), units[int(i)])
	}
} // }}}

/**
 * General statistics.
 * [runtime.MemStats.Alloc] - bytes allocated and still in use
 *
 * Main allocation heap statistics.
 * [runtime.MemStats.HeapAlloc] - bytes allocated and still in use
 */
func (this *Context) MemoryUsage() string { // {{{
	ms := &runtime.MemStats{}
	runtime.ReadMemStats(ms)

	return this.MemoryFormat(ms.Alloc)
} // }}}
