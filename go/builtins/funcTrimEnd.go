package builtins

import (
	"strings"
	"unicode"

	"github.com/jplorg/jpl/go/jpl"
	"github.com/jplorg/jpl/go/library"
)

var funcTrimEnd = enclose(func(runtime jpl.JPLRuntime, signal jpl.JPLRuntimeSignal, next jpl.JPLPiper, input any, args ...any) ([]any, error) {
	t, err := library.Type(input)
	if err != nil {
		return nil, err
	}

	switch t {
	case jpl.JPLT_STRING:
		alteredValue, err := library.AlterValue(input, jpl.JPLModifierFunc(func(value any) (any, jpl.JPLError) {
			return strings.TrimRightFunc(value.(string), unicode.IsSpace), nil
		}))
		if err != nil {
			return nil, err
		}
		return next.Pipe(alteredValue)
	}

	return nil, nil
})
