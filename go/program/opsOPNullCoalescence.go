package program

import (
	"github.com/jplorg/jpl/go/v2/definition"
	"github.com/jplorg/jpl/go/v2/jpl"
	"github.com/jplorg/jpl/go/v2/library"
)

type opNullCoalescence struct{}

// { pipes: [[op]] }
func (opNullCoalescence) OP(runtime jpl.JPLRuntime, input any, params definition.JPLInstructionParams, scope jpl.JPLRuntimeScope, next jpl.JPLScopedPiper) ([]any, jpl.JPLError) {
	var iter func(from int) ([]any, jpl.JPLError)
	iter = func(from int) ([]any, jpl.JPLError) {
		if err := scope.Signal().CheckHealth(); err != nil {
			return nil, err
		}

		if from >= len(params.Pipes) {
			return next.Pipe(nil, scope)
		}

		pipe := params.Pipes[from]

		return runtime.ExecuteInstructions(pipe, []any{input}, scope, jpl.JPLScopedPiperFunc(func(output any, _ jpl.JPLRuntimeScope) ([]any, jpl.JPLError) {
			t, err := library.Type(output)
			if err != nil {
				return nil, err
			}

			switch t {
			case jpl.JPLT_NULL:
				return iter(from + 1)

			default:
				return next.Pipe(output, scope)
			}
		}))
	}

	return iter(0)
}

// { pipes: [function] }
func (opNullCoalescence) Map(runtime jpl.JPLRuntime, params jpl.JPLInstructionParams) (result definition.JPLInstructionParams, err jpl.JPLError) {
	if result.Pipes, err = library.MuxOne([][]jpl.JPLFunc{params.Pipes}, jpl.IOMuxerFunc[jpl.JPLFunc, definition.Pipe](func(args ...jpl.JPLFunc) (definition.Pipe, jpl.JPLError) {
		return call(args[0]), nil
	})); err != nil {
		return
	}
	return
}
