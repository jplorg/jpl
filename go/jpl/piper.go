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

// JPL scoped piper
type JPLScopedPiper interface {
	Pipe(input any, scope JPLRuntimeScope) (outputs []any, err JPLError)
}

type JPLScopedPiperFunc func(input any, scope JPLRuntimeScope) (outputs []any, err JPLError)

// JPLPiperFunc implements JPLPiper
var _ JPLScopedPiper = JPLScopedPiperFunc(nil)

func (p JPLScopedPiperFunc) Pipe(input any, scope JPLRuntimeScope) ([]any, JPLError) {
	return p(input, scope)
}
