package jpl

import (
	"github.com/2manyvcos/jpl/go/config"
	"github.com/2manyvcos/jpl/go/definition"
)

// JPL program
type JPLProgram interface {
	// Run the program with the provided inputs and runtime options.
	// The program throws a JPLExecutionError for runtime failures.
	// Other errors may be thrown when execution fails.
	Run(inputs []any, options *config.JPLProgramConfig) ([]any, error)

	// Return the program's definition.
	// The definition can be serialized as JSON to be reused in other JPL implementations.
	Definition() definition.JPLDefinition
}
