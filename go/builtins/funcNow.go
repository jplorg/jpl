package builtins

import (
	"time"

	"github.com/jplorg/jpl/go/v2/jpl"
)

var funcNow = enclose(func(runtime jpl.JPLRuntime, signal jpl.JPLRuntimeSignal, next jpl.JPLPiper, input any, args ...any) ([]any, error) {
	return next.Pipe(float64(time.Now().UnixMilli()))
})
