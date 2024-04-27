package builtins

import (
	"github.com/2manyvcos/jpl/go/jpl"
	"github.com/2manyvcos/jpl/go/library"
)

var funcType jpl.JPLFunc = func(runtime jpl.JPLRuntime, signal jpl.JPLRuntimeSignal, next jpl.JPLPiper, input any, args ...any) ([]any, error) {
	t, err := library.Type(input)
	if err != nil {
		return nil, err
	}
	return next.Pipe(string(t))
}
