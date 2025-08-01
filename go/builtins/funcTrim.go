package builtins

import (
	"strings"

	"github.com/jplorg/jpl/go/jpl"
	"github.com/jplorg/jpl/go/library"
)

var funcTrim = enclose(func(runtime jpl.JPLRuntime, signal jpl.JPLRuntimeSignal, next jpl.JPLPiper, input any, args ...any) ([]any, error) {
	t, err := library.Type(input)
	if err != nil {
		return nil, err
	}

	switch t {
	case jpl.JPLT_STRING:
		alteredValue, err := library.AlterValue(input, jpl.JPLModifierFunc(func(value any) (any, jpl.JPLError) {
			return strings.TrimSpace(value.(string)), nil
		}))
		if err != nil {
			return nil, err
		}
		return next.Pipe(alteredValue)
	}

	u, err := library.UnwrapValue(input)
	if err != nil {
		return nil, err
	}
	return nil, library.ThrowAny(library.NewTypeError("%s (%*<100v) cannot be trimmed", string(t), u))
})
