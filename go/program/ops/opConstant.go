package ops

import (
	"github.com/2manyvcos/jpl/go/definition"
	"github.com/2manyvcos/jpl/go/jpl"
	"github.com/2manyvcos/jpl/go/library"
)

type opConstant struct{}

// { value: any }
func (opConstant) OP(runtime jpl.JPLRuntime, input any, params definition.JPLInstructionParams, scope jpl.JPLRuntimeScope, next jpl.JPLScopedPiper) ([]any, jpl.JPLError) {
	return next.Pipe(params.Value, scope)
}

// { value: any }
func (opConstant) Map(runtime jpl.JPLRuntime, params jpl.JPLInstructionParams) (definition.JPLInstructionParams, jpl.JPLError) {
	normalized, err := library.NormalizeValue(params.Value)
	if err != nil {
		return definition.JPLInstructionParams{}, err
	}
	return definition.JPLInstructionParams{
		Value: normalized,
	}, nil
}
