package program

import (
	"github.com/2manyvcos/jpl/go/config"
	"github.com/2manyvcos/jpl/go/definition"
)

var defaultOptions = config.JPLProgramOptions{}

type Options struct {
	Program config.JPLProgramOptions
	Runtime config.JPLRuntimeOptions
}

// JPL program
type JPLProgram interface {
	// Run the program with the provided inputs and runtime options.
	// The program throws a JPLExecutionError for runtime failures.
	// Other errors may be thrown when execution fails.
	Run(inputs []any, options *Options) ([]any, error)
}

func NewProgram(instructions definition.Pipe, options *Options) JPLProgram {
	if options == nil {
		options = new(Options)
	}
	return &program{
		Options: config.ApplyProgramDefaults(options.Program, defaultOptions),

		Instructions: instructions,
		OPs:          OPs,

		RuntimeOptions: options.Runtime,
	}
}

type program struct {
	Options config.JPLProgramOptions

	Instructions definition.Pipe
	OPs          map[string]JPLOP

	RuntimeOptions config.JPLRuntimeOptions
}

func (p *program) Run(inputs []any, options *Options) ([]any, error) {
	panic("TODO:")
}
