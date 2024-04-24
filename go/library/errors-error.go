package library

// Generic JPL error type
type JPLError interface {
	error
	JPLErrorName() string
}

func NewJPLError(message string, name string) JPLError {
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
