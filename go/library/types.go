package library

// Normalize the specified external value
func Normalize(value any) any {
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
func DisplayValue(value any) string {
	panic("TODO:")
}
