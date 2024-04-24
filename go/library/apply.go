package library

type Entry[Key any] struct {
	Key     Key
	NoValue bool
	Value   any
}
type ObjectEntry = Entry[string]
type ArrayEntry = Entry[int]

// TODO: this will not work like expected
// because for example slices cannot be compared
// and are not passed by pointer but are copied.
//
// One possible workaround would be a deep comparison
// but this would be very imperformant,
// the other workaround would be to pass arrays and
// objects as slice / map pointers instead of them directly,
// however that would make the transition to
// other JSON packages rather complicated and
// expensive in terms of memory.
//
// A third (and maybe best) workaround would be
// to use unsafe.SliceData
// as described here:
// https://stackoverflow.com/questions/53009686/equality-identity-of-go-slices
// However, there still probably needs to be found a way to compare maps (if they cannot simply be compared directly).
// UPDATE:
// reflect.ValueOf(a).Pointer() == reflect.ValueOf(b).Pointer()
// works for both slices and maps

// Apply all changes immutably to the source object.
func applyObject(source map[string]any, changes []ObjectEntry) map[string]any {
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

		if value, exists := result[change.Key]; exists && value == change.Value {
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
func applyArray(source []any, changes []ArrayEntry, filler any) []any {
	result := source
	unchanged := true

	for _, change := range changes {
		i := change.Key
		if change.Key < 0 {
			i = len(result) + change.Key
		}

		if i >= 0 && i < len(result) && result[i] == change.Value {
			continue
		}
		if change.Key >= 0 {
			suf := i + 1 - len(result)
			if suf > 0 {
				total := len(result) + suf
				nextResult := make([]any, total)
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
				nextResult := make([]any, total)
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
			nextResult := make([]any, len(result))
			copy(nextResult, result)
			result = nextResult
			unchanged = false
		}
		result[i] = change.Value
	}

	return result
}

type combinationIndex struct {
	count   int
	current int
	values  []any
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
func applyCombinations(combinations [][]any) []any {
	l := len(combinations)
	total := 1
	indices := make([]*combinationIndex, l)
	for i, entry := range combinations {
		count := len(entry)
		total *= count
		indices[i] = &combinationIndex{count, 0, entry}
	}
	if total == 0 {
		return nil
	}
	s := make([]any, l)
	out := make([]any, total)
	var c int
	for {
		result := s
		unchanged := true
		for i, index := range indices {
			v := index.values[index.current]
			if result[i] == v {
				return
			}
		}
	}
	return out
}

// Copies the source map
func copyMap[Value any](source map[string]Value) (result map[string]Value) {
	result = make(map[string]Value, len(source))
	for key, value := range source {
		result[key] = value
	}
	return
}
