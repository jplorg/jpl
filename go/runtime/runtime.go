package runtime

import (
	"github.com/2manyvcos/jpl/go/config"
	"github.com/2manyvcos/jpl/go/definition"
	"github.com/2manyvcos/jpl/go/jpl"
	"github.com/2manyvcos/jpl/go/library"
)

var defaultOptions = config.JPLRuntimeOptions{}

func NewRuntime(program jpl.JPLProgram, options *config.JPLRuntimeConfig) jpl.JPLRuntime {
	if options == nil {
		options = new(config.JPLRuntimeConfig)
	}

	return &runtime{
		options: config.ApplyRuntimeDefaults(options.Runtime, defaultOptions),

		program: program,
	}
}

type runtime struct {
	options config.JPLRuntimeOptions

	program jpl.JPLProgram
}

func (r *runtime) Options() config.JPLRuntimeOptions {
	return r.options
}

func (r *runtime) Program() jpl.JPLProgram {
	return r.program
}

func (r *runtime) CreateScope(presets *library.RuntimeScopeConfig) library.RuntimeScope {
	return library.NewRuntimeScope(presets)
}

func (r *runtime) Execute(inputs []any) ([]any, library.JPLError) {
	panic("TODO:")
}

func (r *runtime) ExecuteInstructions(instructions definition.Pipe, inputs []any, scope library.RuntimeScope, next library.JPLPiper) ([]any, library.JPLError) {
	panic("TODO:")
}

func (r *runtime) OP(op definition.JPLOP, params map[string]any, inputs []any, scope library.RuntimeScope, next library.JPLPiper) ([]any, library.JPLError) {
	panic("TODO:")
}
