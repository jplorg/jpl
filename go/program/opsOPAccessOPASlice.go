package program

import (
	"github.com/jplorg/jpl/go/definition"
	"github.com/jplorg/jpl/go/jpl"
	"github.com/jplorg/jpl/go/library"
)

type opaSlice struct{}

// { from: [op], to: [op], optional: boolean }
func (opaSlice) OP(runtime jpl.JPLRuntime, input any, target any, params definition.JPLSelectorParams, scope jpl.JPLRuntimeScope, next jpl.JPLPiper) ([]any, jpl.JPLError) {
	value, err := library.UnwrapValue(target)
	if err != nil {
		return nil, err
	}
	tv, err := library.Type(value)
	if err != nil {
		return nil, err
	}

	froms, err := runtime.ExecuteInstructions(params.From, []any{input}, scope, nil)
	if err != nil {
		return nil, err
	}
	tos, err := runtime.ExecuteInstructions(params.To, []any{input}, scope, nil)
	if err != nil {
		return nil, err
	}

	unwrappedFroms, err := library.UnwrapValues(froms, "")
	if err != nil {
		return nil, err
	}
	unwrappedTos, err := library.UnwrapValues(tos, "")
	if err != nil {
		return nil, err
	}
	return library.MuxAll([][]any{unwrappedFroms, unwrappedTos}, jpl.IOMuxerFunc[any, []any](func(args ...any) ([]any, jpl.JPLError) {
		from := args[0]
		to := args[1]
		tf, err := library.Type(from)
		if err != nil {
			return nil, err
		}
		tt, err := library.Type(to)
		if err != nil {
			return nil, err
		}
		switch tv {
		case jpl.JPLT_ARRAY:
			if (tf == jpl.JPLT_NUMBER || tf == jpl.JPLT_NULL) && (tt == jpl.JPLT_NUMBER || tt == jpl.JPLT_NULL) {
				v := value.([]any)
				l := len(v)
				vf := 0
				if tf == jpl.JPLT_NUMBER {
					vf = int(from.(float64))
				}
				vt := l
				if tt == jpl.JPLT_NUMBER {
					vt = int(to.(float64))
				}
				return next.Pipe(library.SubSlice(v, vf, vt))
			}

		case jpl.JPLT_STRING:
			if (tf == jpl.JPLT_NUMBER || tf == jpl.JPLT_NULL) && (tt == jpl.JPLT_NUMBER || tt == jpl.JPLT_NULL) {
				chars := []rune(value.(string))
				l := len(chars)
				vf := 0
				if tf == jpl.JPLT_NUMBER {
					vf = int(from.(float64))
				}
				vt := l
				if tt == jpl.JPLT_NUMBER {
					vt = int(to.(float64))
				}
				return next.Pipe(string(library.SubSlice(chars, vf, vt)))
			}

		default:
		}

		if params.Optional {
			return nil, nil
		}
		return nil, library.ThrowAny(library.NewTypeError("cannot slice %s (%*<100v) with %s (%*<100v) and %s (%*<100v)", string(tv), value, string(tf), from, string(tt), to))
	}))
}

// { from: function, to: function, optional: boolean }
func (opaSlice) Map(runtime jpl.JPLRuntime, params jpl.JPLSelectorParams) (result definition.JPLSelectorParams, err jpl.JPLError) {
	result.From = call(params.From)
	result.To = call(params.To)
	result.Optional = params.Optional
	return
}
