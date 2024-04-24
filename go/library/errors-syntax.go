package library

// JPL error type for syntax errors
type JPLSyntaxError interface {
	JPLError
	JPLSyntaxError() bool
}

func NewJPLSyntaxError(message string) JPLSyntaxError {
	return jplSyntaxError(message)
}

type jplSyntaxError string

func (e jplSyntaxError) Error() string {
	return string(e)
}

func (e jplSyntaxError) Name() string {
	return "JPLSyntaxError"
}

func (e jplSyntaxError) JPLError() bool {
	return true
}

func (e jplSyntaxError) JPLSyntaxError() bool {
	return true
}
