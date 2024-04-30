package library

import "github.com/jplorg/jpl/go/jpl"

func NewSyntaxError(message string) jpl.JPLSyntaxError {
	return syntaxError(message)
}

type syntaxError string

func (e syntaxError) Error() string {
	return e.JPLErrorName() + ": " + e.JPLErrorMessage()
}

func (syntaxError) JPLErrorName() string {
	return "JPLSyntaxError"
}

func (e syntaxError) JPLErrorMessage() string {
	return string(e)
}

func (syntaxError) IsJPLSyntaxError() {}
