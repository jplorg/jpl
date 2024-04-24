package ops

import (
	"github.com/2manyvcos/jpl/go/definition"
	"github.com/2manyvcos/jpl/go/jpl"
)

type opArrayConstructor struct{}

// { pipe: [op] }
func (opArrayConstructor) OP(runtime jpl.JPLRuntime, input any, params definition.JPLInstructionParams, scope jpl.JPLRuntimeScope, next jpl.JPLScopedPiper) ([]any, jpl.JPLError) {
	outputs, err := runtime.ExecuteInstructions(params.Pipe, []any{input}, scope, nil)
	if err != nil {
		return nil, err
	}
	return next.Pipe(outputs, scope)
}

// { pipe: function }
func (opArrayConstructor) Map(runtime jpl.JPLRuntime, params jpl.JPLInstructionParams) (definition.JPLInstructionParams, jpl.JPLError) {
	return definition.JPLInstructionParams{
		Pipe: call(params.Pipe),
	}, nil
}
