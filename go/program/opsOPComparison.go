package program

import (
	"github.com/2manyvcos/jpl/go/definition"
	"github.com/2manyvcos/jpl/go/jpl"
	"github.com/2manyvcos/jpl/go/library"
)

type opComparison struct{}

// { pipe: [op], comparisons: [opc] }
func (opComparison) OP(runtime jpl.JPLRuntime, input any, params definition.JPLInstructionParams, scope jpl.JPLRuntimeScope, next jpl.JPLScopedPiper) ([]any, jpl.JPLError) {
	var iter func(from int, value any) ([]any, jpl.JPLError)
	iter = func(from int, value any) ([]any, jpl.JPLError) {
		if err := scope.Signal().CheckHealth(); err != nil {
			return nil, err
		}

		if from >= len(params.Comparisons) {
			return next.Pipe(value, scope)
		}

		comparison := params.Comparisons[from]
		operator, ok := opcs[comparison.OP]
		if !ok {
			return nil, library.NewFatalError("invalid OPC '" + string(comparison.OP) + "'")
		}

		return operator.OP(runtime, input, value, comparison.Params, scope, jpl.JPLPiperFunc(func(output any) ([]any, jpl.JPLError) {
			return iter(from+1, output)
		}))
	}

	return runtime.ExecuteInstructions(params.Pipe, []any{input}, scope, jpl.JPLScopedPiperFunc(func(output any, _ jpl.JPLRuntimeScope) ([]any, jpl.JPLError) {
		return iter(0, output)
	}))
}

// { pipe: function, comparisons: [opc] }
func (opComparison) Map(runtime jpl.JPLRuntime, params jpl.JPLInstructionParams) (result definition.JPLInstructionParams, err jpl.JPLError) {
	result.Pipe = call(params.Pipe)
	if result.Comparisons, err = library.MuxOne([][]jpl.JPLComparison{params.Comparisons}, jpl.IOMuxerFunc[jpl.JPLComparison, definition.JPLComparison](func(args ...jpl.JPLComparison) (result definition.JPLComparison, err jpl.JPLError) {
		comparison := args[0]
		operator, ok := opcs[comparison.OP]
		if !ok {
			err = library.NewFatalError("invalid OPC '" + string(comparison.OP) + "'")
			return
		}

		result.OP = comparison.OP
		if result.Params, err = operator.Map(runtime, comparison.Params); err != nil {
			return
		}
		return
	})); err != nil {
		return
	}
	return
}

type JPLOPCHandler = jpl.JPLOPSubHandler[definition.JPLComparisonParams, jpl.JPLComparisonParams]

var opcs = map[definition.JPLOPC]JPLOPCHandler{
	definition.OPC_EQUAL:        opcEqual{},
	definition.OPC_GREATER:      opcGreater{},
	definition.OPC_GREATEREQUAL: opcGreaterEqual{},
	definition.OPC_LESS:         opcLess{},
	definition.OPC_LESSEQUAL:    opcLessEqual{},
	definition.OPC_UNEQUAL:      opcUnequal{},
}
