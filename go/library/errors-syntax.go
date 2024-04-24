package library

// JPL error type for syntax errors
type JPLSyntaxError interface {
	JPLError
	JPLSyntaxError() bool
}

func NewJPLSyntaxError(message string) JPLSyntaxError {
	return syntaxError(message)
}

type syntaxError string

func (e syntaxError) Error() string {
	return string(e)
}

func (e syntaxError) JPLErrorName() string {
	return "JPLSyntaxError"
}

func (e syntaxError) JPLSyntaxError() bool {
	return true
}
