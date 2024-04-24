package library

import "reflect"

// Check if both values are the same.
//
// Slices and maps are compared by memory address,
// primitives by value.
func Is(a, b any) bool {
	// TODO: what about JPLTypes and functions?

	_, aOk := a.(map[string]any)
	_, bOk := b.(map[string]any)
	if aOk || bOk {
		return aOk && bOk && reflect.ValueOf(a).Pointer() == reflect.ValueOf(b).Pointer()
	}

	aV, aOk := a.([]any)
	bV, bOk := b.([]any)
	if aOk || bOk {
		return aOk && bOk && len(aV) == len(bV) && reflect.ValueOf(a).Pointer() == reflect.ValueOf(b).Pointer()
	}

	aR := reflect.ValueOf(a)
	bR := reflect.ValueOf(b)
	if aV, bV := aR.IsValid(), bR.IsValid(); !aV || !bV {
		return !aV && !bV
	}
	if !aR.Comparable() || !bR.Comparable() {
		return false
	}

	return a == b
}
