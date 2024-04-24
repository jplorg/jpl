package library

import (
	"encoding/json"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/2manyvcos/jpl/go/jpl"
)

func NewType(value any) (jpl.JPLType, jpl.JPLError) {
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
func MarshalJPLType(t jpl.JPLType) ([]byte, jpl.JPLError) {
	v, err := t.JSON()
	if err != nil {
		return nil, err
	}
	if json, err := json.Marshal(v); err != nil {
		return nil, AdaptError(err)
	} else {
		return json, nil
	}
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

// Unwrap the specified normalized value
func Unwrap(value any) (any, jpl.JPLError) {
	if t, ok := value.(jpl.JPLType); ok {
		return t.Value()
	}
	return value, nil
}

// Resolve the type of the specified normalized value
func TypeOf(value any) (jpl.JPLDataType, jpl.JPLError) {
	v, err := Unwrap(value)
	if err != nil {
		return "", err
	}
	switch v.(type) {
	case jpl.JPLFunc:
		return jpl.JPLT_FUNCTION, nil
	case map[string]any:
		return jpl.JPLT_OBJECT, nil
	case []any:
		return jpl.JPLT_ARRAY, nil
	case nil:
		return jpl.JPLT_NULL, nil
	case bool:
		return jpl.JPLT_BOOLEAN, nil
	case float64:
		return jpl.JPLT_NUMBER, nil
	case string:
		return jpl.JPLT_STRING, nil

	default:
		return "", NewFatalError(fmt.Sprintf("invalid type %T (%+v)", v, v))
	}
}

// Assert the type for the specified unwrapped value
func AssertType(value any, expectedType jpl.JPLDataType) (any, jpl.JPLError) {
	if _, ok := value.(jpl.JPLType); ok {
		if message, err := Template("unexpected type: JPLTypes (%*<100v) are not allowed here", value); err != nil {
			return nil, err
		} else {
			return nil, NewFatalError(message)
		}
	}
	t, err := TypeOf(value)
	if err != nil {
		return nil, err
	}
	if t != expectedType {
		if message, err := Template("unexpected type: %s (%*<100v) cannot be used as %s", t, value, expectedType); err != nil {
			return nil, err
		} else {
			return nil, NewFatalError(message)
		}
	}
	return value, nil
}

var functionReplacer = jpl.JPLReplacerFunc(func(k string, v any) (result any, remove bool, err jpl.JPLError) {
	u, err := Unwrap(v)
	if err != nil {
		return nil, false, err
	}
	if _, ok := u.(jpl.JPLFunc); ok {
		return "<function>", false, nil
	}
	return v, false, nil
})

// Stringify the specified normalized value
func Stringify(value any, unescapeString bool, escapeFunctions bool) (string, jpl.JPLError) {
	var rawValue any
	if escapeFunctions {
		var err jpl.JPLError
		rawValue, err = Strip(value, functionReplacer, nil)
		if err != nil {
			return "", err
		}
	} else {
		var err jpl.JPLError
		rawValue, err = Strip(value, nil, nil)
		if err != nil {
			return "", err
		}
	}
	if unescapeString {
		if s, ok := rawValue.(string); ok {
			return s, nil
		}
	}
	if json, err := json.Marshal(rawValue); err != nil {
		return "", AdaptError(err)
	} else {
		return string(json), nil
	}
}

var placeholder = regexp.MustCompile(`%(?P<Flags>[*\-<]+)?(?P<Width>[1-9][0-9]*)?(?P<Verb>.)`)

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
	var i int
	v, err := DisplayValue(tmpl)
	if err != nil {
		return "", err
	}
	result := placeholder.ReplaceAllStringFunc(v, func(match string) string {
		if err != nil {
			return match
		}
		parts := placeholder.FindStringSubmatch(match)
		if len(parts) < 4 {
			return match
		}
		flags := parts[1]
		width := parts[2]
		verb := parts[3]
		// verbs without replacement
		switch verb {
		case "%":
			return "%"
		default:
		}
		// verbs with replacement
		var value any
		if len(replacements) > i {
			value = replacements[i]
		}
		var result string
		switch verb {
		case "s":
			result, err = DisplayValue(value)
			if err != nil {
				return match
			}
		case "v":
			result, err = StrictDisplayValue(value)
			if err != nil {
				return match
			}
		default:
			err = NewFatalError("format " + match + " has unknown verb " + verb)
			return match
		}
		pad := true
		padRight := false
		trunc := false
		for _, flag := range flags {
			switch flag {
			case '*':
				pad = false
			case '-':
				padRight = true
			case '<':
				trunc = true
			default:
				err = NewFatalError("format " + match + " has unknown flag " + string(flag))
				return match
			}
		}
		var w int
		if parsed, err := strconv.ParseInt(width, 10, 64); width != "" && err == nil {
			w = int(parsed)
		}
		if w > 0 {
			rl := len([]rune(result))
			if pad && rl < w {
				padding := strings.Repeat(" ", w-rl)
				if padRight {
					result = result + padding
				} else {
					result = padding + result
				}
			} else if trunc && rl > w {
				result = string([]rune(result)[0:w-1]) + "…"
			}
		}
		return result
	})
	if err != nil {
		return "", err
	}
	return result, nil
}

// Format the specified normalized value as a string
func DisplayValue(value any) (string, jpl.JPLError) {
	return Stringify(value, true, true)
}

// Format the specified normalized value as a string, without removing escaping
func StrictDisplayValue(value any) (string, jpl.JPLError) {
	return Stringify(value, false, true)
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
	iter = func(k *string, v any) (result any, remove bool, err jpl.JPLError) {
		r := v
		if replacer != nil {
			var key string
			if k != nil {
				key = *k
			}
			var err jpl.JPLError
			if r, remove, err = replacer.Replace(key, r); err != nil {
				return nil, false, err
			} else if remove {
				return nil, true, nil
			}
		}
		return stripper.Strip(k, r, iter)
	}
	result, remove, err := iter(nil, value)
	if err != nil {
		return nil, err
	} else if remove {
		return nil, nil
	}
	return result, nil
}

// Stripper that allows JPLTypes and normalized values
var JPLTypedStripper = jpl.JPLStripperFunc(func(k *string, v any, iter jpl.IterFunc) (result any, remove bool, err jpl.JPLError) {
	if _, ok := v.(jpl.JPLType); ok {
		return v, false, nil
	}
	if _, ok := v.(jpl.JPLFunc); ok {
		return v, false, nil
	}
	return RawStripper(k, v, iter)
})

// Stripper that only allows normalized values and unwraps JPLTypes
var JPLStripper = jpl.JPLStripperFunc(func(k *string, v any, iter jpl.IterFunc) (result any, remove bool, err jpl.JPLError) {
	r := v
	if t, ok := v.(jpl.JPLType); ok {
		var err jpl.JPLError
		if r, err = t.JSON(); err != nil {
			return nil, false, err
		}
	}
	if _, ok := v.(jpl.JPLFunc); ok {
		return r, false, nil
	}
	return RawStripper(k, r, iter)
})

// Stripper that allows JSON like values and unwraps JPLTypes
var JPLJSONStripper = jpl.JPLStripperFunc(func(k *string, v any, iter jpl.IterFunc) (result any, remove bool, err jpl.JPLError) {
	r := v
	if t, ok := v.(jpl.JPLType); ok {
		var err jpl.JPLError
		if r, err = t.JSON(); err != nil {
			return nil, false, err
		}
	}
	return RawStripper(k, r, iter)
})

// Stripper that allows JSON like values and parses `json.Marshaler` interfaces
var JSONStripper = jpl.JPLStripperFunc(func(k *string, v any, iter jpl.IterFunc) (result any, remove bool, err jpl.JPLError) {
	r := v
	if m, ok := r.(json.Marshaler); ok {
		bytes, err := m.MarshalJSON()
		if err != nil {
			return nil, false, AdaptError(err)
		}
		err = json.Unmarshal(bytes, &r)
		if err != nil {
			return nil, false, AdaptError(err)
		}
	}
	return RawStripper(k, r, iter)
})

// Stripper that allows JSON like values
var RawStripper = jpl.JPLStripperFunc(func(k *string, v any, iter jpl.IterFunc) (result any, remove bool, err jpl.JPLError) {
	top := k == nil

	switch v := v.(type) {
	case jpl.JPLFunc:
		return nil, !top, nil
	case string, bool:
		return v, false, nil
	case float64:
		if math.IsNaN(v) || math.IsInf(v, 0) {
			return nil, false, nil
		}
		return v, false, nil
	case []any:
		changes := make([]*ArrayEntry[any], len(v))
		for i, entry := range v {
			key := strconv.Itoa(i)
			result, remove, err := iter(&key, entry)
			if err != nil {
				return nil, false, err
			} else if remove {
				result = nil
			}
			changes[i] = &ArrayEntry[any]{i, result}
		}
		return ApplyArray(v, changes, nil), false, nil
	case map[string]any:
		changes := make([]*ObjectEntry[any], 0, len(v))
		for i, entry := range v {
			result, remove, err := iter(&i, entry)
			if err != nil {
				return nil, false, err
			} else if remove {
				result = nil
			}
			changes = append(changes, &ObjectEntry[any]{i, result, remove})
		}
		return ApplyObject(v, changes), false, nil
	case nil:
		return nil, false, nil
	default:
		return nil, false, NewFatalError(fmt.Sprintf("unexpected %T", v))
	}
})
