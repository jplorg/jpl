package program

import (
	"github.com/2manyvcos/jpl/go/definition"
	"github.com/2manyvcos/jpl/go/jpl"
	"github.com/2manyvcos/jpl/go/library"
)

type opOr struct{}

// { pipes: [[op]] }
func (opOr) OP(runtime jpl.JPLRuntime, input any, params definition.JPLInstructionParams, scope jpl.JPLRuntimeScope, next jpl.JPLScopedPiper) ([]any, jpl.JPLError) {
	var iter func(from int) ([]any, jpl.JPLError)
	iter = func(from int) ([]any, jpl.JPLError) {
		if err := scope.Signal().CheckHealth(); err != nil {
			return nil, err
		}

		if from >= len(params.Pipes) {
			return next.Pipe(false, scope)
		}

		pipe := params.Pipes[from]

		return runtime.ExecuteInstructions(pipe, []any{input}, scope, jpl.JPLScopedPiperFunc(func(output any, _ jpl.JPLRuntimeScope) ([]any, jpl.JPLError) {
			if truthy, err := library.Truthy(output); err != nil {
				return nil, err
			} else if truthy {
				return next.Pipe(true, scope)
			}

			return iter(from + 1)
		}))
	}

	return iter(0)
}

// { pipes: [function] }
func (opOr) Map(runtime jpl.JPLRuntime, params jpl.JPLInstructionParams) (result definition.JPLInstructionParams, err jpl.JPLError) {
	if result.Pipes, err = library.MuxOne([][]jpl.JPLFunc{params.Pipes}, jpl.IOMuxerFunc[jpl.JPLFunc, definition.Pipe](func(args ...jpl.JPLFunc) (definition.Pipe, jpl.JPLError) {
		return call(args[0]), nil
	})); err != nil {
		return
	}
	return
}
