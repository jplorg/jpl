package ops

import (
	"github.com/2manyvcos/jpl/go/definition"
	"github.com/2manyvcos/jpl/go/jpl"
	"github.com/2manyvcos/jpl/go/library"
)

type opFunctionDefinition struct{}

// { argNames: [string], pipe: [op] }
func (opFunctionDefinition) OP(runtime jpl.JPLRuntime, input any, params definition.JPLInstructionParams, scope jpl.JPLRuntimeScope, next jpl.JPLScopedPiper) ([]any, jpl.JPLError) {
	return next.Pipe(library.ScopedFunction(params.ArgNames, params.Pipe, scope), scope)
}

// { argNames: [string], pipe: function }
func (opFunctionDefinition) Map(runtime jpl.JPLRuntime, params jpl.JPLInstructionParams) (definition.JPLInstructionParams, jpl.JPLError) {
	return definition.JPLInstructionParams{
		ArgNames: params.ArgNames,
		Pipe:     call(params.Pipe),
	}, nil
}
