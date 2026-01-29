package builtins

import (
	"github.com/jplorg/jpl/go/v2/jpl"
	"github.com/jplorg/jpl/go/v2/library"
)

var funcToJSON = enclose(func(runtime jpl.JPLRuntime, signal jpl.JPLRuntimeSignal, next jpl.JPLPiper, input any, args ...any) ([]any, error) {
	json, err := library.StringifyJSON(input, false)
	if err != nil {
		return nil, err
	}
	return next.Pipe(json)
})
