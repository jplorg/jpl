package gojpl

import (
	"github.com/2manyvcos/jpl/go/builtins"
	"github.com/2manyvcos/jpl/go/config"
	"github.com/2manyvcos/jpl/go/interpreter"
	"github.com/2manyvcos/jpl/go/jpl"
)

var Options = &config.JPLInterpreterConfig{
	Interpreter: config.JPLInterpreterOptions{},

	Program: config.JPLProgramOptions{},

	Runtime: config.JPLRuntimeOptions{
		Vars: builtins.GetBuiltins(),
	},
}

func Parse(source string, options *config.JPLInterpreterConfig) (jpl.JPLProgram, error) {
	if options == nil {
		options = new(config.JPLInterpreterConfig)
	}
	interpreter := interpreter.NewInterpreter(&config.JPLInterpreterConfig{
		Interpreter: config.ApplyInterpreterDefaults(options.Interpreter, Options.Interpreter),
		Program:     config.ApplyProgramDefaults(options.Program, Options.Program),
		Runtime:     config.ApplyRuntimeDefaults(options.Runtime, Options.Runtime),
	})
	return interpreter.Parse(source, nil)
}

func Run(source string, inputs []any, options *config.JPLInterpreterConfig) ([]any, error) {
	program, err := Parse(source, options)
	if err != nil {
		return nil, err
	}
	return program.Run(inputs, nil)
}
