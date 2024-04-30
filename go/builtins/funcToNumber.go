package builtins

import (
	"strconv"

	"github.com/jplorg/jpl/go/jpl"
	"github.com/jplorg/jpl/go/library"
)

var funcToNumber jpl.JPLFunc = func(runtime jpl.JPLRuntime, signal jpl.JPLRuntimeSignal, next jpl.JPLPiper, input any, args ...any) ([]any, error) {
	value, err := library.UnwrapValue(input)
	if err != nil {
		return nil, err
	}
	t, err := library.Type(value)
	if err != nil {
		return nil, err
	}

	switch t {
	case jpl.JPLT_NUMBER:
		return next.Pipe(value)

	case jpl.JPLT_STRING:
		parsed, err := strconv.ParseFloat(value.(string), 64)
		if err != nil {
			return nil, library.ThrowAny(library.NewTypeConversionError("%s (%*<100v) does not contain a valid number", string(t), value))
		}
		result, err := library.NormalizeValue(parsed)
		if err != nil {
			return nil, err
		}
		return next.Pipe(result)

	default:
	}

	return nil, library.ThrowAny(library.NewTypeError("%s (%*<100v) cannot be parsed as a number", string(t), value))
}
