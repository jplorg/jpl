package interpreter

import (
	"github.com/2manyvcos/jpl/go/config"
	"github.com/2manyvcos/jpl/go/library/definition"
	"github.com/2manyvcos/jpl/go/program"
)

var defaultOptions = config.JPLInterpreterOptions{}

type Options struct {
	Interpreter config.JPLInterpreterOptions
	Program     config.JPLProgramOptions
	Runtime     config.JPLRuntimeOptions
}

// JPL interpreter
type JPLInterpreter interface {
	// Parse the specified source program string into a reusable JPLProgram instance
	Parse(source string, options *Options) (program.JPLProgram, error)

	// Parse the specified source program string
	ParseInstructions(source string) (definition.Pipe, error)
}

func NewInterpreter(options *Options) JPLInterpreter {
	if options == nil {
		options = new(Options)
	}
	return &interpreter{
		Options: config.ApplyInterpreterDefaults(options.Interpreter, defaultOptions),

		ProgramOptions: options.Program,
		RuntimeOptions: options.Runtime,
	}
}

type interpreter struct {
	Options config.JPLInterpreterOptions

	ProgramOptions config.JPLProgramOptions
	RuntimeOptions config.JPLRuntimeOptions
}

func (i *interpreter) Parse(source string, options *Options) (program.JPLProgram, error) {
	if options == nil {
		options = new(Options)
	}

	instructions, err := i.ParseInstructions(source)
	if err != nil {
		return nil, err
	}

	definition := definition.JPLDefinition{
		Version:      definition.DEFINITION_VERSION,
		Instructions: instructions,
	}

	return program.NewProgram(definition, &program.Options{
		Program: config.ApplyProgramDefaults(options.Program, i.ProgramOptions),
		Runtime: config.ApplyRuntimeDefaults(options.Runtime, i.RuntimeOptions),
	})
}

func (i *interpreter) ParseInstructions(source string) (definition.Pipe, error) {
	_, instructions, err := parseEntrypoint(source, 0, &ParserContext{Interpreter: i})
	return instructions, err
}
