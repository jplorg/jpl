package jpl

import (
	"github.com/2manyvcos/jpl/go/config"
	"github.com/2manyvcos/jpl/go/definition"
	"github.com/2manyvcos/jpl/go/library"
)

// JPL runtime
type JPLRuntime interface {
	// Return the runtime's options
	Options() config.JPLRuntimeOptions

	// Return the runtime's program
	Program() JPLProgram

	// Create a new orphan scope
	CreateScope(presets *library.RuntimeScopeConfig) library.RuntimeScope

	// Execute a new dedicated program
	Execute(inputs []any) ([]any, error)

	// Execute the specified instructions
	ExecuteInstructions(instructions definition.Pipe, inputs []any, scope library.RuntimeScope, next library.JPLPiper) ([]any, error)

	// Execute the specified OP
	OP(op definition.JPLOP, params map[string]any, inputs []any, scope library.RuntimeScope, next library.JPLPiper) ([]any, error)
}
