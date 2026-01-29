package builtins

import (
	"github.com/jplorg/jpl/go/v2/jpl"
	"github.com/jplorg/jpl/go/v2/library"
)

var funcIn = enclose(func(runtime jpl.JPLRuntime, signal jpl.JPLRuntimeSignal, next jpl.JPLPiper, input any, args ...any) ([]any, error) {
	key, err := library.UnwrapValue(input)
	if err != nil {
		return nil, err
	}
	tk, err := library.Type(key)
	if err != nil {
		return nil, err
	}

	var value any
	if len(args) > 0 {
		value, err = library.UnwrapValue(args[0])
		if err != nil {
			return nil, err
		}
	}
	tv, err := library.Type(value)
	if err != nil {
		return nil, err
	}

	switch tv {
	case jpl.JPLT_ARRAY:
		if tk == jpl.JPLT_NUMBER {
			k := key.(float64)
			return next.Pipe(k >= 0 && int(k) < len(value.([]any)))
		}

	case jpl.JPLT_OBJECT:
		if tk == jpl.JPLT_STRING {
			_, ok := value.(map[string]any)[key.(string)]
			return next.Pipe(ok)
		}

	default:
	}

	return nil, library.ThrowAny(library.NewTypeError("%s (%*<100v) cannot have %s (%*<100v) as key", string(tv), value, string(tk), key))
})
