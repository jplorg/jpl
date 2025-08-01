package library

import (
	"reflect"

	"github.com/jplorg/jpl/go/jpl"
)

// Check if both values are the same.
//
// JPLTypes are compared by type and `JPLType.IsSame`,
// JPLFuncs, []any and map[string]any are compared by memory address,
// and other comparable values using `==`.
func IsSame(a, b any) bool {
	aType, aOk := a.(jpl.JPLType)
	bType, bOk := a.(jpl.JPLType)
	if aOk || bOk {
		return aOk && bOk && reflect.TypeOf(a) == reflect.TypeOf(b) && aType.IsSame(bType)
	}

	aFunc, aOk := a.(jpl.JPLFunc)
	bFunc, bOk := b.(jpl.JPLFunc)
	if aOk || bOk {
		return aOk && bOk && reflect.TypeOf(a) == reflect.TypeOf(b) && aFunc.IsSame(bFunc)
	}

	_, aOk = a.(map[string]any)
	_, bOk = b.(map[string]any)
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
