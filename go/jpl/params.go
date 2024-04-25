package jpl

type JPLInstructionParams struct {
	After    string
	ArgNames []string
	// Assignment     *JPLAssignment
	Catch JPLFunc
	// Comparisons    []JPLComparison
	Else           JPLFunc
	Fields         []JPLField
	Ifs            []JPLIfThen
	Interpolations []JPLInterpolation
	Name           string
	// Operations     []JPLOperation
	Pipe  JPLFunc
	Pipes []JPLFunc
	// Selectors      []JPLSelector
	Try    JPLFunc
	Value  any
	String string
	Number float64
}

type JPLIfThen struct {
	If   JPLFunc
	Then JPLFunc
}

type JPLField struct {
	Key      JPLFunc
	Value    JPLFunc
	Optional bool
}

type JPLInterpolation struct {
	Before string
	Pipe   JPLFunc
}
