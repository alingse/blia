package blia

import (
	"bytes"
	"io"
	"net/http"
)

func HTTPFullLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body := &bytes.Buffer{}
		defer r.Body.Close()
		r.Body = io.NopCloser(io.TeeReader(r.Body, body))
		next.ServeHTTP(w, r)
		// TODO tee on w
		log := map[string]any{
			"header":          r.Header,
			"path":            r.URL.Path,
			"host":            r.Host,
			"method":          r.Method,
			"requst_body":     body.String(),
			"response_header": w.Header(),
		}
		std.Info(r.Context(), "%+v", log)
	})
}
