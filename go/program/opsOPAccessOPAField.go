package program

import (
	"github.com/2manyvcos/jpl/go/definition"
	"github.com/2manyvcos/jpl/go/jpl"
	"github.com/2manyvcos/jpl/go/library"
)

type opaField struct{}

// { pipe: [op], optional: boolean }
func (opaField) OP(runtime jpl.JPLRuntime, input any, target any, params definition.JPLSelectorParams, scope jpl.JPLRuntimeScope, next jpl.JPLPiper) ([]any, jpl.JPLError) {
	value, err := library.UnwrapValue(target)
	if err != nil {
		return nil, err
	}
	tv, err := library.Type(value)
	if err != nil {
		return nil, err
	}

	return runtime.ExecuteInstructions(params.Pipe, []any{input}, scope, jpl.JPLScopedPiperFunc(func(output any, _ jpl.JPLRuntimeScope) ([]any, jpl.JPLError) {
		field, err := library.UnwrapValue(output)
		if err != nil {
			return nil, err
		}
		tf, err := library.Type(field)
		if err != nil {
			return nil, err
		}
		switch tv {
		case jpl.JPLT_NULL:
			if tf == jpl.JPLT_STRING || tf == jpl.JPLT_NUMBER {
				return next.Pipe(nil)
			}

		case jpl.JPLT_OBJECT:
			if tf == jpl.JPLT_STRING {
				return next.Pipe(value.(map[string]any)[field.(string)])
			}

		case jpl.JPLT_ARRAY:
			if tf == jpl.JPLT_NUMBER {
				i := int(field.(float64))
				v := value.([]any)
				l := len(v)
				vi := i
				if i < 0 {
					vi = l + i
				}
				if vi >= 0 && vi < l {
					return next.Pipe(v[vi])
				} else {
					return next.Pipe(nil)
				}
			}

		case jpl.JPLT_STRING:
			if tf == jpl.JPLT_NUMBER {
				i := int(field.(float64))
				chars := []rune(value.(string))
				l := len(chars)
				vi := i
				if i < 0 {
					vi = l + i
				}
				if vi >= 0 && vi < l {
					return next.Pipe(string(chars[vi]))
				} else {
					return next.Pipe(nil)
				}
			}

		default:
		}

		if params.Optional {
			return nil, nil
		}
		return nil, library.ThrowAny(library.NewTypeError("cannot access field of %s (%*<100v) with %s (%*<100v)", string(tv), value, string(tf), field))
	}))
}

// { pipe: function, optional: boolean }
func (opaField) Map(runtime jpl.JPLRuntime, params jpl.JPLSelectorParams) (result definition.JPLSelectorParams, err jpl.JPLError) {
	result.Pipe = call(params.Pipe)
	result.Optional = params.Optional
	return
}
