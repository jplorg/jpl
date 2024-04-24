package jpl

// JPL piper
type JPLPiper interface {
	Pipe(input any) (outputs []any, err JPLError)
}

type JPLPiperFunc func(input any) (outputs []any, err JPLError)

// JPLPiperFunc implements JPLPiper
var _ JPLPiper = JPLPiperFunc(nil)

func (p JPLPiperFunc) Pipe(input any) ([]any, JPLError) {
	return p(input)
}
