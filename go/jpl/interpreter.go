package jpl

import (
	"github.com/jplorg/jpl/go/definition"
)

type JPLInterpreterConfig struct {
	Interpreter JPLInterpreterOptions
	Program     JPLProgramOptions
	Runtime     JPLRuntimeOptions
}

type JPLInterpreterOptions struct{}

func ApplyInterpreterDefaults(options JPLInterpreterOptions, defaults JPLInterpreterOptions) (result JPLInterpreterOptions) {
	return
}

// JPL interpreter
type JPLInterpreter interface {
	// Return the interpreter's options
	Options() JPLInterpreterOptions

	// Parse the specified source program string into a reusable JPLProgram instance
	Parse(source string, options *JPLInterpreterConfig) (JPLProgram, JPLError)

	// Parse the specified source program string
	ParseInstructions(source string) (definition.Pipe, JPLSyntaxError)
}
