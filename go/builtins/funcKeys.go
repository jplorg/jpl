package builtins

import (
	"github.com/jplorg/jpl/go/v2/jpl"
	"github.com/jplorg/jpl/go/v2/library"
)

var funcKeys = enclose(func(runtime jpl.JPLRuntime, signal jpl.JPLRuntimeSignal, next jpl.JPLPiper, input any, args ...any) ([]any, error) {
	value, err := library.UnwrapValue(input)
	if err != nil {
		return nil, err
	}
	t, err := library.Type(value)
	if err != nil {
		return nil, err
	}

	switch t {
	case jpl.JPLT_ARRAY:
		v := value.([]any)
		result := make([]any, len(v))
		for i := range v {
			result[i] = float64(i)
		}
		return next.Pipe(result)

	case jpl.JPLT_OBJECT:
		v := value.(map[string]any)
		result := make([]any, 0, len(v))
		for k := range v {
			result = append(result, k)
		}
		return next.Pipe(result)

	default:
	}

	return nil, library.ThrowAny(library.NewTypeError("%s (%*<100v) has no keys", string(t), value))
})
