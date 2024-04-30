package library

import "github.com/jplorg/jpl/go/jpl"

func NewError(message string, name string) jpl.JPLError {
	return jplError{message: message, name: name}
}

type jplError struct {
	message string
	name    string
}

func (e jplError) Error() string {
	return e.JPLErrorName() + ": " + e.JPLErrorMessage()
}

func (e jplError) JPLErrorName() string {
	if e.name == "" {
		return "JPLError"
	}
	return e.name
}

func (e jplError) JPLErrorMessage() string {
	return e.message
}
