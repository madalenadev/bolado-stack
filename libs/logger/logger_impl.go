package logger

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// ILog interface
type ILog interface {
	Info(arguments ...interface{})
	Error(arguments ...interface{})
	Success(arguments ...interface{})
}

// Config of a logger
type Config struct {
	RequestID string
}

type level string

func (lvl level) toString() string {
	return string(lvl)
}

const (
	infoLevel    level = "\033[32m[INFO]\033[0m"
	successLevel level = "\033[32m[SUCCESS]\033[0m"
	errorLevel   level = "\033[31m[ERROR]\033[0m"
)

type contextKey string

const (
	keyRequestID contextKey = contextKey("request-id")
)

// WithRequestID generate request id
func WithRequestID(ctx context.Context) context.Context {
	return context.WithValue(ctx, keyRequestID, uuid.New().String())
}

// FromContext return a new instance of ILog
func FromContext(ctx context.Context) ILog {
	val := ctx.Value(keyRequestID)
	id, ok := val.(string)
	if !ok {
		id = "empty"
	}
	return &log{id}
}

type log struct {
	requestID string
}

func (l log) Info(arguments ...interface{}) {
	l.send(infoLevel, arguments...)
}

func (l log) Success(arguments ...interface{}) {
	l.send(successLevel, arguments...)
}

func (l log) Error(arguments ...interface{}) {
	l.send(errorLevel, arguments...)
}

func (l log) send(lvl level, arguments ...interface{}) {
	fmt.Printf("%+v request_id: %+v ", lvl.toString(), l.requestID)
	fmt.Println(arguments...)
}
