package library

import "github.com/2manyvcos/jpl/go/jpl"

func NewExecutionError(message string, name string) jpl.JPLExecutionError {
	if name == "" {
		name = "JPLExecutionError"
	}

	var value any
	if message != "" {
		value = message
	}

	return executionError{
		JPLError: NewError(message, name),
		value:    value,
	}
}

type executionError struct {
	jpl.JPLError
	value any
}

func (e executionError) JPLErrorValue() any {
	return e.value
}
