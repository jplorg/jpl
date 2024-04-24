package jpl

import (
	"github.com/2manyvcos/jpl/go/definition"
)

type JPLRuntimeConfig struct {
	Runtime JPLRuntimeOptions
}

type JPLRuntimeOptions struct {
	Vars map[string]any

	AdjustResult JPLScopedPiper
}

func ApplyRuntimeDefaults(options JPLRuntimeOptions, defaults JPLRuntimeOptions) (result JPLRuntimeOptions) {
	result.Vars = make(map[string]any)
	for k, v := range defaults.Vars {
		result.Vars[k] = v
	}
	for k, v := range options.Vars {
		result.Vars[k] = v
	}

	return
}

// JPL runtime
type JPLRuntime interface {
	// Return the runtime's options
	Options() JPLRuntimeOptions

	// Return the runtime's program
	Program() JPLProgram

	// Create a new orphan scope
	CreateScope(presets *JPLRuntimeScopeConfig) JPLRuntimeScope

	// Execute a new dedicated program
	Execute(inputs []any) ([]any, JPLError)

	// Execute the specified instructions
	ExecuteInstructions(instructions definition.Pipe, inputs []any, scope JPLRuntimeScope, next JPLScopedPiper) ([]any, JPLError)

	// Execute the specified OP
	OP(op definition.JPLOP, params map[string]any, inputs []any, scope JPLRuntimeScope, next JPLScopedPiper) ([]any, JPLError)
}
