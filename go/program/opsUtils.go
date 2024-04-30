package program

import (
	"github.com/jplorg/jpl/go/definition"
	"github.com/jplorg/jpl/go/jpl"
)

// Create OP_CONSTANT operation for the specified normalized value
func constant(value any) definition.Pipe {
	return definition.Pipe{{OP: definition.OP_CONSTANT, Params: definition.JPLInstructionParams{Value: value}}}
}

// Create OPA_FUNCTION operation for the specified function
func call(function jpl.JPLFunc) definition.Pipe {
	return definition.Pipe{
		{
			OP: definition.OP_ACCESS,
			Params: definition.JPLInstructionParams{
				Pipe:      constant(function),
				Selectors: []definition.JPLSelector{{OP: definition.OPA_FUNCTION, Params: definition.JPLSelectorParams{Args: nil, Bound: false, Optional: false}}},
			},
		},
	}
}
