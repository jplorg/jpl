package program

import (
	"strings"

	definition "github.com/jplorg/jpl/go/v2/definition"
	jpl "github.com/jplorg/jpl/go/v2/jpl"
	library "github.com/jplorg/jpl/go/v2/library"
)

type opaAssignIter struct{}

// { optional: boolean }
func (opaAssignIter) OP(runtime jpl.JPLRuntime, input any, target any, params definition.JPLSelectorParams, scope jpl.JPLRuntimeScope, next jpl.JPLPiper) ([]any, jpl.JPLError) {
	vt, err := library.UnwrapValue(target)
	if err != nil {
		return nil, err
	}
	tt, err := library.Type(vt)
	if err != nil {
		return nil, err
	}
	switch tt {
	case jpl.JPLT_NULL:
		return []any{unchanged{}}, nil

	case jpl.JPLT_OBJECT:
		items := library.ObjectEntries(vt.(map[string]any))
		outputs, err := library.MuxOne([][]*library.ObjectEntry[any]{items}, jpl.IOMuxerFunc[*library.ObjectEntry[any], []*library.ObjectEntry[any]](func(args ...*library.ObjectEntry[any]) ([]*library.ObjectEntry[any], jpl.JPLError) {
			item := args[0]
			values, err := next.Pipe(item.Value)
			if err != nil {
				return nil, err
			}
			return library.MuxOne([][]any{values}, jpl.IOMuxerFunc[any, *library.ObjectEntry[any]](func(args ...any) (*library.ObjectEntry[any], jpl.JPLError) {
				output := args[0]
				if _, ok := output.(unchanged); ok {
					return item, nil
				}
				return &library.ObjectEntry[any]{Key: item.Key, Value: output}, nil
			}))
		}))
		if err != nil {
			return nil, err
		}
		return library.MuxOne([][][]*library.ObjectEntry[any]{library.ApplyCombinations(items, outputs)}, jpl.IOMuxerFunc[[]*library.ObjectEntry[any], any](func(args ...[]*library.ObjectEntry[any]) (any, jpl.JPLError) {
			results := args[0]
			if library.IsSame(items, results) {
				return target, nil
			}
			return library.AlterValue(target, jpl.JPLModifierFunc(func(any) (any, jpl.JPLError) { return library.ObjectFromEntries(results), nil }))
		}))

	case jpl.JPLT_ARRAY:
		items := vt.([]any)
		outputs, err := library.MuxOne([][]any{items}, jpl.IOMuxerFunc[any, []any](func(args ...any) ([]any, jpl.JPLError) {
			item := args[0]
			values, err := next.Pipe(item)
			if err != nil {
				return nil, err
			}
			return library.MuxOne([][]any{values}, jpl.IOMuxerFunc[any, any](func(args ...any) (any, jpl.JPLError) {
				output := args[0]
				if _, ok := output.(unchanged); ok {
					return item, nil
				}
				return output, nil
			}))
		}))
		if err != nil {
			return nil, err
		}
		return library.MuxOne([][][]any{library.ApplyCombinations(items, outputs)}, jpl.IOMuxerFunc[[]any, any](func(args ...[]any) (any, jpl.JPLError) {
			results := args[0]
			if library.IsSame(items, results) {
				return target, nil
			}
			return library.AlterValue(target, jpl.JPLModifierFunc(func(any) (any, jpl.JPLError) { return results, nil }))
		}))

	case jpl.JPLT_STRING:
		items := strings.Split(vt.(string), "")
		outputs, err := library.MuxOne([][]string{items}, jpl.IOMuxerFunc[string, []string](func(args ...string) ([]string, jpl.JPLError) {
			item := args[0]
			values, err := next.Pipe(item)
			if err != nil {
				return nil, err
			}
			return library.MuxOne([][]any{values}, jpl.IOMuxerFunc[any, string](func(args ...any) (string, jpl.JPLError) {
				output := args[0]
				if _, ok := output.(unchanged); ok {
					return item, nil
				}
				r, err := library.UnwrapValue(output)
				if err != nil {
					return "", err
				}
				tr, err := library.Type(r)
				if err != nil {
					return "", err
				}
				switch tr {
				case jpl.JPLT_NULL, jpl.JPLT_STRING:
					value := " "
					if tr == jpl.JPLT_STRING {
						value = r.(string)
					}
					return value, nil

				default:
				}

				return "", library.ThrowAny(library.NewTypeError("cannot assign %s (%*<100v) to string (%*<100v)", string(tr), r, vt))
			}))
		}))
		if err != nil {
			return nil, err
		}
		return library.MuxOne([][][]string{library.ApplyCombinations(items, outputs)}, jpl.IOMuxerFunc[[]string, any](func(args ...[]string) (any, jpl.JPLError) {
			results := args[0]
			if library.IsSame(items, results) {
				return target, nil
			}
			return library.AlterValue(target, jpl.JPLModifierFunc(func(any) (any, jpl.JPLError) { return strings.Join(results, ""), nil }))
		}))

	default:
	}

	if params.Optional {
		return []any{unchanged{}}, nil
	}
	return nil, library.ThrowAny(library.NewTypeError("cannot iterate over %s (%*<100v) (assignment)", string(tt), vt))
}

// { optional: boolean }
func (opaAssignIter) Map(runtime jpl.JPLRuntime, params jpl.JPLSelectorParams) (result definition.JPLSelectorParams, err jpl.JPLError) {
	result.Optional = params.Optional
	return
}
