package builtins

import (
	"github.com/2manyvcos/jpl/go/jpl"
)

var funcVoid jpl.JPLFunc = func(runtime jpl.JPLRuntime, signal jpl.JPLRuntimeSignal, next jpl.JPLPiper, input any, args ...any) ([]any, error) {
	return nil, nil
}
