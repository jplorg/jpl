package program

import (
	"github.com/2manyvcos/jpl/go/definition"
	"github.com/2manyvcos/jpl/go/jpl"
)

type opAssignment struct{}

// { pipe: [op], selectors: [opa], assignment: [opu] }
func (opAssignment) OP(runtime jpl.JPLRuntime, input any, params definition.JPLInstructionParams, scope jpl.JPLRuntimeScope, next jpl.JPLScopedPiper) ([]any, jpl.JPLError) {
	panic("TODO:")
}

// { pipe: function, selectors: [opa], assignment: opu }
func (opAssignment) Map(runtime jpl.JPLRuntime, params jpl.JPLInstructionParams) (result definition.JPLInstructionParams, err jpl.JPLError) {
	panic("TODO:")
	// result.Pipe = call(params.Pipe)
	// if result.Selectors, err = library.MuxOne([][]jpl.JPLSelector{params.Selectors}, jpl.IOMuxerFunc[jpl.JPLSelector, definition.JPLSelector](func(args ...jpl.JPLSelector) (result definition.JPLSelector, err jpl.JPLError) {
	// 	selector := args[0]
	// 	operator, ok := opas[selector.OP]
	// 	if !ok {
	// 		err = library.NewFatalError("invalid OPA '" + string(selector.OP) + "'")
	// 		return
	// 	}

	// 	result.OP = selector.OP
	// 	if result.Params, err = operator.Map(runtime, selector.Params); err != nil {
	// 		return
	// 	}
	// 	return
	// })); err != nil {
	// 	return
	// }
	// return
}
