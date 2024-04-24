package library

func format(value any, replacements ...any) (string, error) {
	if len(replacements) > 0 {
		return Template(value, replacements...)
	}
	return DisplayValue(value), nil
}

// JPL error type for generic runtime errors.
type JPLRuntimeError interface {
	JPLExecutionError

	JPLRuntimeError() bool
}

// `value` can by of any type.
// If at least one replacement is specified, the value is formatted as a template.
func NewJPLRuntimeError(value any, replacements ...any) (JPLRuntimeError, error) {
	message, err := format(value, replacements)
	if err != nil {
		return nil, err
	}

	if len(replacements) == 0 {
		value = message
	} else {
		value = Normalize(value)
	}

	return jplRuntimeError{
		JPLError: NewJPLError(message, "JPLRuntimeError"),
		value:    value,
	}, nil
}

type jplRuntimeError struct {
	JPLError
	value any
}

func (e jplRuntimeError) JPLErrorValue() any {
	return e.value
}

func (e jplRuntimeError) JPLRuntimeError() bool {
	return true
}

// `value` can by of any type.
// If at least one replacement is specified, the value is formatted as a template.
func NewJPLTypeError(value any, replacements ...any) (JPLTypeError, error) {
	message, err := format(value, replacements)
	if err != nil {
		return nil, err
	}
	runtimeErr, err := NewJPLRuntimeError("TypeError - " + message)
	if err != nil {
		return nil, err
	}
	return JPLTypeError(runtimeErr), nil
}

// JPL runtime error type for type errors.
type JPLTypeError JPLRuntimeError

// `value` can by of any type.
// If at least one replacement is specified, the value is formatted as a template.
func NewJPLReferenceError(value any, replacements ...any) (JPLReferenceError, error) {
	message, err := format(value, replacements)
	if err != nil {
		return nil, err
	}
	runtimeErr, err := NewJPLRuntimeError("ReferenceError - " + message)
	if err != nil {
		return nil, err
	}
	return JPLTypeError(runtimeErr), nil
}

// JPL runtime error Reference for reference errors.
type JPLReferenceError JPLRuntimeError

// `value` can by of any type.
// If at least one replacement is specified, the value is formatted as a template.
func NewJPLZeroDivisionError(value any, replacements ...any) (JPLZeroDivisionError, error) {
	message, err := format(value, replacements)
	if err != nil {
		return nil, err
	}
	runtimeErr, err := NewJPLRuntimeError("ZeroDivisionError - " + message)
	if err != nil {
		return nil, err
	}
	return JPLTypeError(runtimeErr), nil
}

// JPL runtime error type for zero division errors.
type JPLZeroDivisionError JPLRuntimeError

// `value` can by of any type.
// If at least one replacement is specified, the value is formatted as a template.
func NewJPLTypeConversionError(value any, replacements ...any) (JPLTypeConversionError, error) {
	message, err := format(value, replacements)
	if err != nil {
		return nil, err
	}
	runtimeErr, err := NewJPLRuntimeError("TypeConversionError - " + message)
	if err != nil {
		return nil, err
	}
	return JPLTypeError(runtimeErr), nil
}

// JPL runtime error type for type conversion errors.
type JPLTypeConversionError JPLRuntimeError
