package blia

import (
	"context"
	"encoding/json"
	"net/http"
)

type HandlerFuncErr func(w http.ResponseWriter, r *http.Request) error

func HandleErr(h HandlerFuncErr) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := h(w, r)
		if err != nil {
			WriteError(w, err)
		}
		_, _ = w.Write([]byte(""))
	}
}

type errorResponse struct {
	Error *HTTPError `json:"error"`
}

func WriteError(w http.ResponseWriter, err error) {
	httpErr := NewHTTPError(err)
	w.WriteHeader(httpErr.ErrStatusCode())
	_ = WriteJSON(w, errorResponse{Error: httpErr})

	std.Error(context.Background(), "httpErr %+v err %+v ", httpErr, err)
}

func WriteJSON(w http.ResponseWriter, v interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)
	return encoder.Encode(v)
}
