package jpl

import "github.com/2manyvcos/jpl/go/definition"

type JPLOPHandler interface {
	OP(runtime JPLRuntime, input any, params definition.JPLInstructionParams, scope JPLRuntimeScope, next JPLScopedPiper) ([]any, JPLError)

	Map(runtime JPLRuntime, params map[string]any) (definition.JPLInstructionParams, JPLError)
}
