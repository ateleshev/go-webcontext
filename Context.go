package webcontext

import (
	"github.com/jinzhu/gorm"
	"sync"
)

const (
	DB = "db" // GORM (https://godoc.org/github.com/jinzhu/gorm)
)

type Context struct {
	sync.Mutex
	data map[string]interface{}
}

func NewContext() *Context { // {{{
	return &Context{data: make(map[string]interface{}, 0)}
} // }}}

func (this *Context) has(key string) bool { // {{{
	_, ok := this.data[key]

	return ok
} // }}}

func (this *Context) set(key string, value interface{}) bool { // {{{
	if this.has(key) {
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

// [DB]

func (this *Context) SetDB(value *gorm.DB) bool { // {{{
	return this.set(DB, value)
} // }}}

func (this *Context) GetDB() *gorm.DB { // {{{
	return this.get(DB).(*gorm.DB)
} // }}}
