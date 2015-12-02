package webcontext

import (
	"net/http"
	"net/url"
)

type RouteInterface interface {
	GetError() error
	BuildOnly() *RouteInterface
	Handler(handler http.Handler) *RouteInterface
	HandlerFunc(f func(http.ResponseWriter, *http.Request)) *RouteInterface
	GetHandler() http.Handler
	Name(name string) *RouteInterface
	GetName() string
	Headers(pairs ...string) *RouteInterface
	HeadersRegexp(pairs ...string) *RouteInterface
	Host(tpl string) *RouteInterface
	Methods(methods ...string) *RouteInterface
	Path(tpl string) *RouteInterface
	PathPrefix(tpl string) *RouteInterface
	Queries(pairs ...string) *RouteInterface
	Schemes(schemes ...string) *RouteInterface
	Subrouter() *RouterInterface
	URL(pairs ...string) (*url.URL, error)
	URLHost(pairs ...string) (*url.URL, error)
	URLPath(pairs ...string) (*url.URL, error)
}
