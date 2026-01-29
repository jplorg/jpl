package builtins

import "github.com/jplorg/jpl/go/v2/library"

var native = library.MergeMaps(
	map[string]any{
		"contains":    funcContains,
		"endsWith":    funcEndsWith,
		"error":       funcError,
		"fromJSON":    funcFromJSON,
		"has":         funcHas,
		"in":          funcIn,
		"keys":        funcKeys,
		"length":      funcLength,
		"now":         funcNow,
		"startsWith":  funcStartsWith,
		"toJSON":      funcToJSON,
		"toLowerCase": funcToLowerCase,
		"toNumber":    funcToNumber,
		"toString":    funcToString,
		"toUpperCase": funcToUpperCase,
		"trim":        funcTrim,
		"trimEnd":     funcTrimEnd,
		"trimStart":   funcTrimStart,
		"type":        funcType,
		"void":        funcVoid,
	},
	funcsMath,
)
