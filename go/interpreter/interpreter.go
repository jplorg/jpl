package interpreter

import (
	"github.com/2manyvcos/jpl/go/config"
	"github.com/2manyvcos/jpl/go/definition"
	"github.com/2manyvcos/jpl/go/jpl"
	"github.com/2manyvcos/jpl/go/library"
	"github.com/2manyvcos/jpl/go/program"
)

var defaultOptions = config.JPLInterpreterOptions{}

func NewInterpreter(options *config.JPLInterpreterConfig) jpl.JPLInterpreter {
	if options == nil {
		options = new(config.JPLInterpreterConfig)
	}

	return &interpreter{
		options:        config.ApplyInterpreterDefaults(options.Interpreter, defaultOptions),
		programOptions: options.Program,
		runtimeOptions: options.Runtime,
	}
}

type interpreter struct {
	options        config.JPLInterpreterOptions
	programOptions config.JPLProgramOptions
	runtimeOptions config.JPLRuntimeOptions
}

func (i *interpreter) Options() config.JPLInterpreterOptions {
	return i.options
}

func (i *interpreter) Parse(source string, options *config.JPLInterpreterConfig) (jpl.JPLProgram, library.JPLError) {
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
		Program: config.ApplyProgramDefaults(options.Program, i.programOptions),
		Runtime: config.ApplyRuntimeDefaults(options.Runtime, i.runtimeOptions),
	})
}

func (i *interpreter) ParseInstructions(source string) (definition.Pipe, library.JPLSyntaxError) {
	_, instructions, err := parseEntrypoint(source, 0, &ParserContext{Interpreter: i})
	return instructions, err
}
