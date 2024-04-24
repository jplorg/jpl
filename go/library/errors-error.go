package library

type JPLError interface {
	error
	Name() string
	JPLError() bool
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

func (e jplError) Name() string {
	return e.name
}

func (e jplError) JPLError() bool {
	return true
}
