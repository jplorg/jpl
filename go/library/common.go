package library

type iteratorIndex[Value any] struct {
	max, current int
	values       []Value
}

func copyMap[Value any](source map[string]Value) (result map[string]Value) {
	result = make(map[string]Value, len(source))
	for key, value := range source {
		result[key] = value
	}
	return
}

func mapKeys[Value any](source map[string]Value) []string {
	result := make([]string, len(source))
	for key := range source {
		result = append(result, key)
	}
	return result
}

func mapValues[Value any](source map[string]Value) []Value {
	result := make([]Value, len(source))
	for _, value := range source {
		result = append(result, value)
	}
	return result
}

func copySlice[Value any](source []Value) (result []Value) {
	result = make([]Value, len(source))
	copy(result, source)
	return
}
