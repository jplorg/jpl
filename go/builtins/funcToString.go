package builtins

import (
	"github.com/jplorg/jpl/go/jpl"
	"github.com/jplorg/jpl/go/library"
)

var funcToString jpl.JPLFunc = func(runtime jpl.JPLRuntime, signal jpl.JPLRuntimeSignal, next jpl.JPLPiper, input any, args ...any) ([]any, error) {
	result, err := library.StringifyJSON(input, true)
	if err != nil {
		return nil, err
	}
	return next.Pipe(result)
}
