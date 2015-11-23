package webcontext

import (
	"net/http"
	"path"

	"github.com/gorilla/mux"

	wp "github.com/ArtemTeleshev/go-webpage"
)

type WebContentController struct {
	Controller

	request *http.Request
}

func NewWebContentController() *WebContentController { // {{{
	return &WebContentController{}
} // }}}

func (this *WebContentController) Register(router *mux.Router) *mux.Route { // {{{
	return router.NewRoute().PathPrefix("/")
} // }}}

func (this *WebContentController) Configure(request *http.Request) error { // {{{
	this.request = request

	return nil
} // }}}

func (this *WebContentController) Render(writer http.ResponseWriter) error { // {{{
	dir := path.Join(this.TemplatePath, wp.TEMPLATE_DIR_MAIN, this.TemplateName, wp.TEMPLATE_DIR_WEB)
	http.FileServer(http.Dir(dir)).ServeHTTP(writer, this.request)

	return nil
} // }}}