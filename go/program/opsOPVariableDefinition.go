package program

import (
	"github.com/jplorg/jpl/go/definition"
	"github.com/jplorg/jpl/go/jpl"
)

type opVariableDefinition struct{}

// { name: string, pipe: [op] }
func (opVariableDefinition) OP(runtime jpl.JPLRuntime, input any, params definition.JPLInstructionParams, scope jpl.JPLRuntimeScope, next jpl.JPLScopedPiper) ([]any, jpl.JPLError) {
	return runtime.ExecuteInstructions(params.Pipe, []any{input}, scope, jpl.JPLScopedPiperFunc(func(output any, scope jpl.JPLRuntimeScope) ([]any, jpl.JPLError) {
		return next.Pipe(input, scope.Next(&jpl.JPLRuntimeScopeConfig{Vars: map[string]any{params.Name: output}}))
	}))
}

// { name: string, pipe: function }
func (opVariableDefinition) Map(runtime jpl.JPLRuntime, params jpl.JPLInstructionParams) (result definition.JPLInstructionParams, err jpl.JPLError) {
	return definition.JPLInstructionParams{
		Name: params.Name,
		Pipe: call(params.Pipe),
	}, nil
}
