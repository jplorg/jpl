package jpl

import "encoding/json"

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

	// Alter the internal value using the specified updater.
	// The result may or may not be wrapped into a new JPLType if desired.
	// The JPLType may return itself if the value did not change.
	//
	// It is recommended to use AlterJPLType here as a starting point.
	Alter(updater JPLModifier) (any, JPLError)

	// Return whether the specified JPLType refers to the same value as the receiver, meaning that both instances are interchangeable.
	// `other` is always expected to be of the same type than the receiver.
	//
	// For a JPLType that is defined as a pointer, both pointers can be simply compared, e.g.
	// `func (t *someType) IsSame(other JPLType) bool { return t == other }`
	IsSame(other JPLType) bool

	// It is recommended to use MarshalJSONJPLType here, which marshals the JPLType based on its json accessor.
	json.Marshaler
}

// Shordhand type for JPL functions
type JPLFunc = func(runtime JPLRuntime, signal JPLRuntimeSignal, next JPLPiper, input any, args ...any) ([]any, error)
