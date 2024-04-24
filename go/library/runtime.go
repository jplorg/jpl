package library

// // Normalize the specified external value
func NormalizeValue(value any) (any, JPLError) {
	return Normalize(value)
}

// Normalize the specified array of external values
func NormalizeValues(values any, name string) ([]any, JPLError) {
	if _, ok := values.([]any); !ok {
		return nil, NewJPLFatalError("expected " + name + " to be an array")
	}
	result, err := NormalizeValue(values)
	if err != nil {
		return nil, err
	}
	return result.([]any), nil
}

// // Unwrap the specified normalized value for usage in JPL operations
// func UnwrapValue(value any) (any, JPLError) {
// 	 panic("TODO:")
// }

// // Unwrap the specified array of normalized values for usage in JPL operations
// func UnwrapValues(values []any, name string) ([]any, JPLError) {
// 	 panic("TODO:")
// }

// // Strip the specified normalized value for usage in JPL operations
// func StripValue(value any) (any, JPLError) {
// 	 panic("TODO:")
// }

// // Strip the specified array of normalized values for usage in JPL operations
// func StripValues(value []any) ([]any, JPLError) {
// 	 panic("TODO:")
// }

// // Alter the specified normalized value using the specified updater
// func AlterValue(value any, updater library.JPLModifier) (any, JPLError) {
// 	 panic("TODO:")
// }

// // Resolve the type of the specified normalized value for JPL operations
// func Type(value any) (library.JPLDataType, JPLError) {
// 	 panic("TODO:")
// }

// // Assert the type for the specified unwrapped value for JPL operations
// func AssertType(value any, assertedType library.JPLDataType) (any, JPLError) {
// 	 panic("TODO:")
// }

// // Determine whether the specified normalized value should be considered as truthy in JPL operations
// func Truthy(value any) (bool, JPLError) {
// 	 panic("TODO:")
// }

// // Compare the specified normalized values
// func Compare(a, b any) (int, JPLError) {
// 	 panic("TODO:")
// }

// // Compare the specified normalized values
// func CompareStrings(a, b any) (int, JPLError) {
// 	 panic("TODO:")
// }

// // Compare the specified normalized values
// func CompareArrays(a, b any) (int, JPLError) {
// 	 panic("TODO:")
// }

// // Compare the specified normalized values
// func CompareObjects(a, b any) (int, JPLError) {
// 	 panic("TODO:")
// }

// // Determine if the specified normalized values can be considered to be equal
// func Equals(a, b any) (bool, JPLError) {
// 	 panic("TODO:")
// }

// // Deep merge the specified normalized values
// func Merge(a, b any) (any, JPLError) {
// 	 panic("TODO:")
// }

// // Stringify the specified normalized value for usage in program outputs
// func StringifyJSON(value any, unescapeString bool) (string, JPLError) {
// 	 panic("TODO:")
// }

// Strip the specified normalized value for usage in program outputs
func StripJSON(value any) (any, JPLError) {
	return Strip(value, nil, nil)
}
