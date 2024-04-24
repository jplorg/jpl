package gojpl

import (
	"github.com/2manyvcos/jpl/go/builtins"
	"github.com/2manyvcos/jpl/go/interpreter"
	"github.com/2manyvcos/jpl/go/jpl"
)

var Options = &jpl.JPLInterpreterConfig{
	Interpreter: jpl.JPLInterpreterOptions{},

	Program: jpl.JPLProgramOptions{},

	Runtime: jpl.JPLRuntimeOptions{
		Vars: builtins.GetBuiltins(),
	},
}

func Parse(source string, options *jpl.JPLInterpreterConfig) (jpl.JPLProgram, jpl.JPLError) {
	if options == nil {
		options = new(jpl.JPLInterpreterConfig)
	}
	interpreter := interpreter.NewInterpreter(&jpl.JPLInterpreterConfig{
		Interpreter: jpl.ApplyInterpreterDefaults(options.Interpreter, Options.Interpreter),
		Program:     jpl.ApplyProgramDefaults(options.Program, Options.Program),
		Runtime:     jpl.ApplyRuntimeDefaults(options.Runtime, Options.Runtime),
	})
	return interpreter.Parse(source, nil)
}

func Run(source string, inputs []any, options *jpl.JPLInterpreterConfig) ([]any, jpl.JPLError) {
	program, err := Parse(source, options)
	if err != nil {
		return nil, err
	}
	return program.Run(inputs, nil)
}
