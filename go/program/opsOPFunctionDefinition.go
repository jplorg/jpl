package program

import (
	definition "github.com/jplorg/jpl/go/v2/definition"
	jpl "github.com/jplorg/jpl/go/v2/jpl"
	library "github.com/jplorg/jpl/go/v2/library"
)

type opFunctionDefinition struct{}

// { argNames: [string], pipe: [op] }
func (opFunctionDefinition) OP(runtime jpl.JPLRuntime, input any, params definition.JPLInstructionParams, scope jpl.JPLRuntimeScope, next jpl.JPLScopedPiper) ([]any, jpl.JPLError) {
	return next.Pipe(library.ScopedFunction(params.ArgNames, params.Pipe, scope), scope)
}

// { argNames: [string], pipe: function }
func (opFunctionDefinition) Map(runtime jpl.JPLRuntime, params jpl.JPLInstructionParams) (result definition.JPLInstructionParams, err jpl.JPLError) {
	return definition.JPLInstructionParams{
		ArgNames: params.ArgNames,
		Pipe:     call(params.Pipe),
	}, nil
}
