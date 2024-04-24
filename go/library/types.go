package library

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
	Value() (any, JPLError)

	// Resolve the JSON value for usage in program outputs
	JSON() (any, JPLError)
}

// Normalize the specified external value
func Normalize(value any) (any, JPLError) {
	panic("TODO:")
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
func Template(tmpl any, replacements ...any) (string, JPLError) {
	panic("TODO:")
}

// Format the specified normalized value as a string
func DisplayValue(value any) (string, JPLError) {
	panic("TODO:")
}

// Unwrap the specified value similar to `JSON.stringify`.
// However, unlike with `JSON.stringify`, object member functions `toJSON` are not supported as they could interfere with user defined members.
// Instead, JPLTypes are being unwrapped by default.
//
// A custom stripper can be provided to customize the behavior.
func Strip(value any, replacer JPLReplacer, stripper JPLStripper) (any, JPLError) {
	if stripper == nil {
		stripper = JPLJSONStripper
	}
	var iter IterFunc
	iter = func(k *string, v any) (any, JPLError) {
		r := v
		if replacer != nil {
			var key string
			if k != nil {
				key = *k
			}
			var err JPLError
			if r, err = replacer.Replace(key, r); err != nil {
				return nil, err
			}
		}
		return stripper.Strip(k, r, iter)
	}
	return iter(nil, value)
}

// Stripper that allows JSON like values and unwraps JPLTypes
var JPLJSONStripper = JPLStripperFunc(func(k *string, v any, iter IterFunc) (any, JPLError) {
	panic("TODO:")
})
