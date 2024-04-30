package builtins

import (
	"time"

	"github.com/jplorg/jpl/go/jpl"
)

var funcNow jpl.JPLFunc = func(runtime jpl.JPLRuntime, signal jpl.JPLRuntimeSignal, next jpl.JPLPiper, input any, args ...any) ([]any, error) {
	return next.Pipe(float64(time.Now().UnixMilli()))
}
