package program

import (
	definition "github.com/jplorg/jpl/go/v2/definition"
	jpl "github.com/jplorg/jpl/go/v2/jpl"
	library "github.com/jplorg/jpl/go/v2/library"
)

type opaFunction struct{}

// { args: [[op]], bound: boolean, optional: boolean }
func (opaFunction) OP(runtime jpl.JPLRuntime, input any, target any, params definition.JPLSelectorParams, scope jpl.JPLRuntimeScope, next jpl.JPLPiper) ([]any, jpl.JPLError) {
	value, err := library.UnwrapValue(target)
	if err != nil {
		return nil, err
	}
	t, err := library.Type(value)
	if err != nil {
		return nil, err
	}
	switch t {
	case jpl.JPLT_FUNCTION:
		args, err := library.MuxOne([][]definition.Pipe{params.Args}, jpl.IOMuxerFunc[definition.Pipe, []any](func(args ...definition.Pipe) ([]any, jpl.JPLError) {
			arg := args[0]
			return runtime.ExecuteInstructions(arg, []any{input}, scope, nil)
		}))
		if err != nil {
			return nil, err
		}

		return library.MuxAll(args, jpl.IOMuxerFunc[any, []any](func(a ...any) ([]any, jpl.JPLError) {
			fnNext := jpl.JPLPiperFunc(func(output any) ([]any, jpl.JPLError) {
				results, err := next.Pipe(output)
				if err != nil {
					return nil, library.NewErrorEnclosure(err)
				}
				return results, nil
			})

			var results []any
			var err error
			if params.Bound {
				if len(a) == 0 {
					results, err = value.(jpl.JPLFunc).Call(runtime, scope.Signal(), fnNext, nil)
				} else {
					results, err = value.(jpl.JPLFunc).Call(runtime, scope.Signal(), fnNext, a[0], a[1:]...)
				}
			} else {
				results, err = value.(jpl.JPLFunc).Call(runtime, scope.Signal(), fnNext, input, a...)
			}
			if err != nil {
				if errorEnclosure, ok := err.(jpl.JPLErrorEnclosure); ok {
					return nil, errorEnclosure.JPLEnclosedError()
				}
				return nil, library.AdaptError(err)
			}
			return results, nil
		}))

	default:
	}

	if params.Optional {
		return nil, nil
	}
	return nil, library.ThrowAny(library.NewTypeError("cannot execute %s (%*<100v)", string(t), value))
}

// { args: [function], bound: boolean, optional: boolean }
func (opaFunction) Map(runtime jpl.JPLRuntime, params jpl.JPLSelectorParams) (result definition.JPLSelectorParams, err jpl.JPLError) {
	if result.Args, err = library.MuxOne([][]jpl.JPLFunc{params.Args}, jpl.IOMuxerFunc[jpl.JPLFunc, definition.Pipe](func(args ...jpl.JPLFunc) (definition.Pipe, jpl.JPLError) {
		return call(args[0]), nil
	})); err != nil {
		return
	}
	result.Bound = params.Bound
	result.Optional = params.Optional
	return
}
