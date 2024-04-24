package interpreter

import (
	"github.com/2manyvcos/jpl/go/options"
)

var defaultOptions = options.JPLInterpreterOptions{}

type Options struct {
	Interpreter options.JPLInterpreterOptions
	Program     options.JPLProgramOptions
	Runtime     options.JPLRuntimeOptions
}

// JPL interpreter
type JPLInterpreter interface {
	// Parse the specified source program string into a reusable JPLProgram instance
	Parse(source string, config *Options) (JPLProgram, error)

	// Parse the specified source program string
	ParseInstructions(source string) ([]JPLInstruction, error)
}

func NewInterpreter(config *Options) JPLInterpreter {
	if config == nil {
		config = new(Options)
	}

	return &interpreter{
		Options:        options.ApplyInterpreterDefaults(config.Interpreter, defaultOptions),
		ProgramOptions: config.Program,
		RuntimeOptions: config.Runtime,
	}
}

type interpreter struct {
	Options        options.JPLInterpreterOptions
	ProgramOptions options.JPLProgramOptions
	RuntimeOptions options.JPLRuntimeOptions
}

func (i *interpreter) Parse(source string, config *Options) (JPLProgram, error) {
	panic("TODO:")
}

func (i *interpreter) ParseInstructions(source string) ([]JPLInstruction, error) {
	panic("TODO:")
}
