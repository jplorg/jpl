package library

type iteratorIndex[Value any] struct {
	max, current int
	values       []Value
}

// Copy the specified map
func CopyMap[Key comparable, Value any](source map[Key]Value) (result map[Key]Value) {
	result = make(map[Key]Value, len(source))
	for key, value := range source {
		result[key] = value
	}
	return
}

// Return the specified map's keys
func GetMapKeys[Key comparable, Value any](source map[Key]Value) []Key {
	result := make([]Key, len(source))
	for key := range source {
		result = append(result, key)
	}
	return result
}

// Return the specified map's values
func GetMapValues[Key comparable, Value any](source map[Key]Value) []Value {
	result := make([]Value, 0, len(source))
	for _, value := range source {
		result = append(result, value)
	}
	return result
}

// Copy the specified slice
func CopySlice[Value any](source []Value) (result []Value) {
	result = make([]Value, len(source))
	copy(result, source)
	return
}

// Return a sub slice of the specified slice.
//
// `from` and `to` are allowed to be negative or to be outside of the range of `source`.
func SubSlice[Value any](source []Value, from int, to int) []Value {
	l := len(source)
	if from < 0 {
		from = l + from
	}
	if to < 0 {
		to = l + to
	}
	from = max(0, min(from, l))
	to = max(from, min(to, l))
	return source[from:to]
}
