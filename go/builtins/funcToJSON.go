package builtins

import (
	"github.com/2manyvcos/jpl/go/jpl"
	"github.com/2manyvcos/jpl/go/library"
)

var funcToJSON jpl.JPLFunc = func(runtime jpl.JPLRuntime, signal jpl.JPLRuntimeSignal, next jpl.JPLPiper, input any, args ...any) ([]any, error) {
	json, err := library.StringifyJSON(input, false)
	if err != nil {
		return nil, err
	}
	return next.Pipe(json)
}
