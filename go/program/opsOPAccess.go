package program

import (
	"github.com/jplorg/jpl/go/v2/definition"
	"github.com/jplorg/jpl/go/v2/jpl"
	"github.com/jplorg/jpl/go/v2/library"
)

type opAccess struct{}

// { pipe: [op], selectors: [opa] }
func (opAccess) OP(runtime jpl.JPLRuntime, input any, params definition.JPLInstructionParams, scope jpl.JPLRuntimeScope, next jpl.JPLScopedPiper) ([]any, jpl.JPLError) {
	var iter func(from int, value any) ([]any, jpl.JPLError)
	iter = func(from int, value any) ([]any, jpl.JPLError) {
		if err := scope.Signal().CheckHealth(); err != nil {
			return nil, err
		}

		if from >= len(params.Selectors) {
			return next.Pipe(value, scope)
		}

		selector := params.Selectors[from]
		operator, ok := opas[selector.OP]
		if !ok {
			return nil, library.NewFatalError("invalid OPA '" + string(selector.OP) + "'")
		}

		return operator.OP(runtime, input, value, selector.Params, scope, jpl.JPLPiperFunc(func(output any) ([]any, jpl.JPLError) {
			return iter(from+1, output)
		}))
	}

	return runtime.ExecuteInstructions(params.Pipe, []any{input}, scope, jpl.JPLScopedPiperFunc(func(output any, _ jpl.JPLRuntimeScope) ([]any, jpl.JPLError) { return iter(0, output) }))
}

// { pipe: function, selectors: [opa] }
func (opAccess) Map(runtime jpl.JPLRuntime, params jpl.JPLInstructionParams) (result definition.JPLInstructionParams, err jpl.JPLError) {
	result.Pipe = call(params.Pipe)
	if result.Selectors, err = library.MuxOne([][]jpl.JPLSelector{params.Selectors}, jpl.IOMuxerFunc[jpl.JPLSelector, definition.JPLSelector](func(args ...jpl.JPLSelector) (result definition.JPLSelector, err jpl.JPLError) {
		selector := args[0]
		operator, ok := opas[selector.OP]
		if !ok {
			err = library.NewFatalError("invalid OPA '" + string(selector.OP) + "'")
			return
		}

		result.OP = selector.OP
		if result.Params, err = operator.Map(runtime, selector.Params); err != nil {
			return
		}
		return
	})); err != nil {
		return
	}
	return
}

type JPLOPAHandler = jpl.JPLOPSubHandler[definition.JPLSelectorParams, jpl.JPLSelectorParams]

var opas = map[definition.JPLOPA]JPLOPAHandler{
	definition.OPA_FIELD:    opaField{},
	definition.OPA_FUNCTION: opaFunction{},
	definition.OPA_ITER:     opaIter{},
	definition.OPA_SLICE:    opaSlice{},
}
