package builtins

import (
	"github.com/2manyvcos/jpl/go/library"
)

var native = library.MergeMaps(
	map[string]any{
		"contains":   funcContains,
		"endsWith":   funcEndsWith,
		"error":      funcError,
		"fromJSON":   funcFromJSON,
		"has":        funcHas,
		"in":         funcIn,
		"keys":       funcKeys,
		"length":     funcLength,
		"startsWith": funcStartsWith,
		"toJSON":     funcToJSON,
		"toNumber":   funcToNumber,
		"toString":   funcToString,
		"trim":       funcTrim,
		"trimEnd":    funcTrimEnd,
		"trimStart":  funcTrimStart,
		"type":       funcType,
		"void":       funcVoid,
		// "date":       nsDate, // TODO:
	},
	funcsMath,
)
