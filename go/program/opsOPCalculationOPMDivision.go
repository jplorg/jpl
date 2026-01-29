package program

import (
	"strings"

	"github.com/jplorg/jpl/go/v2/definition"
	"github.com/jplorg/jpl/go/v2/jpl"
	"github.com/jplorg/jpl/go/v2/library"
)

type opmDivision struct{}

// { by: [op] }
func (opmDivision) OP(runtime jpl.JPLRuntime, input any, target any, params definition.JPLOperationParams, scope jpl.JPLRuntimeScope, next jpl.JPLPiper) ([]any, jpl.JPLError) {
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
					va := a.(float64)
					vb := b.(float64)
					if vb == 0 {
						return nil, library.ThrowAny(library.NewZeroDivisionError("%s (%*<100v) cannot be divided by zero", string(ta), a))
					}
					return va / vb, nil
				}

			case jpl.JPLT_STRING:
				if tb == jpl.JPLT_STRING {
					parts := strings.Split(a.(string), b.(string))
					result := make([]any, len(parts))
					for i, v := range parts {
						result[i] = v
					}
					return result, nil
				}

			default:
			}

			return nil, library.ThrowAny(library.NewTypeError("%s (%*<100v) and %s (%*<100v) cannot be added together", string(ta), a, string(tb), b))
		}))
		if err != nil {
			return nil, err
		}
		return next.Pipe(alteredValue)
	}))
}

// { by: function }
func (opmDivision) Map(runtime jpl.JPLRuntime, params jpl.JPLOperationParams) (result definition.JPLOperationParams, err jpl.JPLError) {
	result.By = call(params.By)
	return
}
