package program

import (
	"github.com/2manyvcos/jpl/go/definition"
	"github.com/2manyvcos/jpl/go/jpl"
)

type opConstantTrue struct{}

// {}
func (opConstantTrue) OP(runtime jpl.JPLRuntime, input any, params definition.JPLInstructionParams, scope jpl.JPLRuntimeScope, next jpl.JPLScopedPiper) ([]any, jpl.JPLError) {
	return next.Pipe(true, scope)
}

// {}
func (opConstantTrue) Map(runtime jpl.JPLRuntime, params jpl.JPLInstructionParams) (result definition.JPLInstructionParams, err jpl.JPLError) {
	return
}
