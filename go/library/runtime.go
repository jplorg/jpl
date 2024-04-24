package library

import (
	"slices"

	"github.com/2manyvcos/jpl/go/jpl"
)

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

// Compare the specified normalized values
func Compare(a, b any) (int, jpl.JPLError) {
	ta, err := Type(a)
	if err != nil {
		return 0, err
	}
	tb, err := Type(b)
	if err != nil {
		return 0, err
	}

	if ta != tb {
		return slices.Index(jpl.TypeOrder, ta) - slices.Index(jpl.TypeOrder, tb), nil
	}

	ua, err := UnwrapValue(a)
	if err != nil {
		return 0, err
	}
	ub, err := UnwrapValue(b)
	if err != nil {
		return 0, err
	}

	switch ta {
	case jpl.JPLT_NULL, jpl.JPLT_FUNCTION:
		return 0, nil

	case jpl.JPLT_BOOLEAN:
		var numA, numB int
		if ua == true {
			numA = 1
		}
		if ub == true {
			numB = 1
		}
		return numA - numB, nil

	case jpl.JPLT_NUMBER:
		numA := ua.(float64)
		numB := ub.(float64)
		if numA < numB {
			return -1, nil
		} else if numA > numB {
			return 1, nil
		} else {
			return 0, nil
		}

	case jpl.JPLT_STRING:
		return compareStrings(ua.(string), ub.(string)), nil

	case jpl.JPLT_ARRAY:
		return compareArrays(ua.([]any), ub.([]any))

	case jpl.JPLT_OBJECT:
		return compareObjects(ua.(map[string]any), ub.(map[string]any))

	default:
		return 0, NewFatalError("unexpected type " + string(ta))
	}
}

// Compare the specified normalized strings based on their unicode code points
func CompareStrings(a, b any) (int, jpl.JPLError) {
	ta, err := Type(a)
	if err != nil {
		return 0, err
	}
	if ta != jpl.JPLT_STRING {
		return 0, NewFatalError("unexpected type " + string(ta))
	}
	tb, err := Type(b)
	if err != nil {
		return 0, err
	}
	if tb != jpl.JPLT_STRING {
		return 0, NewFatalError("unexpected type " + string(tb))
	}
	ua, err := UnwrapValue(a)
	if err != nil {
		return 0, err
	}
	ub, err := UnwrapValue(b)
	if err != nil {
		return 0, err
	}
	return compareStrings(ua.(string), ub.(string)), nil
}

// Compare the specified normalized strings based on their unicode code points
func compareStrings(a, b string) int {
	ua := []rune(a)
	ub := []rune(b)
	min := min(len(ua), len(ub))
	for i, cp1 := range ua {
		if i >= min {
			break
		}
		cp2 := ub[i]
		if cp1 < cp2 {
			return -1
		} else if cp1 > cp2 {
			return 1
		}
	}
	return len(ua) - len(ub)
}

// Compare the specified normalized arrays based on their lexical order
func CompareArrays(a, b any) (int, jpl.JPLError) {
	ta, err := Type(a)
	if err != nil {
		return 0, err
	}
	if ta != jpl.JPLT_ARRAY {
		return 0, NewFatalError("unexpected type " + string(ta))
	}
	tb, err := Type(b)
	if err != nil {
		return 0, err
	}
	if tb != jpl.JPLT_ARRAY {
		return 0, NewFatalError("unexpected type " + string(tb))
	}
	ua, err := UnwrapValue(a)
	if err != nil {
		return 0, err
	}
	ub, err := UnwrapValue(b)
	if err != nil {
		return 0, err
	}
	return compareArrays(ua.([]any), ub.([]any))
}

// Compare the specified normalized arrays based on their lexical order
func compareArrays[Value any](a, b []Value) (int, jpl.JPLError) {
	min := min(len(a), len(b))
	for i := 0; i < min; i += 1 {
		c, err := Compare(a[i], b[i])
		if err != nil {
			return 0, err
		}
		if c != 0 {
			return c, nil
		}
	}
	return len(a) - len(b), nil
}

// Compare the specified normalized objects
func CompareObjects(a, b any) (int, jpl.JPLError) {
	ta, err := Type(a)
	if err != nil {
		return 0, err
	}
	if ta != jpl.JPLT_OBJECT {
		return 0, NewFatalError("unexpected type " + string(ta))
	}
	tb, err := Type(b)
	if err != nil {
		return 0, err
	}
	if tb != jpl.JPLT_OBJECT {
		return 0, NewFatalError("unexpected type " + string(tb))
	}
	ua, err := UnwrapValue(a)
	if err != nil {
		return 0, err
	}
	ub, err := UnwrapValue(b)
	if err != nil {
		return 0, err
	}
	return compareObjects(ua.(map[string]any), ub.(map[string]any))
}

// Compare the specified normalized objects
func compareObjects(a, b map[string]any) (int, jpl.JPLError) {
	aKeys := mapKeys(a)
	slices.SortFunc(aKeys, compareStrings)
	bKeys := mapKeys(b)
	slices.SortFunc(bKeys, compareStrings)
	order, err := compareArrays(aKeys, bKeys)
	if err != nil {
		return 0, err
	}
	if order != 0 {
		return order, nil
	}
	for _, key := range aKeys {
		order, err = Compare(a[key], b[key])
		if err != nil {
			return 0, err
		}
		if order != 0 {
			return order, nil
		}
	}
	return 0, nil
}

// Determine if the specified normalized values can be considered to be equal
func Equals(a, b any) (bool, jpl.JPLError) {
	c, err := Compare(a, b)
	if err != nil {
		return false, err
	}
	return c == 0, nil
}

// Deep merge the specified normalized values
func Merge(a, b any) (any, jpl.JPLError) {
	ta, err := Type(a)
	if err != nil {
		return nil, err
	}
	tb, err := Type(b)
	if err != nil {
		return nil, err
	}
	if ta != "object" || tb != "object" {
		return b, nil
	}

	return AlterValue(a, jpl.JPLModifierFunc(func(input any) (output any, err jpl.JPLError) {
		value := input.(map[string]any)

		ub, err := UnwrapValue(b)
		if err != nil {
			return nil, err
		}
		changes := ObjectEntries(ub.(map[string]any))
		for _, change := range changes {
			change.Value, err = Merge(value[change.Key], change.Value)
			if err != nil {
				return nil, err
			}
		}

		return ApplyObject(value, changes), nil
	}))
}

// Stringify the specified normalized value for usage in program outputs
func StringifyJSON(value any, unescapeString bool) (string, jpl.JPLError) {
	return Stringify(value, unescapeString, false)
}

// Strip the specified normalized value for usage in program outputs
func StripJSON(value any) (any, jpl.JPLError) {
	return Strip(value, nil, nil)
}
