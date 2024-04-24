package library

import "github.com/2manyvcos/jpl/go/jpl"

func format(value any, replacements ...any) (string, jpl.JPLError) {
	if len(replacements) > 0 {
		return Template(value, replacements...)
	}
	return DisplayValue(value)
}

// `value` can by of any type.
// If at least one replacement is specified, the value is formatted as a template.
func NewJPLRuntimeError(value any, replacements ...any) (jpl.JPLRuntimeError, jpl.JPLError) {
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
		JPLError: NewJPLError(message, "JPLRuntimeError"),
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
func NewJPLTypeError(value any, replacements ...any) (jpl.JPLTypeError, jpl.JPLError) {
	message, err := format(value, replacements)
	if err != nil {
		return nil, err
	}
	runtimeErr, err := NewJPLRuntimeError("TypeError - " + message)
	if err != nil {
		return nil, err
	}
	return typeError{runtimeErr}, nil
}

type typeError struct{ jpl.JPLRuntimeError }

func (typeError) IsJPLTypeError() {}

// `value` can by of any type.
// If at least one replacement is specified, the value is formatted as a template.
func NewJPLReferenceError(value any, replacements ...any) (jpl.JPLReferenceError, jpl.JPLError) {
	message, err := format(value, replacements)
	if err != nil {
		return nil, err
	}
	runtimeErr, err := NewJPLRuntimeError("ReferenceError - " + message)
	if err != nil {
		return nil, err
	}
	return referenceError{runtimeErr}, nil
}

type referenceError struct{ jpl.JPLRuntimeError }

func (referenceError) IsJPLReferenceError() {}

// `value` can by of any type.
// If at least one replacement is specified, the value is formatted as a template.
func NewJPLZeroDivisionError(value any, replacements ...any) (jpl.JPLZeroDivisionError, jpl.JPLError) {
	message, err := format(value, replacements)
	if err != nil {
		return nil, err
	}
	runtimeErr, err := NewJPLRuntimeError("ZeroDivisionError - " + message)
	if err != nil {
		return nil, err
	}
	return zeroDivisionError{runtimeErr}, nil
}

type zeroDivisionError struct{ jpl.JPLRuntimeError }

func (zeroDivisionError) IsJPLZeroDivisionError() {}

// `value` can by of any type.
// If at least one replacement is specified, the value is formatted as a template.
func NewJPLTypeConversionError(value any, replacements ...any) (jpl.JPLTypeConversionError, jpl.JPLError) {
	message, err := format(value, replacements)
	if err != nil {
		return nil, err
	}
	runtimeErr, err := NewJPLRuntimeError("TypeConversionError - " + message)
	if err != nil {
		return nil, err
	}
	return typeConversionError{runtimeErr}, nil
}

type typeConversionError struct{ jpl.JPLRuntimeError }

func (typeConversionError) IsJPLTypeConversionError() {}
