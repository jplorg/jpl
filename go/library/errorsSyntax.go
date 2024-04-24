package library

import "github.com/2manyvcos/jpl/go/jpl"

func NewJPLSyntaxError(message string) jpl.JPLSyntaxError {
	return syntaxError(message)
}

type syntaxError string

func (e syntaxError) Error() string {
	return string(e)
}

func (syntaxError) JPLErrorName() string {
	return "JPLSyntaxError"
}

func (syntaxError) IsJPLSyntaxError() {}
