package library

import "github.com/2manyvcos/jpl/go/jpl"

func NewJPLError(message string, name string) jpl.JPLError {
	if name == "" {
		name = "JPLError"
	}

	return jplError{message: message, name: name}
}

type jplError struct {
	message string
	name    string
}

func (e jplError) Error() string {
	return e.message
}

func (e jplError) JPLErrorName() string {
	return e.name
}
