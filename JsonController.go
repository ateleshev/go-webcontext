package webcontext

import (
	"encoding/json"
	"net/http"
)

const (
	CONTENT_TYPE_JSON = "application/json"
)

type JsonController struct {
	Controller

	Data interface{}
}

func (this *JsonController) Render(writer http.ResponseWriter) error { // {{{
	writer.Header().Set("Content-Type", CONTENT_TYPE_JSON)
	// Write headers
	writer.WriteHeader(http.StatusOK)

	json, err := json.Marshal(this.Data)
	if err != nil {
		return err
	}

	// Write body
	writer.Write(json)

	return nil
} // }}}
