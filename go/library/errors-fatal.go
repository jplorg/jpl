package library

// JPL error type for unrecoverable errors
type JPLFatalError interface {
	JPLError
	JPLFatalError() bool
}

func NewJPLFatalError(message string) JPLFatalError {
	return fatalError(message)
}

type fatalError string

func (e fatalError) Error() string {
	return string(e)
}

func (e fatalError) JPLErrorName() string {
	return "JPLFatalError"
}

func (e fatalError) JPLFatalError() bool {
	return true
}
