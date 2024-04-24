package ops

import (
	"github.com/2manyvcos/jpl/go/definition"
	"github.com/2manyvcos/jpl/go/jpl"
	"github.com/2manyvcos/jpl/go/library"
)

type opOutputConcat struct{}

// { pipes: [[op]] }
func (opOutputConcat) OP(runtime jpl.JPLRuntime, input any, params definition.JPLInstructionParams, scope jpl.JPLRuntimeScope, next jpl.JPLScopedPiper) ([]any, jpl.JPLError) {
	return library.MuxAll([][]definition.Pipe{params.Pipes}, jpl.IOMuxerFunc[definition.Pipe, []any](func(args ...definition.Pipe) ([]any, jpl.JPLError) {
		return runtime.ExecuteInstructions(args[0], []any{input}, scope, jpl.JPLScopedPiperFunc(func(output any, _ jpl.JPLRuntimeScope) ([]any, jpl.JPLError) {
			return next.Pipe(output, scope)
		}))
	}))
}

// { pipes: [function] }
func (opOutputConcat) Map(runtime jpl.JPLRuntime, params jpl.JPLInstructionParams) (definition.JPLInstructionParams, jpl.JPLError) {
	pipes, err := library.MuxOne([][]jpl.JPLFunc{params.Pipes}, jpl.IOMuxerFunc[jpl.JPLFunc, definition.Pipe](func(args ...jpl.JPLFunc) (definition.Pipe, jpl.JPLError) {
		return call(args[0]), nil
	}))
	if err != nil {
		return definition.JPLInstructionParams{}, err
	}

	return definition.JPLInstructionParams{
		Pipes: pipes,
	}, nil
}
