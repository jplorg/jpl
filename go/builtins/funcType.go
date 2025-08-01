package builtins

import (
	"github.com/jplorg/jpl/go/jpl"
	"github.com/jplorg/jpl/go/library"
)

var funcType = enclose(func(runtime jpl.JPLRuntime, signal jpl.JPLRuntimeSignal, next jpl.JPLPiper, input any, args ...any) ([]any, error) {
	t, err := library.Type(input)
	if err != nil {
		return nil, err
	}
	return next.Pipe(string(t))
})
