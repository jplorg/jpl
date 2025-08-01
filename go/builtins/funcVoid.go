package builtins

import (
	"github.com/jplorg/jpl/go/jpl"
)

var funcVoid = enclose(func(runtime jpl.JPLRuntime, signal jpl.JPLRuntimeSignal, next jpl.JPLPiper, input any, args ...any) ([]any, error) {
	return nil, nil
})
