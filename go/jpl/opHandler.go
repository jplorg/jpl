package jpl

import "github.com/2manyvcos/jpl/go/definition"

type JPLOPHandler interface {
	OP(runtime JPLRuntime, input any, params definition.JPLInstructionParams, scope JPLRuntimeScope, next JPLScopedPiper) ([]any, JPLError)

	Map(runtime JPLRuntime, params JPLInstructionParams) (definition.JPLInstructionParams, JPLError)
}

type JPLOPSubHandler[DefinitionParams any, JPLParams any] interface {
	OP(runtime JPLRuntime, input any, target any, params DefinitionParams, scope JPLRuntimeScope, next JPLPiper) ([]any, JPLError)

	Map(runtime JPLRuntime, params JPLParams) (DefinitionParams, JPLError)
}
