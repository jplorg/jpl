package library

import "github.com/2manyvcos/jpl/go/jpl"

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

		return MuxAll([][]any{results}, NewPiperMuxer(next))
	})
}
