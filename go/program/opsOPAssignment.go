package program

import (
	"github.com/2manyvcos/jpl/go/definition"
	"github.com/2manyvcos/jpl/go/jpl"
	"github.com/2manyvcos/jpl/go/library"
)

type unchanged struct{}

type opAssignment struct{}

// { pipe: [op], selectors: [opa], assignment: [opu] }
func (opAssignment) OP(runtime jpl.JPLRuntime, input any, params definition.JPLInstructionParams, scope jpl.JPLRuntimeScope, next jpl.JPLScopedPiper) ([]any, jpl.JPLError) {
	if params.Assignment == nil {
		params.Assignment = new(definition.JPLAssignment)
	}

	var iter func(from int, value any) ([]any, jpl.JPLError)
	iter = func(from int, value any) ([]any, jpl.JPLError) {
		if err := scope.Signal().CheckHealth(); err != nil {
			return nil, err
		}

		if from >= len(params.Selectors) {
			operator, ok := opus[params.Assignment.OP]
			if !ok {
				return nil, library.NewFatalError("invalid OPU '" + string(params.Assignment.OP) + "'")
			}

			return operator.OP(runtime, input, value, params.Assignment.Params, scope, jpl.JPLPiperFunc(func(output any) ([]any, jpl.JPLError) {
				return []any{output}, nil
			}))
		}

		selector := params.Selectors[from]
		operator, ok := opasAssign[selector.OP]
		if !ok {
			return nil, library.NewFatalError("invalid OPA '" + string(selector.OP) + "' (assignment)")
		}

		return operator.OP(runtime, input, value, selector.Params, scope, jpl.JPLPiperFunc(func(output any) ([]any, jpl.JPLError) {
			return iter(from+1, output)
		}))
	}

	return runtime.ExecuteInstructions(params.Pipe, []any{input}, scope, jpl.JPLScopedPiperFunc(func(output any, _ jpl.JPLRuntimeScope) ([]any, jpl.JPLError) {
		values, err := iter(0, output)
		if err != nil {
			return nil, err
		}
		return library.MuxAll([][]any{values}, jpl.IOMuxerFunc[any, []any](func(args ...any) ([]any, jpl.JPLError) {
			result := args[0]
			if _, ok := result.(unchanged); ok {
				result = output
			}
			return next.Pipe(result, scope)
		}))
	}))
}

// { pipe: function, selectors: [opa], assignment: opu }
func (opAssignment) Map(runtime jpl.JPLRuntime, params jpl.JPLInstructionParams) (result definition.JPLInstructionParams, err jpl.JPLError) {
	result.Pipe = call(params.Pipe)
	if result.Selectors, err = library.MuxOne([][]jpl.JPLSelector{params.Selectors}, jpl.IOMuxerFunc[jpl.JPLSelector, definition.JPLSelector](func(args ...jpl.JPLSelector) (result definition.JPLSelector, err jpl.JPLError) {
		selector := args[0]
		operator, ok := opasAssign[selector.OP]
		if !ok {
			err = library.NewFatalError("invalid OPA '" + string(selector.OP) + "' (assignment)")
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

	if params.Assignment == nil {
		params.Assignment = new(jpl.JPLAssignment)
	}
	operator, ok := opus[params.Assignment.OP]
	if !ok {
		err = library.NewFatalError("invalid OPU '" + string(params.Assignment.OP) + "'")
		return
	}
	result.Assignment = new(definition.JPLAssignment)
	result.Assignment.OP = params.Assignment.OP
	if result.Assignment.Params, err = operator.Map(runtime, params.Assignment.Params); err != nil {
		return
	}
	return
}

var opasAssign = map[definition.JPLOPA]JPLOPAHandler{
	definition.OPA_FIELD: opaAssignField{},
	definition.OPA_ITER:  opaAssignIter{},
	definition.OPA_SLICE: opaAssignSlice{},
}

type JPLOPUHandler = jpl.JPLOPSubHandler[definition.JPLAssignmentParams, jpl.JPLAssignmentParams]

var opus = map[definition.JPLOPU]JPLOPUHandler{
	definition.OPU_ADDITION:         opuAddition{},
	definition.OPU_DIVISION:         opuDivision{},
	definition.OPU_MULTIPLICATION:   opuMultiplication{},
	definition.OPU_NULL_COALESCENCE: opuNullCoalescence{},
	definition.OPU_REMAINDER:        opuRemainder{},
	definition.OPU_SET:              opuSet{},
	definition.OPU_SUBTRACTION:      opuSubtraction{},
	definition.OPU_UPDATE:           opuUpdate{},
}
