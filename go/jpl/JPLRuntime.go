package jpl

import (
	"github.com/2manyvcos/jpl/go/definition"
	"github.com/2manyvcos/jpl/go/library"
)

// JPL runtime
type JPLRuntime interface {
	// Create a new orphan scope
	CreateScope(presets *library.RuntimeScopePresets) *library.RuntimeScope

	// Execute a new dedicated program
	Execute(inputs []any) ([]any, error)

	// Execute the specified instructions
	ExecuteInstructions(instructions definition.Pipe, inputs []any, scope *library.RuntimeScope, next library.JPLPiper) ([]any, error)

	// Execute the specified OP
	OP(op definition.JPLOP, params map[string]any, inputs []any, scope *library.RuntimeScope, next library.JPLPiper) ([]any, error)

	// Normalize the specified external value
	NormalizeValue(value any) (any, error)

	// Normalize the specified array of external values
	NormalizeValues(values []any, name string) ([]any, error)

	// Unwrap the specified normalized value for usage in JPL operations
	UnwrapValue(value any) (any, error)

	// Unwrap the specified array of normalized values for usage in JPL operations
	UnwrapValues(values []any, name string) ([]any, error)

	// Strip the specified normalized value for usage in JPL operations
	StripValue(value any) (any, error)

	// Strip the specified array of normalized values for usage in JPL operations
	StripValues(value []any) ([]any, error)

	// Alter the specified normalized value using the specified updater
	AlterValue(value any, updater library.JPLModifier) (any, error)

	// Resolve the type of the specified normalized value for JPL operations
	Type(value any) (library.JPLDataType, error)

	// Assert the type for the specified unwrapped value for JPL operations
	AssertType(value any, assertedType library.JPLDataType) (any, error)

	// Determine whether the specified normalized value should be considered as truthy in JPL operations
	Truthy(value any) (bool, error)

	// Compare the specified normalized values
	Compare(a, b any) (int, error)

	// Compare the specified normalized values
	CompareStrings(a, b any) (int, error)

	// Compare the specified normalized values
	CompareArrays(a, b any) (int, error)

	// Compare the specified normalized values
	CompareObjects(a, b any) (int, error)

	// Determine if the specified normalized values can be considered to be equal
	Equals(a, b any) (bool, error)

	// Deep merge the specified normalized values
	Merge(a, b any) (any, error)

	// Stringify the specified normalized value for usage in program outputs
	StringifyJSON(value any, unescapeString bool) (string, error)

	// Strip the specified normalized value for usage in program outputs
	StripJSON(value any) (any, error)

	// Multiplex the specified array of arguments by calling cb for all possible combinations of arguments.
	//
	// `mux([[1,2], [3,4]], cb)` for example yields:
	// - `cb(1, 3)`
	// - `cb(1, 4)`
	// - `cb(2, 3)`
	// - `cb(2, 4)`
	Mux(args [][]any, cb library.JPLMuxer) error

	// Multiplex the specified array of arguments and return the results produced by the callbacks
	MuxOne(args [][]any, cb library.JPLMuxerOne) ([]any, error)

	// Multiplex the specified array of arguments and return a single array of all merged result arrays produced by the callbacks
	MuxAll(args [][]any, cb library.JPLMuxerAll) ([]any, error)
}
