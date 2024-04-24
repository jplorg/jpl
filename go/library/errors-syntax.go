package library

func NewJPLSyntaxError(message string) JPLSyntaxError {
	return JPLSyntaxError(message)
}

// JPL error type for syntax errors
type JPLSyntaxError string

// JPLSyntaxError implements JPLError
var _ JPLError = JPLSyntaxError("")

func (e JPLSyntaxError) Error() string {
	return string(e)
}

func (e JPLSyntaxError) JPLErrorName() string {
	return "JPLSyntaxError"
}
