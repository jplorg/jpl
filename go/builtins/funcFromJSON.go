package builtins

import (
	"encoding/json"

	"github.com/2manyvcos/jpl/go/jpl"
	"github.com/2manyvcos/jpl/go/library"
)

var funcFromJSON jpl.JPLFunc = func(runtime jpl.JPLRuntime, signal jpl.JPLRuntimeSignal, next jpl.JPLPiper, input any, args ...any) ([]any, error) {
	value, err := library.UnwrapValue(input)
	if err != nil {
		return nil, err
	}
	t, err := library.Type(value)
	if err != nil {
		return nil, err
	}

	switch t {
	case jpl.JPLT_STRING:
		var result any
		err := json.Unmarshal([]byte(value.(string)), &result)
		if err != nil {
			return nil, library.ThrowAny(library.NewRuntimeError(err.Error()))
		}
		return next.Pipe(result)

	default:
	}

	return nil, library.ThrowAny(library.NewTypeError("%s (%*<100v) cannot be parsed as JSON", string(t), value))
}
