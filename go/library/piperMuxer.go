package library

import "github.com/2manyvcos/jpl/go/jpl"

// A JPLPiper can be used as an IOMuxer.
// For each argument, `JPLPiper.Pipe` is called and the results are combined into a single array.
type JPLPiperMuxer func(input any) (outputs []any, err jpl.JPLError)

func NewPiperMuxer(piper jpl.JPLPiper) JPLPiperMuxer {
	return JPLPiperMuxer(piper.Pipe)
}

var _ jpl.IOMuxer[any, []any] = JPLPiperMuxer(nil)

func (m JPLPiperMuxer) Mux(args ...any) ([]any, jpl.JPLError) {
	if len(args) == 0 {
		return nil, nil
	}
	if len(args) == 1 {
		return m(args[0])
	}
	segments := make([][]any, len(args))
	for i, arg := range args {
		var err jpl.JPLError
		segments[i], err = m(arg)
		if err != nil {
			return nil, err
		}
	}
	return MergeSegments(segments), nil
}
