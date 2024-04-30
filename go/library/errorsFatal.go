package library

import "github.com/jplorg/jpl/go/jpl"

func NewFatalError(message string) jpl.JPLFatalError {
	return fatalError(message)
}

type fatalError string

func (e fatalError) Error() string {
	return e.JPLErrorName() + ": " + e.JPLErrorMessage()
}

func (fatalError) JPLErrorName() string {
	return "JPLFatalError"
}

func (e fatalError) JPLErrorMessage() string {
	return string(e)
}

func (fatalError) IsJPLFatalError() {}
