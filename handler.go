package blia

import (
	"encoding/json"
	"io"
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
}

func WriteJSON(w io.Writer, v interface{}) error {
	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)
	return encoder.Encode(v)
}
