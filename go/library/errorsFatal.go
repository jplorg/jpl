package library

import "github.com/2manyvcos/jpl/go/jpl"

func NewJPLFatalError(message string) jpl.JPLFatalError {
	return fatalError(message)
}

type fatalError string

func (e fatalError) Error() string {
	return string(e)
}

func (fatalError) JPLErrorName() string {
	return "JPLFatalError"
}

func (fatalError) IsJPLFatalError() {}
