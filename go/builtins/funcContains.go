package builtins

import (
	"strings"

	"github.com/2manyvcos/jpl/go/jpl"
	"github.com/2manyvcos/jpl/go/library"
)

var funcContains jpl.JPLFunc = func(runtime jpl.JPLRuntime, signal jpl.JPLRuntimeSignal, next jpl.JPLPiper, input any, args ...any) ([]any, error) {
	value, err := library.UnwrapValue(input)
	if err != nil {
		return nil, err
	}
	t, err := library.Type(value)
	if err != nil {
		return nil, err
	}

	var token any
	if len(args) > 0 {
		token, err = library.UnwrapValue(args[0])
		if err != nil {
			return nil, err
		}
	}
	tt, err := library.Type(token)
	if err != nil {
		return nil, err
	}

	switch t {
	case jpl.JPLT_ARRAY:
		for _, item := range value.([]any) {
			equals, err := library.Equals(item, token)
			if err != nil {
				return nil, err
			}
			if equals {
				return next.Pipe(true)
			}
		}
		return next.Pipe(false)

	case jpl.JPLT_STRING:
		if tt == jpl.JPLT_STRING {
			return next.Pipe(strings.Contains(value.(string), token.(string)))
		}

	case jpl.JPLT_OBJECT:
		for _, item := range value.(map[string]any) {
			equals, err := library.Equals(item, token)
			if err != nil {
				return nil, err
			}
			if equals {
				return next.Pipe(true)
			}
		}
		return next.Pipe(false)

	default:
	}

	return nil, library.ThrowAny(library.NewTypeError("%s (%*<100v) cannot contain %s (%*<100v)", string(t), value, string(tt), token))
}
