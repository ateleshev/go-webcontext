package webcontext

import (
	"net/http"
	"os"

	wp "github.com/ArtemTeleshev/go-webpage"
)

type Controller struct {
	Context *Context
	Page    *wp.Page

	// Template
	TemplatePath string
	TemplateName string
}

func (this *Controller) Initialize(context *Context) { // {{{
	this.Context = context

	// Template

	if len(this.TemplateName) == 0 {
		this.TemplateName = wp.DEFAULT_TEMPLATE_NAME
	}

	if len(this.TemplatePath) == 0 {
		this.TemplatePath, _ = os.Getwd()
	}
} // }}}

/**
 * func (this *Controller) Configure(request *http.Request) error {
 *   // Controller configuration
 * }
 */

func (this *Controller) Render(writer http.ResponseWriter) error { // {{{
	template := wp.NewPageTemplate(this.TemplateName, this.TemplatePath)
	if err := template.Execute(writer, this.Page); err != nil {
		return err
	}

	return nil
} // }}}

func (this *Controller) Error(writer http.ResponseWriter) { // {{{
	http.Error(writer, "Forbidden", http.StatusForbidden)
} // }}}
