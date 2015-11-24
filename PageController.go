package webcontext

import (
	"net/http"
	"os"

	wp "github.com/ArtemTeleshev/go-webpage"
)

type PageController struct {
	Controller

	// Page
	Page *wp.Page

	// Template
	TemplatePath string
	TemplateName string
}

func (this *PageController) Initialize(context *Context) { // {{{
	this.context = context

	// Template

	if len(this.TemplateName) == 0 {
		this.TemplateName = wp.DEFAULT_TEMPLATE_NAME
	}

	if len(this.TemplatePath) == 0 {
		this.TemplatePath, _ = os.Getwd()
	}
} // }}}

func (this *PageController) Render(writer http.ResponseWriter) error { // {{{
	template := wp.NewPageTemplate(this.TemplateName, this.TemplatePath)
	if err := template.Execute(writer, this.Page); err != nil {
		return err
	}

	return nil
} // }}}
