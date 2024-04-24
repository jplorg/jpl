package library

import (
	"github.com/2manyvcos/jpl/go/definition"
	"github.com/2manyvcos/jpl/go/jpl"
)

type jplEnclosure struct {
	argNames []string
	pipe     definition.Pipe
	scope    jpl.JPLRuntimeScope
}

func (e *jplEnclosure) Call(runtime jpl.JPLRuntime, signal jpl.JPLRuntimeSignal, next jpl.JPLPiper, input any, args ...any) ([]any, error) {
	argCount := len(args)
	vars := make(map[string]any, len(e.argNames))
	for i, name := range e.argNames {
		if i < argCount {
			vars[name] = args[i]
		} else {
			vars[name] = nil
		}
	}
	return runtime.ExecuteInstructions(
		e.pipe,
		[]any{input},
		e.scope.Next(&jpl.JPLRuntimeScopeConfig{
			Signal: signal,
			Vars:   vars,
		}),
		NewPiperWithScope(next),
	)
}

// Create a scoped JPL function from the specified instructions.
//
// The function is bound to the specified scope.
func ScopedFunction(argNames []string, instructions definition.Pipe, scope jpl.JPLRuntimeScope) jpl.JPLFunc {
	return (&jplEnclosure{argNames: argNames, pipe: instructions, scope: scope}).Call
}

// Create an orphan JPL function from the specified instructions.
//
// Some optional scope presets may be specified, e.g. for allowing the function access to some specified variables.
// Other than that, the function does not have access to any external variables.
func OrphanFunction(argNames []string, instructions definition.Pipe, presets *jpl.JPLRuntimeScopeConfig) jpl.JPLFunc {
	return ScopedFunction(argNames, instructions, NewRuntimeScope(presets))
}

// Wrap the specified function to allow it to be used as a JPL function.
//
// `next` must be called for each single result, which produces an array of results itself.
// All resulting results must be returned as a single array.
//
// It is recommended to check the provided JPLRuntimeSignal in asynchronous routines to stop the routine when execution has been canceled:
// ```
// if err := signal.CheckHealth(); err != nil { return nil, err }
// ```
func NativeFunction(fn func(runtime jpl.JPLRuntime, input any, args ...any) ([]any, error)) jpl.JPLFunc {
	return jpl.JPLFunc(func(runtime jpl.JPLRuntime, signal jpl.JPLRuntimeSignal, next jpl.JPLPiper, input any, args ...any) ([]any, error) {
		results, err := fn(runtime, input, args...)
		if err != nil {
			return nil, err
		}

		return MuxAll([][]any{results}, NewPiperMuxer(next, nil))
	})
}
