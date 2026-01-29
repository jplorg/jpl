package library

import "github.com/jplorg/jpl/go/v2/jpl"

// Wrap error in a JPLExecutionError if it is not already a JPLError
func AdaptError(err error) jpl.JPLError {
	if e, ok := err.(jpl.JPLError); !ok {
		return NewExecutionError(err.Error(), "")
	} else {
		return e
	}
}
