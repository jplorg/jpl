package program

import (
	"strings"

	"github.com/2manyvcos/jpl/go/definition"
	"github.com/2manyvcos/jpl/go/jpl"
	"github.com/2manyvcos/jpl/go/library"
)

type opmSubtraction struct{}

// { by: [op] }
func (opmSubtraction) OP(runtime jpl.JPLRuntime, input any, target any, params definition.JPLOperationParams, scope jpl.JPLRuntimeScope, next jpl.JPLPiper) ([]any, jpl.JPLError) {
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
					return a.(float64) - b.(float64), nil
				}

			case jpl.JPLT_ARRAY:
				if tb == jpl.JPLT_ARRAY {
					va := a.([]any)
					vb := b.([]any)
					la := len(va)
					lb := len(vb)
					if la == 0 || lb == 0 {
						return a, nil
					}
					filtered := make([]any, 0, len(va))
				L:
					for _, v := range va {
						for _, entry := range vb {
							equals, err := library.Equals(v, entry)
							if err != nil {
								return nil, err
							}
							if equals {
								continue L
							}
						}
						filtered = append(filtered, v)
					}
					if la == len(filtered) {
						return a, nil
					}
					return filtered, nil
				}

			case jpl.JPLT_STRING:
				if tb == jpl.JPLT_STRING {
					return strings.ReplaceAll(a.(string), b.(string), ""), nil
				}

			case jpl.JPLT_OBJECT:
				switch tb {
				case jpl.JPLT_ARRAY:
					va := a.(map[string]any)
					vb := b.([]any)
					if len(vb) == 0 {
						return va, nil
					}
					entries := make([]*library.ObjectEntry[any], 0, len(va))
					for k, v := range va {
						for _, entry := range vb {
							equals, err := library.Equals(v, entry)
							if err != nil {
								return nil, err
							}
							if equals {
								entries = append(entries, &library.ObjectEntry[any]{Key: k, NoValue: true})
							}
						}
					}
					return library.ApplyObject(va, entries), nil

				case jpl.JPLT_STRING:
					return library.ApplyObject(a.(map[string]any), []*library.ObjectEntry[any]{{Key: b.(string), NoValue: true}}), nil

				default:
				}

			default:
			}

			return nil, library.ThrowAny(library.NewTypeError("%s (%*<100v) and %s (%*<100v) cannot be subtracted", string(ta), a, string(tb), b))
		}))
		if err != nil {
			return nil, err
		}
		return next.Pipe(alteredValue)
	}))
}

// { by: function }
func (opmSubtraction) Map(runtime jpl.JPLRuntime, params jpl.JPLOperationParams) (result definition.JPLOperationParams, err jpl.JPLError) {
	result.By = call(params.By)
	return
}
