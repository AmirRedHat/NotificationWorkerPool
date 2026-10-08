package handler

import (
	"encoding/json"
	"net/http"
	"sync"
)

type BaseHandler struct {
	mut sync.Mutex
}

func (h *BaseHandler) MakeResponse(response http.ResponseWriter, data any) error {
	response.Header().Set("Content-Type", "application/json")
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, err = response.Write(jsonData)
	return err
}
