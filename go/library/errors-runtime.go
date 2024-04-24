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

func NewJPLRuntimeError(value any) JPLRuntimeError {
	return jplRuntimeError{
		JPLError: NewJPLError(DisplayValue(value), "JPLRuntimeError"),
		value:    Normalize(value),
	}
}

// `value` can by of any type.
// If at least one replacement is specified, the value is formatted as a template.
func NewJPLRuntimeErrorF(value any, replacements ...any) (JPLRuntimeError, error) {
	message, err := format(value, replacements)
	if err != nil {
		return nil, err
	}

	if len(replacements) > 0 {
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

func NewJPLTypeError(value any) JPLTypeError {
	return JPLTypeError{NewJPLRuntimeError("TypeError - " + DisplayValue(value))}
}

// `value` can by of any type.
// If at least one replacement is specified, the value is formatted as a template.
func NewJPLTypeErrorF(value any, replacements ...any) (JPLTypeError, error) {
	message, err := format(value, replacements)
	if err != nil {
		return JPLTypeError{}, err
	}
	return JPLTypeError{NewJPLRuntimeError("TypeError - " + message)}, nil
}

// JPL runtime error type for type errors.
type JPLTypeError struct{ JPLRuntimeError }

func NewJPLReferenceError(value any) JPLReferenceError {
	return JPLReferenceError{NewJPLRuntimeError("ReferenceError - " + DisplayValue(value))}
}

// `value` can by of any type.
// If at least one replacement is specified, the value is formatted as a template.
func NewJPLReferenceErrorF(value any, replacements ...any) (JPLReferenceError, error) {
	message, err := format(value, replacements)
	if err != nil {
		return JPLReferenceError{}, err
	}
	return JPLReferenceError{NewJPLRuntimeError("ReferenceError - " + message)}, nil
}

// JPL runtime error Reference for reference errors.
type JPLReferenceError struct{ JPLRuntimeError }

func NewJPLZeroDivisionError(value any) JPLZeroDivisionError {
	return JPLZeroDivisionError{NewJPLRuntimeError("ZeroDivisionError - " + DisplayValue(value))}
}

// `value` can by of any type.
// If at least one replacement is specified, the value is formatted as a template.
func NewJPLZeroDivisionErrorF(value any, replacements ...any) (JPLZeroDivisionError, error) {
	message, err := format(value, replacements)
	if err != nil {
		return JPLZeroDivisionError{}, err
	}
	return JPLZeroDivisionError{NewJPLRuntimeError("ZeroDivisionError - " + message)}, nil
}

// JPL runtime error type for zero division errors.
type JPLZeroDivisionError struct{ JPLRuntimeError }

func NewJPLTypeConversionError(value any) JPLTypeConversionError {
	return JPLTypeConversionError{NewJPLRuntimeError("TypeConversionError - " + DisplayValue(value))}
}

// `value` can by of any type.
// If at least one replacement is specified, the value is formatted as a template.
func NewJPLTypeConversionErrorF(value any, replacements ...any) (JPLTypeConversionError, error) {
	message, err := format(value, replacements)
	if err != nil {
		return JPLTypeConversionError{}, err
	}
	return JPLTypeConversionError{NewJPLRuntimeError("TypeConversionError - " + message)}, nil
}

// JPL runtime error type for type conversion errors.
type JPLTypeConversionError struct{ JPLRuntimeError }
