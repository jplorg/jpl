package jpl

import "github.com/2manyvcos/jpl/go/definition"

type JPLInstructionParams struct {
	After    string
	ArgNames []string
	// Assignment     *JPLAssignment
	Catch          JPLFunc
	Comparisons    []JPLComparison
	Else           JPLFunc
	Fields         []JPLField
	Ifs            []JPLIfThen
	Interpolations []JPLInterpolation
	Name           string
	Operations     []JPLOperation
	Pipe           JPLFunc
	Pipes          []JPLFunc
	// Selectors      []JPLSelector
	Try    JPLFunc
	Value  any
	String string
	Number float64
}

type JPLComparison struct {
	OP     definition.JPLOPC
	Params JPLComparisonParams
}

type JPLComparisonParams struct {
	By JPLFunc
}

type JPLOperation struct {
	OP     definition.JPLOPM
	Params JPLOperationParams
}

type JPLOperationParams struct {
	By JPLFunc
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
