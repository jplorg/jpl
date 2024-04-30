package builtins

import (
	"slices"

	"github.com/jplorg/jpl/go/jpl"
	"github.com/jplorg/jpl/go/library"
)

func sortFn(a, b any) (int, jpl.JPLError) {
	ua, err := library.UnwrapValue(a)
	if err != nil {
		return 0, err
	}
	ub, err := library.UnwrapValue(b)
	if err != nil {
		return 0, err
	}
	return library.CompareArrays(ua.([]any)[0], ub.([]any)[0])
}

var funcSortEntries = library.NativeFunction(func(runtime jpl.JPLRuntime, input any, args ...any) ([]any, error) {
	alteredValue, err := library.AlterValue(input, jpl.JPLModifierFunc(func(value any) (any, jpl.JPLError) {
		result := library.CopySlice(value.([]any))
		var sortErr jpl.JPLError
		slices.SortFunc(result, func(a, b any) int {
			c, err := sortFn(a, b)
			if err != nil {
				if sortErr == nil {
					sortErr = err
				}
				return 0
			}
			return c
		})
		if sortErr != nil {
			return nil, sortErr
		}
		return result, nil
	}))
	if err != nil {
		return nil, err
	}
	return []any{alteredValue}, nil
})
