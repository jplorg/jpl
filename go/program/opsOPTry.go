package program

import (
	definition "github.com/jplorg/jpl/go/v2/definition"
	jpl "github.com/jplorg/jpl/go/v2/jpl"
	library "github.com/jplorg/jpl/go/v2/library"
)

type opTry struct{}

// { try: [op], catch: [op] }
func (opTry) OP(runtime jpl.JPLRuntime, input any, params definition.JPLInstructionParams, scope jpl.JPLRuntimeScope, next jpl.JPLScopedPiper) ([]any, jpl.JPLError) {
	nextScope := scope.Next(&jpl.JPLRuntimeScopeConfig{Signal: scope.Signal().Next()})
	results, err := runtime.ExecuteInstructions(params.Try, []any{input}, nextScope, jpl.JPLScopedPiperFunc(func(output any, _ jpl.JPLRuntimeScope) ([]any, jpl.JPLError) {
		results, err := next.Pipe(output, scope)
		if err != nil {
			return nil, library.NewErrorEnclosure(err)
		}
		return results, nil
	}))
	if err != nil {
		if errorEnclosure, ok := err.(jpl.JPLErrorEnclosure); ok {
			return nil, errorEnclosure.JPLEnclosedError()
		}
		if executionErr, ok := err.(jpl.JPLExecutionError); !ok {
			return nil, err
		} else {
			nextScope.Signal().Exit()
			return runtime.ExecuteInstructions(params.Catch, []any{executionErr.JPLErrorValue()}, scope, jpl.JPLScopedPiperFunc(func(output any, _ jpl.JPLRuntimeScope) ([]any, jpl.JPLError) {
				return next.Pipe(output, scope)
			}))
		}
	}
	return results, nil
}

// { try: function, catch: function }
func (opTry) Map(runtime jpl.JPLRuntime, params jpl.JPLInstructionParams) (result definition.JPLInstructionParams, err jpl.JPLError) {
	return definition.JPLInstructionParams{
		Try:   call(params.Try),
		Catch: call(params.Catch),
	}, nil
}
