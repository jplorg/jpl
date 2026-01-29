package program

import (
	"github.com/jplorg/jpl/go/v2/definition"
	"github.com/jplorg/jpl/go/v2/jpl"
	"github.com/jplorg/jpl/go/v2/library"
)

type opuMultiplication struct{}

// { pipe: [op] }
func (opuMultiplication) OP(runtime jpl.JPLRuntime, input any, target any, params definition.JPLAssignmentParams, scope jpl.JPLRuntimeScope, next jpl.JPLPiper) ([]any, jpl.JPLError) {
	return runtime.ExecuteInstructions(definition.Pipe{{
		OP: definition.OP_CALCULATION,
		Params: definition.JPLInstructionParams{
			Pipe:       constant(target),
			Operations: []definition.JPLOperation{{OP: definition.OPM_MULTIPLICATION, Params: definition.JPLOperationParams{By: params.Pipe}}},
		},
	}}, []any{input}, scope, library.NewPiperWithScope(next))
}

// { pipe: function }
func (opuMultiplication) Map(runtime jpl.JPLRuntime, params jpl.JPLAssignmentParams) (result definition.JPLAssignmentParams, err jpl.JPLError) {
	result.Pipe = call(params.Pipe)
	return
}
