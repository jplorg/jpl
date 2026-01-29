package program

import (
	"github.com/jplorg/jpl/go/v2/definition"
	"github.com/jplorg/jpl/go/v2/jpl"
	"github.com/jplorg/jpl/go/v2/library"
)

type opcLess struct{}

// { by: [op] }
func (opcLess) OP(runtime jpl.JPLRuntime, input any, target any, params definition.JPLComparisonParams, scope jpl.JPLRuntimeScope, next jpl.JPLPiper) ([]any, jpl.JPLError) {
	return runtime.ExecuteInstructions(params.By, []any{input}, scope, jpl.JPLScopedPiperFunc(func(by any, _ jpl.JPLRuntimeScope) ([]any, jpl.JPLError) {
		c, err := library.Compare(target, by)
		if err != nil {
			return nil, err
		}
		return next.Pipe(c < 0)
	}))
}

// { by: function }
func (opcLess) Map(runtime jpl.JPLRuntime, params jpl.JPLComparisonParams) (result definition.JPLComparisonParams, err jpl.JPLError) {
	result.By = call(params.By)
	return
}
