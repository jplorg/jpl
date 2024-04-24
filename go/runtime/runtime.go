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
		Options: config.ApplyRuntimeDefaults(options.Runtime, defaultOptions),

		Program: program,
	}
}

type runtime struct {
	Options config.JPLRuntimeOptions

	Program jpl.JPLProgram
}
