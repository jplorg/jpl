package jpl

import (
	"github.com/2manyvcos/jpl/go/config"
	"github.com/2manyvcos/jpl/go/definition"
	"github.com/2manyvcos/jpl/go/library"
)

// JPL interpreter
type JPLInterpreter interface {
	// Return the interpreter's options
	Options() config.JPLInterpreterOptions

	// Parse the specified source program string into a reusable JPLProgram instance
	Parse(source string, options *config.JPLInterpreterConfig) (JPLProgram, library.JPLError)

	// Parse the specified source program string
	ParseInstructions(source string) (definition.Pipe, library.JPLSyntaxError)
}
