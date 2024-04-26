package program

import (
	"github.com/2manyvcos/jpl/go/definition"
	"github.com/2manyvcos/jpl/go/jpl"
	"github.com/2manyvcos/jpl/go/library"
)

type opaIter struct{}

// { optional: boolean }
func (opaIter) OP(runtime jpl.JPLRuntime, input any, target any, params definition.JPLSelectorParams, scope jpl.JPLRuntimeScope, next jpl.JPLPiper) ([]any, jpl.JPLError) {
	value, err := library.UnwrapValue(target)
	if err != nil {
		return nil, err
	}
	t, err := library.Type(value)
	if err != nil {
		return nil, err
	}
	switch t {
	case jpl.JPLT_OBJECT:
		return library.MuxAll([][]any{library.GetMapValues(value.(map[string]any))}, library.NewPiperMuxer(next))

	case jpl.JPLT_ARRAY:
		return library.MuxAll([][]any{value.([]any)}, library.NewPiperMuxer(next))

	case jpl.JPLT_STRING:
		return library.MuxAll([][]rune{[]rune(value.(string))}, jpl.IOMuxerFunc[rune, []any](func(args ...rune) ([]any, jpl.JPLError) {
			return next.Pipe(string(args[0]))
		}))

	default:
	}

	if params.Optional {
		return nil, nil
	}
	return nil, library.ThrowAny(library.NewTypeError("cannot iterate over %s (%*<100v)", string(t), value))
}

// { optional: boolean }
func (opaIter) Map(runtime jpl.JPLRuntime, params jpl.JPLSelectorParams) (result definition.JPLSelectorParams, err jpl.JPLError) {
	result.Optional = params.Optional
	return
}
