package jpl

import (
	"github.com/2manyvcos/jpl/go/builtins"
	"github.com/2manyvcos/jpl/go/config"
	"github.com/2manyvcos/jpl/go/interpreter"
	"github.com/2manyvcos/jpl/go/program"
)

var Options = &interpreter.Options{
	Interpreter: config.JPLInterpreterOptions{},

	Program: config.JPLProgramOptions{},

	Runtime: config.JPLRuntimeOptions{
		Vars: builtins.GetBuiltins(),
	},
}

func Parse(source string, options *interpreter.Options) (program.JPLProgram, error) {
	if options == nil {
		options = new(interpreter.Options)
	}
	interpreter := interpreter.NewInterpreter(&interpreter.Options{
		Interpreter: config.ApplyInterpreterDefaults(options.Interpreter, Options.Interpreter),
		Program:     config.ApplyProgramDefaults(options.Program, Options.Program),
		Runtime:     config.ApplyRuntimeDefaults(options.Runtime, Options.Runtime),
	})
	return interpreter.Parse(source, nil)
}

func Run(source string, inputs []any, options *interpreter.Options) ([]any, error) {
	program, err := Parse(source, options)
	if err != nil {
		return nil, err
	}
	return program.Run(inputs, nil)
}
