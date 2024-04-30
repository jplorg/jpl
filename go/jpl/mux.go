package jpl

// JPL muxer
type IMuxer[Input any] interface {
	Mux(args ...Input) JPLError
}

type IMuxerFunc[Input any] func(args ...Input) JPLError

// IMuxerFunc implements IMuxer
var _ IMuxer[any] = IMuxerFunc[any](nil)

func (p IMuxerFunc[Input]) Mux(args ...Input) JPLError {
	return p(args...)
}

// JPL muxer
type IOMuxer[Input any, Output any] interface {
	Mux(args ...Input) (Output, JPLError)
}

type IOMuxerFunc[Input any, Output any] func(args ...Input) (Output, JPLError)

// IOMuxerFunc implements IOMuxer
var _ IOMuxer[any, any] = IOMuxerFunc[any, any](nil)

func (p IOMuxerFunc[Input, Output]) Mux(args ...Input) (Output, JPLError) {
	return p(args...)
}
