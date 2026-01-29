package library

import jpl "github.com/jplorg/jpl/go/v2/jpl"

func NewErrorEnclosure(inner jpl.JPLError) jpl.JPLErrorEnclosure {
	return errorEnclosure{inner}
}

type errorEnclosure struct{ jpl.JPLError }

func (e errorEnclosure) JPLEnclosedError() jpl.JPLError {
	return e.JPLError
}
