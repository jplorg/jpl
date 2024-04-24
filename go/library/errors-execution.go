package library

// JPL error type for execution errors.
//
// All error types that infer this type can be caught in a program.
type JPLExecutionError interface {
	JPLError
	JPLErrorValue() any
}

func NewJPLExecutionError(message string, name string) JPLExecutionError {
	if name == "" {
		name = "JPLExecutionError"
	}

	var value any
	if message != "" {
		value = message
	}

	return jplExecutionError{
		JPLError: NewJPLError(message, name),
		value:    value,
	}
}

type jplExecutionError struct {
	JPLError
	value any
}

func (e jplExecutionError) JPLErrorValue() any {
	return e.value
}
