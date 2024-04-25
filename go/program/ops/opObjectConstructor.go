package ops

import (
	"github.com/2manyvcos/jpl/go/definition"
	"github.com/2manyvcos/jpl/go/jpl"
	"github.com/2manyvcos/jpl/go/library"
)

type opObjectConstructor struct{}

// { fields: [{ key: [op], value: [op], optional: boolean }] }
func (opObjectConstructor) OP(runtime jpl.JPLRuntime, input any, params definition.JPLInstructionParams, scope jpl.JPLRuntimeScope, next jpl.JPLScopedPiper) ([]any, jpl.JPLError) {
	fields, err := library.MuxOne([][]definition.JPLField{params.Fields}, jpl.IOMuxerFunc[definition.JPLField, []*library.ObjectEntry[any]](func(args ...definition.JPLField) ([]*library.ObjectEntry[any], jpl.JPLError) {
		field := args[0]
		keys, err := runtime.ExecuteInstructions(field.Key, []any{input}, scope, nil)
		if err != nil {
			return nil, err
		}
		values, err := runtime.ExecuteInstructions(field.Value, []any{input}, scope, nil)
		if err != nil {
			return nil, err
		}

		unwrappedKeys, err := library.UnwrapValues(keys, "")
		if err != nil {
			return nil, err
		}
		return library.MuxAll([][]any{unwrappedKeys, values}, jpl.IOMuxerFunc[any, []*library.ObjectEntry[any]](func(args ...any) ([]*library.ObjectEntry[any], jpl.JPLError) {
			key := args[0]
			value := args[1]
			t, err := library.Type(key)
			if err != nil {
				return nil, err
			}
			switch t {
			case jpl.JPLT_STRING:
				return []*library.ObjectEntry[any]{{Key: key.(string), Value: value}}, nil

			default:
			}

			if field.Optional {
				return nil, nil
			}
			return nil, library.ThrowAny(library.NewTypeError("cannot use %s (%*<100v) as object key", string(t), key))
		}))
	}))
	if err != nil {
		return nil, err
	}

	return library.MuxAll(fields, jpl.IOMuxerFunc[*library.ObjectEntry[any], []any](func(entries ...*library.ObjectEntry[any]) ([]any, jpl.JPLError) {
		return next.Pipe(library.ObjectFromEntries(entries), scope)
	}))
}

// { fields: [{ key: function, value: function, optional: boolean }] }
func (opObjectConstructor) Map(runtime jpl.JPLRuntime, params jpl.JPLInstructionParams) (result definition.JPLInstructionParams, err jpl.JPLError) {
	if result.Fields, err = library.MuxOne([][]jpl.JPLField{params.Fields}, jpl.IOMuxerFunc[jpl.JPLField, definition.JPLField](func(args ...jpl.JPLField) (definition.JPLField, jpl.JPLError) {
		entry := args[0]
		return definition.JPLField{
			Key:      call(entry.Key),
			Value:    call(entry.Value),
			Optional: entry.Optional,
		}, nil
	})); err != nil {
		return
	}
	return
}
