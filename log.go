package blia

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Logger interface {
	Info(context.Context, string, ...interface{})
	Warn(context.Context, string, ...interface{})
	Error(context.Context, string, ...interface{})
}

var std Logger

func SetLogger(log Logger) {
	std = log
}

func init() {
	SetLogger(new(fmtLogger))
}

type fmtLogger struct{}

func (log *fmtLogger) Info(_ context.Context, message string, args ...any) {
	fmt.Printf("[Info] %+v ", time.Now())
	fmt.Printf(message, args...)
	fmt.Println("")
}

func (log *fmtLogger) Warn(_ context.Context, message string, args ...any) {
	fmt.Printf("[Warn] %+v ", time.Now())
	fmt.Printf(message, args...)
	fmt.Println("")
}

func (log *fmtLogger) Error(_ context.Context, message string, args ...any) {
	fmt.Printf("[Error] %+v ", time.Now())
	fmt.Printf(message, args...)
	fmt.Println("")
}

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
