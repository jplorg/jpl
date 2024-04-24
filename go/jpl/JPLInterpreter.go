package jpl

import (
	"github.com/2manyvcos/jpl/go/config"
	"github.com/2manyvcos/jpl/go/definition"
)

// JPL interpreter
type JPLInterpreter interface {
	// Parse the specified source program string into a reusable JPLProgram instance
	Parse(source string, options *config.JPLInterpreterConfig) (JPLProgram, error)

	// Parse the specified source program string
	ParseInstructions(source string) (definition.Pipe, error)
}
