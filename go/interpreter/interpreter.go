package interpreter

import (
	definition "github.com/jplorg/jpl/go/v2/definition"
	jpl "github.com/jplorg/jpl/go/v2/jpl"
	program "github.com/jplorg/jpl/go/v2/program"
)

var defaultOptions = jpl.JPLInterpreterOptions{}

func NewInterpreter(options *jpl.JPLInterpreterConfig) jpl.JPLInterpreter {
	if options == nil {
		options = new(jpl.JPLInterpreterConfig)
	}

	return &interpreter{
		options:        jpl.ApplyInterpreterDefaults(options.Interpreter, defaultOptions),
		programOptions: options.Program,
		runtimeOptions: options.Runtime,
	}
}

type interpreter struct {
	options        jpl.JPLInterpreterOptions
	programOptions jpl.JPLProgramOptions
	runtimeOptions jpl.JPLRuntimeOptions
}

func (i *interpreter) Options() jpl.JPLInterpreterOptions {
	return i.options
}

func (i *interpreter) Parse(source string, options *jpl.JPLInterpreterConfig) (jpl.JPLProgram, jpl.JPLError) {
	if options == nil {
		options = new(jpl.JPLInterpreterConfig)
	}

	instructions, err := i.ParseInstructions(source)
	if err != nil {
		return nil, err
	}

	definition := definition.JPLDefinition{
		Version:      definition.DEFINITION_VERSION,
		Instructions: instructions,
	}

	return program.NewProgram(definition, &jpl.JPLProgramConfig{
		Program: jpl.ApplyProgramDefaults(options.Program, i.programOptions),
		Runtime: jpl.ApplyRuntimeDefaults(options.Runtime, i.runtimeOptions),
	})
}

func (i *interpreter) ParseInstructions(source string) (definition.Pipe, jpl.JPLSyntaxError) {
	_, instructions, err := parseEntrypoint(source, 0, &ParserContext{Interpreter: i})
	return instructions, err
}
