package library

// JPL muxer
type JPLMuxer interface {
	Mux(args ...any) error
}

type JPLMuxerFunc func(args ...any) error

// JPLMuxerFunc implements JPLMuxer
var _ JPLMuxer = JPLMuxerFunc(nil)

func (p JPLMuxerFunc) Mux(args ...any) error {
	return p(args...)
}

// JPL muxer
type JPLMuxerOne interface {
	MuxOne(args ...any) (any, error)
}

type JPLMuxerOneFunc func(args ...any) (any, error)

// JPLMuxerOneFunc implements JPLMuxerOne
var _ JPLMuxerOne = JPLMuxerOneFunc(nil)

func (p JPLMuxerOneFunc) MuxOne(args ...any) (any, error) {
	return p(args...)
}

// JPL muxer
type JPLMuxerAll interface {
	MuxAll(args ...any) ([]any, error)
}

type JPLMuxerAllFunc func(args ...any) ([]any, error)

// JPLMuxerAllFunc implements JPLMuxerAll
var _ JPLMuxerAll = JPLMuxerAllFunc(nil)

func (p JPLMuxerAllFunc) MuxAll(args ...any) ([]any, error) {
	return p(args...)
}
