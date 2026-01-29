package program

import (
	"strings"

	definition "github.com/jplorg/jpl/go/v2/definition"
	jpl "github.com/jplorg/jpl/go/v2/jpl"
	library "github.com/jplorg/jpl/go/v2/library"
)

type opmMultiplication struct{}

// { by: [op] }
func (opmMultiplication) OP(runtime jpl.JPLRuntime, input any, target any, params definition.JPLOperationParams, scope jpl.JPLRuntimeScope, next jpl.JPLPiper) ([]any, jpl.JPLError) {
	return runtime.ExecuteInstructions(params.By, []any{input}, scope, jpl.JPLScopedPiperFunc(func(by any, _ jpl.JPLRuntimeScope) ([]any, jpl.JPLError) {
		alteredValue, err := library.AlterValue(target, jpl.JPLModifierFunc(func(a any) (any, jpl.JPLError) {
			b, err := library.UnwrapValue(by)
			if err != nil {
				return nil, err
			}

			ta, err := library.Type(a)
			if err != nil {
				return nil, err
			}
			tb, err := library.Type(b)
			if err != nil {
				return nil, err
			}

			switch ta {
			case jpl.JPLT_NUMBER:
				if tb == jpl.JPLT_NUMBER {
					return a.(float64) * b.(float64), nil
				}

			case jpl.JPLT_STRING:
				if tb == jpl.JPLT_NUMBER {
					vb := b.(float64)
					if vb < 1 {
						return nil, nil
					}
					return strings.Repeat(a.(string), int(vb)), nil
				}

			case jpl.JPLT_OBJECT:
				if tb == jpl.JPLT_OBJECT {
					return library.Merge(a, b)
				}

			default:
			}

			return nil, library.ThrowAny(library.NewTypeError("%s (%*<100v) and %s (%*<100v) cannot be multiplied together", string(ta), a, string(tb), b))
		}))
		if err != nil {
			return nil, err
		}
		return next.Pipe(alteredValue)
	}))
}

// { by: function }
func (opmMultiplication) Map(runtime jpl.JPLRuntime, params jpl.JPLOperationParams) (result definition.JPLOperationParams, err jpl.JPLError) {
	result.By = call(params.By)
	return
}
