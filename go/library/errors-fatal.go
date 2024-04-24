package library

func NewJPLFatalError(message string) JPLFatalError {
	return JPLFatalError(message)
}

// JPL error type for unrecoverable errors
type JPLFatalError string

// JPLFatalError implements JPLError
var _ JPLError = JPLFatalError("")

func (e JPLFatalError) Error() string {
	return string(e)
}

func (e JPLFatalError) JPLErrorName() string {
	return "JPLFatalError"
}
