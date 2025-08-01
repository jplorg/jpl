package builtins

import (
	"github.com/jplorg/jpl/go/jpl"
	"github.com/jplorg/jpl/go/library"
)

var funcError = enclose(func(runtime jpl.JPLRuntime, signal jpl.JPLRuntimeSignal, next jpl.JPLPiper, input any, args ...any) ([]any, error) {
	return nil, library.ThrowAny(library.NewRuntimeError(input))
})
