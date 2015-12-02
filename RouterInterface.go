package webcontext

import (
	"net/http"
)

type RouterInterface interface {
	ServeHTTP(w http.ResponseWriter, req *http.Request)
	Get(name string) *RouteInterface
	GetRoute(name string) *RouteInterface
	StrictSlash(value bool) *RouterInterface
	NewRoute() *RouteInterface
	Handle(path string, handler http.Handler) *RouteInterface
	HandleFunc(path string, f func(http.ResponseWriter, *http.Request)) *RouteInterface
	Headers(pairs ...string) *RouteInterface
	Host(tpl string) *RouteInterface
	Methods(methods ...string) *RouteInterface
	Path(tpl string) *RouteInterface
	PathPrefix(tpl string) *RouteInterface
	Queries(pairs ...string) *RouteInterface
	Schemes(schemes ...string) *RouteInterface
}
