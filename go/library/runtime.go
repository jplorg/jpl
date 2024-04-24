package library

import "github.com/2manyvcos/jpl/go/jpl"

// // Normalize the specified external value
func NormalizeValue(value any) (any, jpl.JPLError) {
	return Normalize(value)
}

// Normalize the specified array of external values
func NormalizeValues(values any, name string) ([]any, jpl.JPLError) {
	if values == nil {
		values = []any(nil)
	}
	if name == "" {
		name = "values"
	}
	if _, ok := values.([]any); !ok {
		return nil, NewFatalError("expected " + name + " to be an array")
	}
	result, err := NormalizeValue(values)
	if err != nil {
		return nil, err
	}
	return result.([]any), nil
}

// Unwrap the specified normalized value for usage in JPL operations
func UnwrapValue(value any) (any, jpl.JPLError) {
	return Unwrap(value)
}

// Unwrap the specified array of normalized values for usage in JPL operations
func UnwrapValues(values any, name string) ([]any, jpl.JPLError) {
	if values == nil {
		values = []any(nil)
	}
	if name == "" {
		name = "values"
	}
	v, ok := values.([]any)
	if !ok {
		return nil, NewFatalError("expected " + name + " to be an array")
	}
	return MuxOne([][]any{v}, jpl.IOMuxerFunc[any, any](func(args ...any) (any, jpl.JPLError) {
		return UnwrapValue(args[0])
	}))
}

var unwrapReplacer = jpl.JPLReplacerFunc(func(k string, v any) (result any, remove bool, err jpl.JPLError) {
	u, err := UnwrapValue(v)
	return u, false, err
})

// Strip the specified normalized value for usage in JPL operations
func StripValue(value any) (any, jpl.JPLError) {
	return Strip(value, unwrapReplacer, nil)
}

// Strip the specified array of normalized values for usage in JPL operations
func StripValues(values any, name string) ([]any, jpl.JPLError) {
	if values == nil {
		values = []any(nil)
	}
	if name == "" {
		name = "values"
	}
	v, ok := values.([]any)
	if !ok {
		return nil, NewFatalError("expected " + name + " to be an array")
	}
	return MuxOne([][]any{v}, jpl.IOMuxerFunc[any, any](func(args ...any) (any, jpl.JPLError) {
		return StripValue(args[0])
	}))
}

// Alter the specified normalized value using the specified updater
func AlterValue(value any, updater jpl.JPLModifier) (any, jpl.JPLError) {
	if t, ok := value.(jpl.JPLType); ok {
		return t.Alter(updater)
	}
	result, err := updater.Modify(value)
	if err != nil {
		return nil, err
	}
	return NormalizeValue(result)
}

// Resolve the type of the specified normalized value for JPL operations
func Type(value any) (jpl.JPLDataType, jpl.JPLError) {
	return TypeOf(value)
}

// Determine whether the specified normalized value should be considered as truthy in JPL operations
func Truthy(value any) (bool, jpl.JPLError) {
	raw, err := UnwrapValue(value)
	if err != nil {
		return false, err
	}
	return !IsSame(raw, nil) && !IsSame(raw, false), nil
}

// // Compare the specified normalized values
// func Compare(a, b any) (int, jpl.JPLError) {
// 	 panic("TODO:")
// }

// // Compare the specified normalized values
// func CompareStrings(a, b any) (int, jpl.JPLError) {
// 	 panic("TODO:")
// }

// // Compare the specified normalized values
// func CompareArrays(a, b any) (int, jpl.JPLError) {
// 	 panic("TODO:")
// }

// // Compare the specified normalized values
// func CompareObjects(a, b any) (int, jpl.JPLError) {
// 	 panic("TODO:")
// }

// // Determine if the specified normalized values can be considered to be equal
// func Equals(a, b any) (bool, jpl.JPLError) {
// 	 panic("TODO:")
// }

// // Deep merge the specified normalized values
// func Merge(a, b any) (any, jpl.JPLError) {
// 	 panic("TODO:")
// }

// Stringify the specified normalized value for usage in program outputs
func StringifyJSON(value any, unescapeString bool) (string, jpl.JPLError) {
	return Stringify(value, unescapeString, false)
}

// Strip the specified normalized value for usage in program outputs
func StripJSON(value any) (any, jpl.JPLError) {
	return Strip(value, nil, nil)
}
