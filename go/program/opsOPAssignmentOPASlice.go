package program

import (
	"github.com/jplorg/jpl/go/definition"
	"github.com/jplorg/jpl/go/jpl"
	"github.com/jplorg/jpl/go/library"
)

type sliceRange struct{ Start, End any }

func shallowCompareArrays[Type any](a, b []Type) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if !library.IsSame(v, b[i]) {
			return false
		}
	}
	return true
}

type opaAssignSlice struct{}

// { from: [op], to: [op], optional: boolean }
func (opaAssignSlice) OP(runtime jpl.JPLRuntime, input any, target any, params definition.JPLSelectorParams, scope jpl.JPLRuntimeScope, next jpl.JPLPiper) ([]any, jpl.JPLError) {
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
	ranges, err := library.MuxOne([][]any{unwrappedFroms, unwrappedTos}, jpl.IOMuxerFunc[any, sliceRange](func(args ...any) (sliceRange, jpl.JPLError) {
		return sliceRange{Start: args[0], End: args[1]}, nil
	}))
	if err != nil {
		return nil, err
	}

	var iter func(from int, source any) ([]any, jpl.JPLError)
	iter = func(from int, source any) ([]any, jpl.JPLError) {
		if err := scope.Signal().CheckHealth(); err != nil {
			return nil, err
		}

		if from >= len(ranges) {
			return []any{source}, nil
		}

		v, err := library.UnwrapValue(source)
		if err != nil {
			return nil, err
		}
		tv, err := library.Type(source)
		if err != nil {
			return nil, err
		}
		r := ranges[from]
		ts, err := library.Type(r.Start)
		if err != nil {
			return nil, err
		}
		te, err := library.Type(r.End)
		if err != nil {
			return nil, err
		}
		switch tv {
		case jpl.JPLT_ARRAY:
			if (ts == jpl.JPLT_NUMBER || ts == jpl.JPLT_NULL) && (te == jpl.JPLT_NUMBER || te == jpl.JPLT_NULL) {
				v := v.([]any)
				l := len(v)
				vs := 0
				if ts == jpl.JPLT_NUMBER {
					vs = int(r.Start.(float64))
				}
				ve := l
				if te == jpl.JPLT_NUMBER {
					ve = int(r.End.(float64))
				}
				values, err := next.Pipe(library.SubSlice(v, vs, ve))
				if err != nil {
					return nil, err
				}
				return library.MuxAll([][]any{values}, jpl.IOMuxerFunc[any, []any](func(args ...any) ([]any, jpl.JPLError) {
					output := args[0]
					if _, ok := output.(unchanged); ok {
						return iter(from+1, source)
					}
					alteredValue, err := library.AlterValue(source, jpl.JPLModifierFunc(func(value any) (any, jpl.JPLError) {
						v := value.([]any)
						result, err := library.UnwrapValue(output)
						if err != nil {
							return nil, err
						}
						tr, err := library.Type(result)
						if err != nil {
							return nil, err
						}
						switch tr {
						case jpl.JPLT_NULL, jpl.JPLT_ARRAY:
							s := vs
							e := ve
							if s >= 0 {
								s = min(l, s)
							} else {
								s = max(0, l+s)
							}
							if e >= 0 {
								e = min(l, e)
							} else {
								e = max(0, l+e)
							}
							e = max(s, e)
							var r []any
							if tr == jpl.JPLT_ARRAY {
								r = result.([]any)
							}
							if shallowCompareArrays(v[s:e], r) {
								return value, nil
							}
							lr := len(r)
							c := make([]any, s+lr+(l-e))
							copy(c, v[:s])
							copy(c[s:], r)
							copy(c[s+lr:], v[e:])
							return c, nil

						default:
						}

						return nil, library.ThrowAny(library.NewTypeError("cannot assign %s (%*<100v) to slice of %s (%*<100v)", string(tr), result, string(tv), value))
					}))
					if err != nil {
						return nil, err
					}
					return iter(from+1, alteredValue)
				}))
			}

		case jpl.JPLT_STRING:
			if (ts == jpl.JPLT_NUMBER || ts == jpl.JPLT_NULL) && (te == jpl.JPLT_NUMBER || te == jpl.JPLT_NULL) {
				chars := []rune(v.(string))
				l := len(chars)
				vs := 0
				if ts == jpl.JPLT_NUMBER {
					vs = int(r.Start.(float64))
				}
				ve := l
				if te == jpl.JPLT_NUMBER {
					ve = int(r.End.(float64))
				}
				values, err := next.Pipe(string(library.SubSlice(chars, vs, ve)))
				if err != nil {
					return nil, err
				}
				return library.MuxAll([][]any{values}, jpl.IOMuxerFunc[any, []any](func(args ...any) ([]any, jpl.JPLError) {
					output := args[0]
					if _, ok := output.(unchanged); ok {
						return iter(from+1, source)
					}
					alteredValue, err := library.AlterValue(source, jpl.JPLModifierFunc(func(value any) (any, jpl.JPLError) {
						result, err := library.UnwrapValue(output)
						if err != nil {
							return nil, err
						}
						tr, err := library.Type(result)
						if err != nil {
							return nil, err
						}
						switch tr {
						case jpl.JPLT_NULL, jpl.JPLT_STRING:
							s := vs
							e := ve
							if s >= 0 {
								s = min(l, s)
							} else {
								s = max(0, l+s)
							}
							if e >= 0 {
								e = min(l, e)
							} else {
								e = max(0, l+e)
							}
							e = max(s, e)
							var r []rune
							if tr == jpl.JPLT_STRING {
								r = []rune(result.(string))
							}
							if shallowCompareArrays(chars[s:e], r) {
								return value, nil
							}
							lr := len(r)
							c := make([]rune, s+lr+(l-e))
							copy(c, chars[:s])
							copy(c[s:], r)
							copy(c[s+lr:], chars[e:])
							return string(c), nil

						default:
						}

						return nil, library.ThrowAny(library.NewTypeError("cannot assign %s (%*<100v) to slice of %s (%*<100v)", string(tr), result, string(tv), value))
					}))
					if err != nil {
						return nil, err
					}
					return iter(from+1, alteredValue)
				}))
			}

		default:
		}

		if params.Optional {
			return iter(from+1, source)
		}
		return nil, library.ThrowAny(library.NewTypeError("cannot slice %s (%*<100v) with %s (%*<100v) and %s (%*<100v) (assignment)", string(tv), v, string(ts), r.Start, string(te), r.End))
	}

	values, err := iter(0, target)
	if err != nil {
		return nil, err
	}
	return library.MuxOne([][]any{values}, jpl.IOMuxerFunc[any, any](func(args ...any) (any, jpl.JPLError) {
		output := args[0]
		t, err := library.Type(output)
		if err != nil {
			return nil, err
		}
		if t == jpl.JPLT_NULL {
			return unchanged{}, nil
		}
		return output, nil
	}))
}

// { from: function, to: function, optional: boolean }
func (opaAssignSlice) Map(runtime jpl.JPLRuntime, params jpl.JPLSelectorParams) (result definition.JPLSelectorParams, err jpl.JPLError) {
	result.From = call(params.From)
	result.To = call(params.To)
	result.Optional = params.Optional
	return
}
