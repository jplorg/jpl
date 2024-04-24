package library

import (
	"encoding/json"

	"github.com/2manyvcos/jpl/go/jpl"
)

func NewJPLType(value any) (jpl.JPLType, jpl.JPLError) {
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

func (t *jplType) IsSame(other jpl.JPLType) bool {
	return t == other
}

func (t *jplType) MarshalJSON() ([]byte, error) {
	return MarshalJPLType(t)
}

// Alter the specified JPLType
func AlterJPLType(t jpl.JPLType, updater jpl.JPLModifier) (any, jpl.JPLError) {
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
func MarshalJPLType(t jpl.JPLType) ([]byte, error) {
	v, err := t.JSON()
	if err != nil {
		return nil, err
	}
	return json.Marshal(v)
}

// Normalize the specified external value to be used in a JPLType
func normalizeInternalValue(value any) (any, jpl.JPLError) {
	if t, ok := value.(jpl.JPLType); ok {
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
