package library

// Normalize the specified external value
func Normalize(value any) (any, error) {
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
func Template(tmpl any, replacements ...any) (string, error) {
	panic("TODO:")
}

// Format the specified normalized value as a string
func DisplayValue(value any) (string, error) {
	panic("TODO:")
}

// Unwrap the specified value similar to `JSON.stringify`.
// However, unlike with `JSON.stringify`, object member functions `toJSON` are not supported as they could interfere with user defined members.
// Instead, JPLTypes are being unwrapped by default.
//
// A custom stripper can be provided to customize the behavior.
func Strip(value any, replacer JPLReplacer, stripper JPLStripper) (any, error) {
	if stripper == nil {
		stripper = JPLJSONStripper
	}
	var iter IterFunc
	iter = func(k *string, v any) (any, error) {
		r := v
		if replacer != nil {
			var key string
			if k != nil {
				key = *k
			}
			var err error
			if r, err = replacer.Replace(key, r); err != nil {
				return nil, err
			}
		}
		return stripper.Strip(k, r, iter)
	}
	return iter(nil, value)
}

// Stripper that allows JSON like values and unwraps JPLTypes
var JPLJSONStripper = JPLStripperFunc(func(k *string, v any, iter IterFunc) (any, error) {
	panic("TODO:")
})
