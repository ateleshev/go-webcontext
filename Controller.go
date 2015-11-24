package webcontext

import (
	"net/http"
)

type Controller struct {
	iController

	context *Context
	request *http.Request
}

// [Context]

func (this *Controller) Context() *Context { // {{{
	return this.context
} // }}}

func (this *Controller) HasContext() bool { // {{{
	return this.context != nil
} // }}}

// [Request]

func (this *Controller) Request() *http.Request { // {{{
	return this.request
} // }}}

func (this *Controller) HasRequest() bool { // {{{
	return this.request != nil
} // }}}

// [Controller]

func (this *Controller) Initialize(context *Context) { // {{{
	this.context = context
} // }}}

func (this *Controller) Configure(request *http.Request) { // {{{
	this.request = request
} // }}}

// [Error]

func (this *Controller) Error(writer http.ResponseWriter) { // {{{
	http.Error(writer, "Forbidden", http.StatusForbidden)
} // }}}
