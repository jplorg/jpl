package library

import "github.com/jplorg/jpl/go/jpl"

func NewErrorEnclosure(inner jpl.JPLError) jpl.JPLErrorEnclosure {
	return errorEnclosure{inner}
}

type errorEnclosure struct{ jpl.JPLError }

func (e errorEnclosure) JPLEnclosedError() jpl.JPLError {
	return e.JPLError
}
