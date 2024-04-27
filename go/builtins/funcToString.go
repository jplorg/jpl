package builtins

import (
	"github.com/2manyvcos/jpl/go/jpl"
	"github.com/2manyvcos/jpl/go/library"
)

var funcToString jpl.JPLFunc = func(runtime jpl.JPLRuntime, signal jpl.JPLRuntimeSignal, next jpl.JPLPiper, input any, args ...any) ([]any, error) {
	result, err := library.StringifyJSON(input, true)
	if err != nil {
		return nil, err
	}
	return next.Pipe(result)
}
