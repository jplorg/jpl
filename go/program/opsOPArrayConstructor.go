package program

import (
	"github.com/jplorg/jpl/go/v2/definition"
	"github.com/jplorg/jpl/go/v2/jpl"
)

type opArrayConstructor struct{}

// { pipe: [op] }
func (opArrayConstructor) OP(runtime jpl.JPLRuntime, input any, params definition.JPLInstructionParams, scope jpl.JPLRuntimeScope, next jpl.JPLScopedPiper) ([]any, jpl.JPLError) {
	outputs, err := runtime.ExecuteInstructions(params.Pipe, []any{input}, scope, nil)
	if err != nil {
		return nil, err
	}
	if outputs == nil {
		outputs = []any{}
	}
	return next.Pipe(outputs, scope)
}

// { pipe: function }
func (opArrayConstructor) Map(runtime jpl.JPLRuntime, params jpl.JPLInstructionParams) (result definition.JPLInstructionParams, err jpl.JPLError) {
	return definition.JPLInstructionParams{
		Pipe: call(params.Pipe),
	}, nil
}
