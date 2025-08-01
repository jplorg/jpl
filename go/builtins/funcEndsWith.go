package builtins

import (
	"strings"

	"github.com/jplorg/jpl/go/jpl"
	"github.com/jplorg/jpl/go/library"
)

var funcEndsWith = enclose(func(runtime jpl.JPLRuntime, signal jpl.JPLRuntimeSignal, next jpl.JPLPiper, input any, args ...any) ([]any, error) {
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
	case jpl.JPLT_STRING:
		if tt == jpl.JPLT_STRING {
			return next.Pipe(strings.HasSuffix(value.(string), token.(string)))
		}

	default:
	}

	return nil, library.ThrowAny(library.NewTypeError("%s (%*<100v) cannot end with %s (%*<100v)", string(t), value, string(tt), token))
})
