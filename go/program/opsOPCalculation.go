package program

import (
	definition "github.com/jplorg/jpl/go/v2/definition"
	jpl "github.com/jplorg/jpl/go/v2/jpl"
	library "github.com/jplorg/jpl/go/v2/library"
)

type opCalculation struct{}

// { pipe: [op], operations: [opm] }
func (opCalculation) OP(runtime jpl.JPLRuntime, input any, params definition.JPLInstructionParams, scope jpl.JPLRuntimeScope, next jpl.JPLScopedPiper) ([]any, jpl.JPLError) {
	var iter func(from int, value any) ([]any, jpl.JPLError)
	iter = func(from int, value any) ([]any, jpl.JPLError) {
		if err := scope.Signal().CheckHealth(); err != nil {
			return nil, err
		}

		if from >= len(params.Operations) {
			return next.Pipe(value, scope)
		}

		operation := params.Operations[from]
		operator, ok := opms[operation.OP]
		if !ok {
			return nil, library.NewFatalError("invalid OPM '" + string(operation.OP) + "'")
		}

		return operator.OP(runtime, input, value, operation.Params, scope, jpl.JPLPiperFunc(func(output any) ([]any, jpl.JPLError) {
			return iter(from+1, output)
		}))
	}

	return runtime.ExecuteInstructions(params.Pipe, []any{input}, scope, jpl.JPLScopedPiperFunc(func(output any, _ jpl.JPLRuntimeScope) ([]any, jpl.JPLError) { return iter(0, output) }))
}

// { pipe: function, operations: [opm] }
func (opCalculation) Map(runtime jpl.JPLRuntime, params jpl.JPLInstructionParams) (result definition.JPLInstructionParams, err jpl.JPLError) {
	result.Pipe = call(params.Pipe)
	if result.Operations, err = library.MuxOne([][]jpl.JPLOperation{params.Operations}, jpl.IOMuxerFunc[jpl.JPLOperation, definition.JPLOperation](func(args ...jpl.JPLOperation) (result definition.JPLOperation, err jpl.JPLError) {
		operation := args[0]
		operator, ok := opms[operation.OP]
		if !ok {
			err = library.NewFatalError("invalid OPM '" + string(operation.OP) + "'")
			return
		}

		result.OP = operation.OP
		if result.Params, err = operator.Map(runtime, operation.Params); err != nil {
			return
		}
		return
	})); err != nil {
		return
	}
	return
}

type JPLOPMHandler = jpl.JPLOPSubHandler[definition.JPLOperationParams, jpl.JPLOperationParams]

var opms = map[definition.JPLOPM]JPLOPMHandler{
	definition.OPM_ADDITION:       opmAddition{},
	definition.OPM_DIVISION:       opmDivision{},
	definition.OPM_MULTIPLICATION: opmMultiplication{},
	definition.OPM_REMAINDER:      opmRemainder{},
	definition.OPM_SUBTRACTION:    opmSubtraction{},
}
