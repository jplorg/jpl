package interpreter

import (
	"github.com/2manyvcos/jpl/go/config"
	"github.com/2manyvcos/jpl/go/definition"
	"github.com/2manyvcos/jpl/go/jpl"
	"github.com/2manyvcos/jpl/go/program"
)

var defaultOptions = config.JPLInterpreterOptions{}

func NewInterpreter(options *config.JPLInterpreterConfig) jpl.JPLInterpreter {
	if options == nil {
		options = new(config.JPLInterpreterConfig)
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

func (i *interpreter) Parse(source string, options *config.JPLInterpreterConfig) (jpl.JPLProgram, error) {
	if options == nil {
		options = new(config.JPLInterpreterConfig)
	}

	instructions, err := i.ParseInstructions(source)
	if err != nil {
		return nil, err
	}

	definition := definition.JPLDefinition{
		Version:      definition.DEFINITION_VERSION,
		Instructions: instructions,
	}

	return program.NewProgram(definition, &config.JPLProgramConfig{
		Program: config.ApplyProgramDefaults(options.Program, i.ProgramOptions),
		Runtime: config.ApplyRuntimeDefaults(options.Runtime, i.RuntimeOptions),
	})
}

func (i *interpreter) ParseInstructions(source string) (definition.Pipe, error) {
	_, instructions, err := parseEntrypoint(source, 0, &ParserContext{Interpreter: i})
	return instructions, err
}
