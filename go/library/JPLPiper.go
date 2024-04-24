package library

// JPL piper
type JPLPiper interface {
	Pipe(input any) (outputs []any, err error)
}

type JPLPiperFunc func(input any) (outputs []any, err error)

// JPLPiperFunc implements JPLPiper
var _ JPLPiper = JPLPiperFunc(nil)

func (p JPLPiperFunc) Pipe(input any) ([]any, error) {
	return p(input)
}
