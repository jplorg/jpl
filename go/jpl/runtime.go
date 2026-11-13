package jpl

import (
	"github.com/jplorg/jpl/go/v2/definition"
)

type JPLRuntimeConfig struct {
	Runtime JPLRuntimeOptions
}

type JPLRuntimeOptions struct {
	Vars map[string]any

	AdjustResult JPLScopedPiper

	Signal <-chan struct{}
}

func ApplyRuntimeDefaults(options JPLRuntimeOptions, defaults JPLRuntimeOptions) (result JPLRuntimeOptions) {
	result.Vars = make(map[string]any)
	for k, v := range defaults.Vars {
		result.Vars[k] = v
	}
	for k, v := range options.Vars {
		result.Vars[k] = v
	}

	if options.AdjustResult != nil {
		result.AdjustResult = options.AdjustResult
	} else {
		result.AdjustResult = defaults.AdjustResult
	}

	if options.Signal != nil || defaults.Signal != nil {
		if options.Signal != nil && defaults.Signal != nil {
			signal := make(chan struct{})
			go func() {
				select {
				case <-options.Signal:
				case <-defaults.Signal:
				}

				close(signal)
			}()
			result.Signal = signal
		} else if options.Signal != nil {
			result.Signal = options.Signal
		} else {
			result.Signal = defaults.Signal
		}
	}

	return
}

// JPL runtime
type JPLRuntime interface {
	// Return the runtime's options
	Options() JPLRuntimeOptions

	// Return the runtime's program
	Program() JPLProgram

	// Return the runtime's root scope
	RootScope() JPLRuntimeScope

	// Create a new orphan scope
	CreateScope(presets *JPLRuntimeScopeConfig) JPLRuntimeScope

	// Execute a new dedicated program
	Execute(inputs []any) ([]any, JPLError)

	// Execute the specified instructions
	ExecuteInstructions(instructions definition.Pipe, inputs []any, scope JPLRuntimeScope, next JPLScopedPiper) ([]any, JPLError)

	// Execute the specified OP
	OP(op definition.JPLOP, params JPLInstructionParams, inputs []any, scope JPLRuntimeScope, next JPLScopedPiper) ([]any, JPLError)
}
