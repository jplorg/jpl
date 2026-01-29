package library

import jpl "github.com/jplorg/jpl/go/v2/jpl"

func format(value any, replacements []any) (string, jpl.JPLError) {
	if len(replacements) > 0 {
		return Template(value, replacements...)
	}
	return DisplayValue(value)
}

// Return any of the specified errors.
func ThrowAny[ErrorType jpl.JPLError](err1 ErrorType, err2 jpl.JPLError) jpl.JPLError {
	if err2 != nil {
		return err2
	}
	return err1
}

// `value` can by of any type.
// If at least one replacement is specified, the value is formatted as a template.
func NewRuntimeError(value any, replacements ...any) (jpl.JPLRuntimeError, jpl.JPLError) {
	message, err := format(value, replacements)
	if err != nil {
		return nil, err
	}

	if len(replacements) > 0 {
		value = message
	} else {
		value, err = Normalize(value)
		if err != nil {
			return nil, err
		}
	}

	return runtimeError{
		JPLError: NewError(message, "JPLRuntimeError"),
		value:    value,
	}, nil
}

type runtimeError struct {
	jpl.JPLError
	value any
}

func (e runtimeError) JPLErrorValue() any {
	return e.value
}

func (runtimeError) IsJPLRuntimeError() {}

// `value` can by of any type.
// If at least one replacement is specified, the value is formatted as a template.
func NewTypeError(value any, replacements ...any) (jpl.JPLTypeError, jpl.JPLError) {
	message, err := format(value, replacements)
	if err != nil {
		return nil, err
	}
	runtimeErr, err := NewRuntimeError("TypeError - " + message)
	if err != nil {
		return nil, err
	}
	return typeError{runtimeErr}, nil
}

type typeError struct{ jpl.JPLRuntimeError }

func (typeError) IsJPLTypeError() {}

// `value` can by of any type.
// If at least one replacement is specified, the value is formatted as a template.
func NewReferenceError(value any, replacements ...any) (jpl.JPLReferenceError, jpl.JPLError) {
	message, err := format(value, replacements)
	if err != nil {
		return nil, err
	}
	runtimeErr, err := NewRuntimeError("ReferenceError - " + message)
	if err != nil {
		return nil, err
	}
	return referenceError{runtimeErr}, nil
}

type referenceError struct{ jpl.JPLRuntimeError }

func (referenceError) IsJPLReferenceError() {}

// `value` can by of any type.
// If at least one replacement is specified, the value is formatted as a template.
func NewZeroDivisionError(value any, replacements ...any) (jpl.JPLZeroDivisionError, jpl.JPLError) {
	message, err := format(value, replacements)
	if err != nil {
		return nil, err
	}
	runtimeErr, err := NewRuntimeError("ZeroDivisionError - " + message)
	if err != nil {
		return nil, err
	}
	return zeroDivisionError{runtimeErr}, nil
}

type zeroDivisionError struct{ jpl.JPLRuntimeError }

func (zeroDivisionError) IsJPLZeroDivisionError() {}

// `value` can by of any type.
// If at least one replacement is specified, the value is formatted as a template.
func NewTypeConversionError(value any, replacements ...any) (jpl.JPLTypeConversionError, jpl.JPLError) {
	message, err := format(value, replacements)
	if err != nil {
		return nil, err
	}
	runtimeErr, err := NewRuntimeError("TypeConversionError - " + message)
	if err != nil {
		return nil, err
	}
	return typeConversionError{runtimeErr}, nil
}

type typeConversionError struct{ jpl.JPLRuntimeError }

func (typeConversionError) IsJPLTypeConversionError() {}
