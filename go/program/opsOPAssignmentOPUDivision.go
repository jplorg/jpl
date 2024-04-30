package program

import (
	"github.com/jplorg/jpl/go/definition"
	"github.com/jplorg/jpl/go/jpl"
	"github.com/jplorg/jpl/go/library"
)

type opuDivision struct{}

// { pipe: [op] }
func (opuDivision) OP(runtime jpl.JPLRuntime, input any, target any, params definition.JPLAssignmentParams, scope jpl.JPLRuntimeScope, next jpl.JPLPiper) ([]any, jpl.JPLError) {
	return runtime.ExecuteInstructions(definition.Pipe{{
		OP: definition.OP_CALCULATION,
		Params: definition.JPLInstructionParams{
			Pipe:       constant(target),
			Operations: []definition.JPLOperation{{OP: definition.OPM_DIVISION, Params: definition.JPLOperationParams{By: params.Pipe}}},
		},
	}}, []any{input}, scope, library.NewPiperWithScope(next))
}

// { pipe: function }
func (opuDivision) Map(runtime jpl.JPLRuntime, params jpl.JPLAssignmentParams) (result definition.JPLAssignmentParams, err jpl.JPLError) {
	result.Pipe = call(params.Pipe)
	return
}
