package ops

import (
	"github.com/2manyvcos/jpl/go/definition"
	"github.com/2manyvcos/jpl/go/jpl"
)

type opConstantFalse struct{}

// {}
func (opConstantFalse) OP(runtime jpl.JPLRuntime, input any, params definition.JPLInstructionParams, scope jpl.JPLRuntimeScope, next jpl.JPLScopedPiper) ([]any, jpl.JPLError) {
	return next.Pipe(false, scope)
}

// {}
func (opConstantFalse) Map(runtime jpl.JPLRuntime, params jpl.JPLInstructionParams) (definition.JPLInstructionParams, jpl.JPLError) {
	return definition.JPLInstructionParams{}, nil
}
