package program

import (
	"github.com/2manyvcos/jpl/go/definition"
	"github.com/2manyvcos/jpl/go/jpl"
	"github.com/2manyvcos/jpl/go/library"
)

type opcUnequal struct{}

// { by: [op] }
func (opcUnequal) OP(runtime jpl.JPLRuntime, input any, target any, params definition.JPLComparisonParams, scope jpl.JPLRuntimeScope, next jpl.JPLPiper) ([]any, jpl.JPLError) {
	return runtime.ExecuteInstructions(params.By, []any{input}, scope, jpl.JPLScopedPiperFunc(func(by any, _ jpl.JPLRuntimeScope) ([]any, jpl.JPLError) {
		equals, err := library.Equals(target, by)
		if err != nil {
			return nil, err
		}
		return next.Pipe(!equals)
	}))
}

// { by: function }
func (opcUnequal) Map(runtime jpl.JPLRuntime, params jpl.JPLComparisonParams) (result definition.JPLComparisonParams, err jpl.JPLError) {
	result.By = call(params.By)
	return
}
