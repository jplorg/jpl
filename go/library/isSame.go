package library

import (
	"reflect"

	"github.com/2manyvcos/jpl/go/jpl"
)

// Check if both values are the same.
//
// []any and map[string]any are compared by memory address,
// primitives by value.
func IsSame(a, b any) bool {
	aType, aOk := a.(jpl.JPLType)
	bType, bOk := a.(jpl.JPLType)
	if aOk || bOk {
		return aOk && bOk && reflect.TypeOf(a) == reflect.TypeOf(b) && aType.IsSame(bType)
	}

	_, aOk = a.(jpl.JPLFunc)
	_, bOk = b.(jpl.JPLFunc)
	if aOk || bOk {
		return aOk && bOk && reflect.ValueOf(a).Pointer() == reflect.ValueOf(b).Pointer()
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
