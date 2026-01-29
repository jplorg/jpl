package builtins

import "github.com/jplorg/jpl/go/v2/jpl"

type enclosure struct {
	fn func(runtime jpl.JPLRuntime, signal jpl.JPLRuntimeSignal, next jpl.JPLPiper, input any, args ...any) ([]any, error)
}

func (e *enclosure) Call(runtime jpl.JPLRuntime, signal jpl.JPLRuntimeSignal, next jpl.JPLPiper, input any, args ...any) ([]any, error) {
	return e.fn(runtime, signal, next, input, args...)
}

func (e *enclosure) IsSame(other jpl.JPLFunc) bool {
	return e == other
}

func enclose(fn func(runtime jpl.JPLRuntime, signal jpl.JPLRuntimeSignal, next jpl.JPLPiper, input any, args ...any) ([]any, error)) jpl.JPLFunc {
	return &enclosure{fn: fn}
}
