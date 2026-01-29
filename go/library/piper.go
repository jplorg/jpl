package library

import jpl "github.com/jplorg/jpl/go/v2/jpl"

// A JPLPiper can be used as an IOMuxer.
// For each argument, `JPLPiper.Pipe` is called and the results are combined into a single array.
type JPLPiperMuxer struct{ jpl.JPLPiper }

func NewPiperMuxer(piper jpl.JPLPiper) JPLPiperMuxer {
	return JPLPiperMuxer{JPLPiper: piper}
}

var _ jpl.IOMuxer[any, []any] = JPLPiperMuxer{}

func (m JPLPiperMuxer) Mux(args ...any) ([]any, jpl.JPLError) {
	if len(args) == 0 {
		return nil, nil
	}
	if len(args) == 1 {
		return m.Pipe(args[0])
	}
	segments := make([][]any, len(args))
	for i, arg := range args {
		var err jpl.JPLError
		segments[i], err = m.Pipe(arg)
		if err != nil {
			return nil, err
		}
	}
	return MergeSegments(segments), nil
}

// A JPLScopedPiper can be used as an IOMuxer.
// For each argument, `JPLScopedPiper.Pipe` is called and the results are combined into a single array.
type JPLScopedPiperMuxer struct {
	jpl.JPLScopedPiper
	Scope jpl.JPLRuntimeScope
}

func NewScopedPiperMuxer(piper jpl.JPLScopedPiper, scope jpl.JPLRuntimeScope) JPLScopedPiperMuxer {
	return JPLScopedPiperMuxer{JPLScopedPiper: piper, Scope: scope}
}

var _ jpl.IOMuxer[any, []any] = JPLScopedPiperMuxer{}

func (m JPLScopedPiperMuxer) Mux(args ...any) ([]any, jpl.JPLError) {
	if len(args) == 0 {
		return nil, nil
	}
	if len(args) == 1 {
		return m.Pipe(args[0], m.Scope)
	}
	segments := make([][]any, len(args))
	for i, arg := range args {
		var err jpl.JPLError
		segments[i], err = m.Pipe(arg, m.Scope)
		if err != nil {
			return nil, err
		}
	}
	return MergeSegments(segments), nil
}

// A JPLPiper can be used in place of a JPLScopedPiper
func NewPiperWithScope(piper jpl.JPLPiper) jpl.JPLScopedPiper {
	return piperWithScope{JPLPiper: piper}
}

type piperWithScope struct{ jpl.JPLPiper }

func (p piperWithScope) Pipe(input any, scope jpl.JPLRuntimeScope) ([]any, jpl.JPLError) {
	return p.JPLPiper.Pipe(input)
}

// A JPLScopedPiper can be used in place of a JPLPiper
func NewPiperWithoutScope(piper jpl.JPLScopedPiper, scope jpl.JPLRuntimeScope) jpl.JPLPiper {
	return piperWithoutScope{JPLScopedPiper: piper, Scope: scope}
}

type piperWithoutScope struct {
	jpl.JPLScopedPiper
	Scope jpl.JPLRuntimeScope
}

func (p piperWithoutScope) Pipe(input any) ([]any, jpl.JPLError) {
	return p.JPLScopedPiper.Pipe(input, p.Scope)
}
