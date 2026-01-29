package program

import (
	"github.com/jplorg/jpl/go/v2/definition"
	"github.com/jplorg/jpl/go/v2/jpl"
)

type opConstantNull struct{}

// {}
func (opConstantNull) OP(runtime jpl.JPLRuntime, input any, params definition.JPLInstructionParams, scope jpl.JPLRuntimeScope, next jpl.JPLScopedPiper) ([]any, jpl.JPLError) {
	return next.Pipe(nil, scope)
}

// {}
func (opConstantNull) Map(runtime jpl.JPLRuntime, params jpl.JPLInstructionParams) (result definition.JPLInstructionParams, err jpl.JPLError) {
	return
}
