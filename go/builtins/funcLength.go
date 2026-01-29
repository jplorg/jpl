package builtins

import (
	jpl "github.com/jplorg/jpl/go/v2/jpl"
	library "github.com/jplorg/jpl/go/v2/library"
)

var funcLength = enclose(func(runtime jpl.JPLRuntime, signal jpl.JPLRuntimeSignal, next jpl.JPLPiper, input any, args ...any) ([]any, error) {
	value, err := library.UnwrapValue(input)
	if err != nil {
		return nil, err
	}
	t, err := library.Type(value)
	if err != nil {
		return nil, err
	}

	switch t {
	case jpl.JPLT_NULL:
		return next.Pipe(0.)

	case jpl.JPLT_ARRAY:
		return next.Pipe(float64(len(value.([]any))))

	case jpl.JPLT_STRING:
		return next.Pipe(float64(len([]rune(value.(string)))))

	case jpl.JPLT_OBJECT:
		return next.Pipe(float64(len(value.(map[string]any))))

	default:
	}

	return nil, library.ThrowAny(library.NewTypeError("%s (%*<100v) has no length", string(t), value))
})
