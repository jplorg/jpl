package program

import (
	"github.com/jplorg/jpl/go/definition"
	"github.com/jplorg/jpl/go/jpl"
	"github.com/jplorg/jpl/go/library"
)

type opInterpolatedString struct{}

// { interpolations: [{ before: string, pipe: [op] }], after: string }
func (opInterpolatedString) OP(runtime jpl.JPLRuntime, input any, params definition.JPLInstructionParams, scope jpl.JPLRuntimeScope, next jpl.JPLScopedPiper) ([]any, jpl.JPLError) {
	interpolations, err := library.MuxOne([][]definition.JPLInterpolation{params.Interpolations}, jpl.IOMuxerFunc[definition.JPLInterpolation, []any](func(args ...definition.JPLInterpolation) ([]any, jpl.JPLError) {
		interpolation := args[0]
		return runtime.ExecuteInstructions(interpolation.Pipe, []any{input}, scope, jpl.JPLScopedPiperFunc(func(output any, _ jpl.JPLRuntimeScope) ([]any, jpl.JPLError) {
			stringified, err := library.StringifyJSON(output, true)
			if err != nil {
				return nil, err
			}
			return []any{interpolation.Before + stringified}, nil
		}))
	}))
	if err != nil {
		return nil, err
	}

	return library.MuxAll(interpolations, jpl.IOMuxerFunc[any, []any](func(parts ...any) ([]any, jpl.JPLError) {
		var result string
		for _, part := range parts {
			result += part.(string)
		}
		return next.Pipe(result+params.After, scope)
	}))
}

// { interpolations: [{ before: string, pipe: function }], after: string }
func (opInterpolatedString) Map(runtime jpl.JPLRuntime, params jpl.JPLInstructionParams) (result definition.JPLInstructionParams, err jpl.JPLError) {
	if result.Interpolations, err = library.MuxOne([][]jpl.JPLInterpolation{params.Interpolations}, jpl.IOMuxerFunc[jpl.JPLInterpolation, definition.JPLInterpolation](func(args ...jpl.JPLInterpolation) (definition.JPLInterpolation, jpl.JPLError) {
		entry := args[0]
		return definition.JPLInterpolation{
			Before: entry.Before,
			Pipe:   call(entry.Pipe),
		}, nil
	})); err != nil {
		return
	}
	result.After = params.After
	return
}
