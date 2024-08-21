package library

import "github.com/jplorg/jpl/go/jpl"

func NewSyntaxError(message string, src string, at int) jpl.JPLSyntaxError {
	return syntaxError{message: message, src: src, at: at}
}

type syntaxError struct {
	message string
	src     string
	at      int
}

func (e syntaxError) Error() string {
	return e.JPLErrorName() + ": " + e.JPLErrorMessage()
}

func (syntaxError) JPLErrorName() string {
	return "JPLSyntaxError"
}

func (e syntaxError) JPLErrorMessage() string {
	return e.message
}

func (e syntaxError) Src() string {
	return e.src
}

func (e syntaxError) At() int {
	return e.at
}
