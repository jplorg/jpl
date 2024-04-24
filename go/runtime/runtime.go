package runtime

import (
	"github.com/2manyvcos/jpl/go/config"
	"github.com/2manyvcos/jpl/go/jpl"
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
