package jpl

// Generic JPL error type
type JPLError interface {
	error
	JPLErrorName() string
}

// JPL error type for syntax errors
type JPLSyntaxError interface {
	JPLError
	IsJPLSyntaxError()
}

// JPL error type for unrecoverable errors
type JPLFatalError interface {
	JPLError
	IsJPLFatalError()
}

// JPL error type for execution errors.
//
// All error types that infer this type can be caught in a program.
type JPLExecutionError interface {
	JPLError
	JPLErrorValue() any
}

// JPL error type for generic runtime errors
type JPLRuntimeError interface {
	JPLExecutionError
	IsJPLRuntimeError()
}

// JPL runtime error type for type errors
type JPLTypeError interface {
	JPLRuntimeError
	IsJPLTypeError()
}

// JPL runtime error type for reference errors
type JPLReferenceError interface {
	JPLRuntimeError
	IsJPLReferenceError()
}

// JPL runtime error type for zero division errors
type JPLZeroDivisionError interface {
	JPLRuntimeError
	IsJPLZeroDivisionError()
}

// JPL runtime error type for type conversion errors
type JPLTypeConversionError interface {
	JPLRuntimeError
	IsJPLTypeConversionError()
}
