package runtime

import (
	"github.com/2manyvcos/jpl/go/definition"
	"github.com/2manyvcos/jpl/go/jpl"
	"github.com/2manyvcos/jpl/go/library"
)

var defaultOptions = jpl.JPLRuntimeOptions{}

func NewRuntime(program jpl.JPLProgram, options *jpl.JPLRuntimeConfig) jpl.JPLRuntime {
	if options == nil {
		options = new(jpl.JPLRuntimeConfig)
	}

	return &runtime{
		options: jpl.ApplyRuntimeDefaults(options.Runtime, defaultOptions),

		program: program,
	}
}

type runtime struct {
	options jpl.JPLRuntimeOptions

	program jpl.JPLProgram
}

func (r *runtime) Options() jpl.JPLRuntimeOptions {
	return r.options
}

func (r *runtime) Program() jpl.JPLProgram {
	return r.program
}

func (r *runtime) CreateScope(presets *jpl.JPLRuntimeScopeConfig) jpl.JPLRuntimeScope {
	return library.NewRuntimeScope(presets)
}

func (r *runtime) Execute(inputs []any) ([]any, jpl.JPLError) {
	panic("TODO:")
}

func (r *runtime) ExecuteInstructions(instructions definition.Pipe, inputs []any, scope jpl.JPLRuntimeScope, next jpl.JPLPiper) ([]any, jpl.JPLError) {
	panic("TODO:")
}

func (r *runtime) OP(op definition.JPLOP, params map[string]any, inputs []any, scope jpl.JPLRuntimeScope, next jpl.JPLPiper) ([]any, jpl.JPLError) {
	panic("TODO:")
}
