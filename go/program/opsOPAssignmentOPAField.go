package program

import (
	"strings"

	"github.com/2manyvcos/jpl/go/definition"
	"github.com/2manyvcos/jpl/go/jpl"
	"github.com/2manyvcos/jpl/go/library"
)

type opaAssignField struct{}

// { pipe: [op], optional: boolean }
func (opaAssignField) OP(runtime jpl.JPLRuntime, input any, target any, params definition.JPLSelectorParams, scope jpl.JPLRuntimeScope, next jpl.JPLPiper) ([]any, jpl.JPLError) {
	var items any
	vt, err := library.UnwrapValue(target)
	if err != nil {
		return nil, err
	}
	tt, err := library.Type(target)
	if err != nil {
		return nil, err
	}
	switch tt {
	case jpl.JPLT_NULL, jpl.JPLT_OBJECT, jpl.JPLT_ARRAY:
		items = vt

	case jpl.JPLT_STRING:
		items = strings.Split(vt.(string), "")

	default:
		if params.Optional {
			return []any{unchanged{}}, nil
		}
		return nil, library.ThrowAny(library.NewTypeError("cannot access fields of %s (%*<100v) (assignment)", string(tt), vt))
	}

	fields, err := runtime.ExecuteInstructions(params.Pipe, []any{input}, scope, jpl.JPLScopedPiperFunc(func(output any, _ jpl.JPLRuntimeScope) ([]any, jpl.JPLError) {
		return []any{output}, nil
	}))
	if err != nil {
		return nil, err
	}

	var iter func(from int, source any) ([]any, jpl.JPLError)
	iter = func(from int, source any) ([]any, jpl.JPLError) {
		if err := scope.Signal().CheckHealth(); err != nil {
			return nil, err
		}

		if from >= len(fields) {
			return []any{source}, nil
		}

		var t jpl.JPLDataType
		if tt == jpl.JPLT_STRING {
			t = tt
		} else {
			t, err = library.Type(source)
			if err != nil {
				return nil, err
			}
		}
		field, err := library.UnwrapValue(fields[from])
		if err != nil {
			return nil, err
		}
		tf, err := library.Type(field)
		if err != nil {
			return nil, err
		}
		switch t {
		case jpl.JPLT_NULL:
			switch tf {
			case jpl.JPLT_STRING:
				values, err := next.Pipe(nil)
				if err != nil {
					return nil, err
				}
				return library.MuxAll([][]any{values}, jpl.IOMuxerFunc[any, []any](func(args ...any) ([]any, jpl.JPLError) {
					output := args[0]
					if _, ok := output.(unchanged); ok {
						return iter(from+1, source)
					}
					return iter(from+1, map[string]any{field.(string): output})
				}))

			case jpl.JPLT_NUMBER:
				i := int(field.(float64))
				values, err := next.Pipe(nil)
				if err != nil {
					return nil, err
				}
				return library.MuxAll([][]any{values}, jpl.IOMuxerFunc[any, []any](func(args ...any) ([]any, jpl.JPLError) {
					output := args[0]
					if _, ok := output.(unchanged); ok {
						return iter(from+1, source)
					}
					return iter(from+1, library.ApplyArray([]any{}, []*library.ArrayEntry[any]{{Index: i, Value: output}}, nil))
				}))

			default:
			}

		case jpl.JPLT_OBJECT:
			if tf == jpl.JPLT_STRING {
				s := source.(map[string]any)
				f := field.(string)
				item := s[f]
				values, err := next.Pipe(item)
				if err != nil {
					return nil, err
				}
				return library.MuxAll([][]any{values}, jpl.IOMuxerFunc[any, []any](func(args ...any) ([]any, jpl.JPLError) {
					output := args[0]
					if _, ok := output.(unchanged); ok {
						return iter(from+1, source)
					}
					return iter(from+1, library.ApplyObject(s, []*library.ObjectEntry[any]{{Key: f, Value: output}}))
				}))
			}

		case jpl.JPLT_ARRAY:
			if tf == jpl.JPLT_NUMBER {
				i := int(field.(float64))
				s := source.([]any)
				l := len(s)
				vi := i
				if i < 0 {
					vi = l + i
				}
				var item any
				if vi >= 0 && vi < l {
					item = s[vi]
				}
				values, err := next.Pipe(item)
				if err != nil {
					return nil, err
				}
				return library.MuxAll([][]any{values}, jpl.IOMuxerFunc[any, []any](func(args ...any) ([]any, jpl.JPLError) {
					output := args[0]
					if _, ok := output.(unchanged); ok {
						return iter(from+1, source)
					}
					return iter(from+1, library.ApplyArray(s, []*library.ArrayEntry[any]{{Index: i, Value: output}}, nil))
				}))
			}

		case jpl.JPLT_STRING:
			if tf == jpl.JPLT_NUMBER {
				i := int(field.(float64))
				s := source.([]string)
				l := len(s)
				vi := i
				if i < 0 {
					vi = l + i
				}
				var item any
				if vi >= 0 && vi < l {
					item = string(s[vi])
				}
				values, err := next.Pipe(item)
				if err != nil {
					return nil, err
				}
				return library.MuxAll([][]any{values}, jpl.IOMuxerFunc[any, []any](func(args ...any) ([]any, jpl.JPLError) {
					output := args[0]
					if _, ok := output.(unchanged); ok {
						return iter(from+1, source)
					}
					r, err := library.UnwrapValue(output)
					if err != nil {
						return nil, err
					}
					tr, err := library.Type(r)
					if err != nil {
						return nil, err
					}
					switch tr {
					case jpl.JPLT_NULL, jpl.JPLT_STRING:
						value := " "
						if tr == jpl.JPLT_STRING {
							value = r.(string)
						}
						return iter(from+1, library.ApplyArray(s, []*library.ArrayEntry[string]{{Index: i, Value: value}}, " "))

					default:
					}

					return nil, library.ThrowAny(library.NewTypeError("cannot assign %s (%*<100v) to string (%*<100v)", string(tr), r, strings.Join(s, "")))
				}))
			}

		default:
		}

		if params.Optional {
			return iter(from+1, source)
		}
		var v any
		if t == jpl.JPLT_STRING {
			v = strings.Join(source.([]string), "")
		} else {
			v = source
		}
		return nil, library.ThrowAny(library.NewTypeError("cannot access field of %s (%*<100v) with %s (%*<100v) (assignment)", string(t), v, string(tf), field))
	}

	values, err := iter(0, items)
	if err != nil {
		return nil, err
	}
	return library.MuxOne([][]any{values}, jpl.IOMuxerFunc[any, any](func(args ...any) (any, jpl.JPLError) {
		results := args[0]
		if library.IsSame(items, results) {
			return target, nil
		}

		var t jpl.JPLDataType
		if tt == jpl.JPLT_STRING {
			t = tt
		} else {
			t, err = library.Type(results)
			if err != nil {
				return nil, err
			}
		}
		switch t {
		case jpl.JPLT_NULL:
			return unchanged{}, nil

		case jpl.JPLT_OBJECT, jpl.JPLT_ARRAY:
			return library.AlterValue(target, jpl.JPLModifierFunc(func(any) (any, jpl.JPLError) { return results, nil }))

		case jpl.JPLT_STRING:
			return library.AlterValue(target, jpl.JPLModifierFunc(func(any) (any, jpl.JPLError) { return strings.Join(results.([]string), ""), nil }))

		default:
		}

		return unchanged{}, nil
	}))
}

// { pipe: function, optional: boolean }
func (opaAssignField) Map(runtime jpl.JPLRuntime, params jpl.JPLSelectorParams) (result definition.JPLSelectorParams, err jpl.JPLError) {
	result.Pipe = call(params.Pipe)
	result.Optional = params.Optional
	return
}
