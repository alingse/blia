package blia

import (
	"context"
	"fmt"
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
