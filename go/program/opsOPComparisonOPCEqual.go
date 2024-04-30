package program

import (
	"github.com/jplorg/jpl/go/definition"
	"github.com/jplorg/jpl/go/jpl"
	"github.com/jplorg/jpl/go/library"
)

type opcEqual struct{}

// { by: [op] }
func (opcEqual) OP(runtime jpl.JPLRuntime, input any, target any, params definition.JPLComparisonParams, scope jpl.JPLRuntimeScope, next jpl.JPLPiper) ([]any, jpl.JPLError) {
	return runtime.ExecuteInstructions(params.By, []any{input}, scope, jpl.JPLScopedPiperFunc(func(by any, _ jpl.JPLRuntimeScope) ([]any, jpl.JPLError) {
		equals, err := library.Equals(target, by)
		if err != nil {
			return nil, err
		}
		return next.Pipe(equals)
	}))
}

// { by: function }
func (opcEqual) Map(runtime jpl.JPLRuntime, params jpl.JPLComparisonParams) (result definition.JPLComparisonParams, err jpl.JPLError) {
	result.By = call(params.By)
	return
}
