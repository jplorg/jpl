package jpl

import (
	"github.com/2manyvcos/jpl/go/interpreter"
	"github.com/2manyvcos/jpl/go/builtins"
	"github.com/2manyvcos/jpl/go/options"
)


var DefaultOptions = &options.JPLOptions{
	Interpreter: options.JPLInterpreterOptions{},

	Program: options.JPLProgramOptions{},

	Runtime: options.JPLRuntimeOptions{
		Vars: builtins.GetBuiltins(),
	},
}


func Parse(source string, options *JPLOptions) (JPLProgram, error) {
	interpreter := &interpreter.JPLInterpreter{
		Options: interpreter.ApplyInterpreterDefaults()
	}

	panic("TODO:")
}

func Run(source string, inputs []any, options *JPLOptions) ([]any, error) {
	program, err := Parse(source, options)
	if err != nil {
		return nil, err
	}

	return program.Run(inputs)
}
