package program

import (
	"github.com/2manyvcos/jpl/go/definition"
	"github.com/2manyvcos/jpl/go/jpl"
	"github.com/2manyvcos/jpl/go/library"
)

type opIf struct{}

// { ifs: [{ if: [op], then: [op] }], else: [op] }
func (opIf) OP(runtime jpl.JPLRuntime, input any, params definition.JPLInstructionParams, scope jpl.JPLRuntimeScope, next jpl.JPLScopedPiper) ([]any, jpl.JPLError) {
	var iter func(from int) ([]any, jpl.JPLError)
	iter = func(from int) ([]any, jpl.JPLError) {
		if err := scope.Signal().CheckHealth(); err != nil {
			return nil, err
		}

		if from >= len(params.Ifs) {
			return runtime.ExecuteInstructions(params.Else, []any{input}, scope, jpl.JPLScopedPiperFunc(func(output any, _ jpl.JPLRuntimeScope) ([]any, jpl.JPLError) {
				return next.Pipe(output, scope)
			}))
		}

		condition := params.Ifs[from]

		return runtime.ExecuteInstructions(condition.If, []any{input}, scope, jpl.JPLScopedPiperFunc(func(result any, _ jpl.JPLRuntimeScope) ([]any, jpl.JPLError) {
			if truthy, err := library.Truthy(result); err != nil {
				return nil, err
			} else if truthy {
				return runtime.ExecuteInstructions(condition.Then, []any{input}, scope, jpl.JPLScopedPiperFunc(func(output any, _ jpl.JPLRuntimeScope) ([]any, jpl.JPLError) {
					return next.Pipe(output, scope)
				}))
			}

			return iter(from + 1)
		}))
	}

	return iter(0)
}

// { ifs: [{ if: function, then: function }], else: function }
func (opIf) Map(runtime jpl.JPLRuntime, params jpl.JPLInstructionParams) (result definition.JPLInstructionParams, err jpl.JPLError) {
	if result.Ifs, err = library.MuxOne([][]jpl.JPLIfThen{params.Ifs}, jpl.IOMuxerFunc[jpl.JPLIfThen, definition.JPLIfThen](func(args ...jpl.JPLIfThen) (definition.JPLIfThen, jpl.JPLError) {
		entry := args[0]
		return definition.JPLIfThen{
			If:   call(entry.If),
			Then: call(entry.Then),
		}, nil
	})); err != nil {
		return
	}
	result.Else = call(params.Else)
	return
}
