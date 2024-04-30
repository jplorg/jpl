package program

import (
	"github.com/jplorg/jpl/go/definition"
	"github.com/jplorg/jpl/go/jpl"
	"github.com/jplorg/jpl/go/library"
)

type opmAddition struct{}

// { by: [op] }
func (opmAddition) OP(runtime jpl.JPLRuntime, input any, target any, params definition.JPLOperationParams, scope jpl.JPLRuntimeScope, next jpl.JPLPiper) ([]any, jpl.JPLError) {
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

			if ta == jpl.JPLT_NULL || tb == jpl.JPLT_NULL {
				if ta != jpl.JPLT_NULL {
					return a, nil
				}
				if tb != jpl.JPLT_NULL {
					return b, nil
				}
				return nil, nil
			}

			switch ta {
			case jpl.JPLT_NUMBER:
				if tb == jpl.JPLT_NUMBER {
					return a.(float64) + b.(float64), nil
				}

			case jpl.JPLT_ARRAY:
				if tb == jpl.JPLT_ARRAY {
					va := a.([]any)
					vb := b.([]any)
					la := len(va)
					lb := len(vb)
					if lb == 0 {
						return a, nil
					}
					if la == 0 {
						return b, nil
					}
					result := make([]any, la+lb)
					copy(result, va)
					copy(result[la:], vb)
					return result, nil
				}

			case jpl.JPLT_STRING:
				if tb == jpl.JPLT_STRING {
					return a.(string) + b.(string), nil
				}

			case jpl.JPLT_OBJECT:
				if tb == jpl.JPLT_OBJECT {
					return library.ApplyObject(a.(map[string]any), library.ObjectEntries(b.(map[string]any))), nil
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
func (opmAddition) Map(runtime jpl.JPLRuntime, params jpl.JPLOperationParams) (result definition.JPLOperationParams, err jpl.JPLError) {
	result.By = call(params.By)
	return
}
