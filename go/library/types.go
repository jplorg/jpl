package library

import (
	"encoding/json"

	"github.com/2manyvcos/jpl/go/jpl"
)

// Generic type for handling special formatting on values.
//
// The value accessor should be used in JPL expressions, where type safety may be considered important.
//
// The json accessor should be used in program outputs, where the format of the value should be applied and the original type may be ignored.
//
// The alter function is called when the value is being changed by an operation.
// The JPLType is responsible to decide whether the resulting value should be wrapped in a new JPLType or not.
// It may also return itself if the value did not change.
// Note that the resulting value may have another type than the original value.
// Commonly, if multiple values wrapped in JPLTypes are involved in an operation, the leftmost operand should be altered whereas all other operands should be unwrapped using their value accessors, effectively loosing their formatting.
//
// However, both accessors should represent the same essential value.
// For example, a JPLType that applies rounding to numbers with a fixed number of decimal digits, may return the rounded numeric value for its value accessor (e.g. `1`), whereas it may return a formatted string for its JSON accessor (e.g. `"1.00"`).
// This allows this JPLType to be processed in JPL operations like generic numbers but resolves to formatted strings in the program output.
type JPLType interface {
	// Resolve the internal value for usage in JPL operations
	Value() (any, jpl.JPLError)

	// Resolve the JSON value for usage in program outputs
	JSON() (any, jpl.JPLError)

	// Alter the internal value using the specified updater.
	// The result may or may not be wrapped into a new JPLType if desired.
	// The JPLType may return itself if the value did not change.
	//
	// It is recommended to use AlterJPLType here as a starting point.
	Alter(updater jpl.JPLModifier) (any, jpl.JPLError)

	// Return whether the specified JPLType refers to the same value as the receiver, meaning that both instances are interchangeable.
	// `other` is always expected to be of the same type than the receiver.
	//
	// For a JPLType that is defined as a pointer, both pointers can be simply compared, e.g.
	// `func (t *someType) IsSame(other JPLType) bool { return t == other }`
	IsSame(other JPLType) bool

	// It is recommended to use MarshalJSONJPLType here, which marshals the JPLType based on its json accessor.
	json.Marshaler
}

func NewJPLType(value any) (JPLType, jpl.JPLError) {
	v, err := normalizeInternalValue(value)
	if err != nil {
		return nil, err
	}
	return &jplType{value: v}, nil
}

type jplType struct{ value any }

func (t *jplType) Value() (any, jpl.JPLError) {
	return t.value, nil
}

func (t *jplType) JSON() (any, jpl.JPLError) {
	return t.value, nil
}

func (t *jplType) Alter(updater jpl.JPLModifier) (any, jpl.JPLError) {
	return AlterJPLType(t, updater)
}

func (t *jplType) IsSame(other JPLType) bool {
	return t == other
}

func (t *jplType) MarshalJSON() ([]byte, error) {
	return MarshalJPLType(t)
}

// Alter the specified JPLType
func AlterJPLType(t JPLType, updater jpl.JPLModifier) (any, jpl.JPLError) {
	v, err := t.Value()
	if err != nil {
		return nil, err
	}
	v, err = updater.Modify(v)
	if err != nil {
		return nil, err
	}
	return Normalize(v)
}

// Marshal the specified JPLType as JSON
func MarshalJPLType(t JPLType) ([]byte, error) {
	v, err := t.JSON()
	if err != nil {
		return nil, err
	}
	return json.Marshal(v)
}

// Normalize the specified external value to be used in a JPLType
func normalizeInternalValue(value any) (any, jpl.JPLError) {
	if t, ok := value.(JPLType); ok {
		return t.Value()
	}
	return Normalize(value)
}

// Normalize the specified external value
func Normalize(value any) (any, jpl.JPLError) {
	return Strip(value, nil, JPLTypedStripper)
}

// Format the specified template string.
// The general form of a format is a percent sign, followed by optional flags, an optional width and a verb.
//
// Examples:
// - `%10s`: Format the next replacement as a string and pads the result at the left with spaces to be at least 10 unicode codepoints long.
// - `%*<10v`: Format the next replacement as a JSON value and truncates it, if it is longer then 10 unicode codepoints.
//
// The width specifies the desired field width and defaults to whatever is necessary to display the full replacement.
// If the width is specified without any other corresponding flags, it is used for padding the field if necessary.
//
// Valid flags:
//
// - `*`: Do not pad the value even if it is shorter than the specified width
// - `-`: Pad the value at the right rather than the left
// - `<`: Truncate the value at the right if it is too long for the specified width
//
// Valid verbs:
//
// - `%`: Returns a literal `%`
// - `s`: Format the next replacement as a string (like JSON, but does not escape strings)
// - `v`: Format the next replacement as a JSON value
func Template(tmpl any, replacements ...any) (string, jpl.JPLError) {
	panic("TODO:")
}

// Format the specified normalized value as a string
func DisplayValue(value any) (string, jpl.JPLError) {
	panic("TODO:")
}

// Unwrap the specified value similar to `JSON.stringify`.
// However, unlike with `JSON.stringify`, object member functions `toJSON` are not supported as they could interfere with user defined members.
// Instead, JPLTypes are being unwrapped by default.
//
// A custom stripper can be provided to customize the behavior.
func Strip(value any, replacer jpl.JPLReplacer, stripper jpl.JPLStripper) (any, jpl.JPLError) {
	if stripper == nil {
		stripper = JPLJSONStripper
	}
	var iter jpl.IterFunc
	iter = func(k *string, v any) (any, jpl.JPLError) {
		r := v
		if replacer != nil {
			var key string
			if k != nil {
				key = *k
			}
			var err jpl.JPLError
			if r, err = replacer.Replace(key, r); err != nil {
				return nil, err
			}
		}
		return stripper.Strip(k, r, iter)
	}
	return iter(nil, value)
}

// Stripper that allows JPLTypes and normalized values
var JPLTypedStripper = jpl.JPLStripperFunc(func(k *string, v any, iter jpl.IterFunc) (any, jpl.JPLError) {
	panic("TODO:")
})

// Stripper that allows JSON like values and unwraps JPLTypes
var JPLJSONStripper = jpl.JPLStripperFunc(func(k *string, v any, iter jpl.IterFunc) (any, jpl.JPLError) {
	panic("TODO:")
})
