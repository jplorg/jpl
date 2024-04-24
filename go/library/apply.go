package library

type ObjectEntry[Value any] struct {
	Key     string
	Value   Value
	NoValue bool
}

type ArrayEntry[Value any] struct {
	Index int
	Value Value
}

func ObjectEntries[Value any](source map[string]Value) []ObjectEntry[Value] {
	result := make([]ObjectEntry[Value], 0, len(source))
	for key, value := range source {
		result = append(result, ObjectEntry[Value]{key, value, false})
	}
	return result
}

func ObjectFromEntries[Value any](entries []ObjectEntry[Value]) map[string]Value {
	result := make(map[string]Value, len(entries))
	for _, entry := range entries {
		if entry.NoValue {
			delete(result, entry.Key)
		} else {
			result[entry.Key] = entry.Value
		}
	}
	return result
}

func ArrayEntries[Value any](source []Value) []ArrayEntry[Value] {
	result := make([]ArrayEntry[Value], 0, len(source))
	for i, value := range source {
		result = append(result, ArrayEntry[Value]{i, value})
	}
	return result
}

// Apply all changes immutably to the source object.
func ApplyObject[Value any](source map[string]Value, changes []ObjectEntry[Value]) map[string]Value {
	result := source
	unchanged := true

	for _, change := range changes {
		if change.NoValue {
			if _, exists := result[change.Key]; !exists {
				continue
			}
			if unchanged {
				result = copyMap(result)
				delete(result, change.Key)
				unchanged = false
			} else {
				delete(result, change.Key)
			}
			continue
		}

		if value, exists := result[change.Key]; exists && IsSame(value, change.Value) {
			continue
		}
		if unchanged {
			result = copyMap(result)
			result[change.Key] = change.Value
			unchanged = false
		} else {
			result[change.Key] = change.Value
		}
	}

	return result
}

// Apply all changes immutably to the source array.
// Indices can be negative to be applied from the end of the array.
func ApplyArray[Value any](source []Value, changes []ArrayEntry[Value], filler Value) []Value {
	result := source
	unchanged := true

	for _, change := range changes {
		i := change.Index
		if change.Index < 0 {
			i = len(result) + change.Index
		}

		if i >= 0 && i < len(result) && IsSame(result[i], change.Value) {
			continue
		}
		if change.Index >= 0 {
			suf := i + 1 - len(result)
			if suf > 0 {
				total := len(result) + suf
				nextResult := make([]Value, total)
				for i := copy(nextResult, result); i < total; i++ {
					nextResult[i] = filler
				}
				result = nextResult
				unchanged = false
			}
		} else {
			pre := -i
			if pre > 0 {
				total := pre + len(result)
				nextResult := make([]Value, total)
				for i := 0; i < pre; i++ {
					nextResult[i] = filler
				}
				copy(nextResult[pre:], result)
				result = nextResult
				unchanged = false
				i = 0
			}
		}
		if unchanged {
			result = copySlice(result)
			unchanged = false
		}
		result[i] = change.Value
	}

	return result
}

// Create all possible combinations immutably.
// If the specified `combinations` array is empty, the resulting array contains a single empty array.
// This function has essentially the same base functionality as the `mux` function, but uses a more performant approach for generating immutable arrays as it reduces the number of necessary array copies.
//
// `applyCombinations([], [[1, 2], [3, 4]])` for example produces:
// - `[1, 3]`
// - `[1, 4]`
// - `[2, 3]`
// - `[2, 4]`
//
// If the values of `source` are equal to the values of one of the combinations, it is used instead of a copy in the output array, e.g.:
// `let i = [1, 2]; applyCombinations(i, [[1], [2]])[0] == i`
// - `true`
func ApplyCombinations[Value any](source []Value, combinations [][]Value) [][]Value {
	l := len(combinations)
	total := 1
	indices := make([]*iteratorIndex[Value], l)
	for i, entry := range combinations {
		count := len(entry)
		total *= count
		indices[i] = &iteratorIndex[Value]{count, 0, entry}
	}
	if total == 0 {
		return nil
	}
	s := source
	if sl := len(s); sl > l {
		s = s[:l]
	} else if sl < l {
		s = make([]Value, l)
		copy(s, source)
	}
	out := make([][]Value, total)
	var c int
	for {
		result := s
		unchanged := true
		for i, index := range indices {
			v := index.values[index.current]
			if IsSame(result[i], v) {
				continue
			}
			if unchanged {
				result = copySlice(result)
				unchanged = false
			}
			result[i] = v
		}
		out[c] = result
		c += 1
		if c >= total {
			break
		}
		for n := l - 1; n >= 0; n -= 1 {
			i := indices[n]
			next := i.current + 1
			if next < i.max {
				i.current = next
				break
			}
			i.current = 0
		}
	}
	return out
}

// Copy the specified map
func copyMap[Value any](source map[string]Value) (result map[string]Value) {
	result = make(map[string]Value, len(source))
	for key, value := range source {
		result[key] = value
	}
	return
}

// Copy the specified slice
func copySlice[Value any](source []Value) (result []Value) {
	result = make([]Value, len(source))
	copy(result, source)
	return
}
