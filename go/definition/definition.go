package definition

const DEFINITION_VERSION_MAJOR = 1
const DEFINITION_VERSION_MINOR = 0
const DEFINITION_VERSION = "1.0"

type JPLDefinition struct {
	Version      string           `json:"version"`
	Instructions []JPLInstruction `json:"instructions"`
}

type Pipe = []JPLInstruction

type JPLInstruction struct {
	OP     JPLOP                `json:"op"`
	Params JPLInstructionParams `json:"params"`
}

type JPLInstructionParams struct {
	After          string             `json:"after,omitempty"`
	ArgNames       []string           `json:"argNames,omitempty"`
	Assignment     *JPLAssignment     `json:"assignment,omitempty"`
	Catch          Pipe               `json:"catch,omitempty"`
	Comparisons    []JPLComparison    `json:"comparisons,omitempty"`
	Else           Pipe               `json:"else,omitempty"`
	Fields         []JPLField         `json:"fields,omitempty"`
	Ifs            []JPLIfThen        `json:"ifs,omitempty"`
	Interpolations []JPLInterpolation `json:"interpolations,omitempty"`
	Name           string             `json:"name,omitempty"`
	Operations     []JPLOperation     `json:"operations,omitempty"`
	Pipe           Pipe               `json:"pipe,omitempty"`
	Pipes          []Pipe             `json:"params,omitempty"`
	Selectors      []JPLSelector      `json:"selectors,omitempty"`
	Try            Pipe               `json:"try,omitempty"`
	Value          any                `json:"value,omitempty"`
	String         string             `json:"string,omitempty"`
	Number         float64            `json:"number,omitempty"`
}

type JPLComparison struct {
	OP     JPLOPC              `json:"op"`
	Params JPLComparisonParams `json:"params"`
}

type JPLComparisonParams struct {
	By Pipe `json:"by"`
}

type JPLSelector struct {
	OP     JPLOPA            `json:"op"`
	Params JPLSelectorParams `json:"params"`
}

type JPLSelectorParams struct {
	Args     []Pipe `json:"args,omitempty"`
	Bound    bool   `json:"bound,omitempty"`
	From     Pipe   `json:"from,omitempty"`
	Optional bool   `json:"optional,omitempty"`
	Pipe     Pipe   `json:"pipe,omitempty"`
	To       Pipe   `json:"to,omitempty"`
}

type JPLOperation struct {
	OP     JPLOPM             `json:"op"`
	Params JPLOperationParams `json:"params"`
}

type JPLOperationParams struct {
	By Pipe `json:"by,omitempty"`
}

type JPLIfThen struct {
	If   Pipe `json:"if"`
	Then Pipe `json:"then"`
}

type JPLAssignment struct {
	OP     JPLOPU              `json:"op"`
	Params JPLAssignmentParams `json:"params"`
}

type JPLAssignmentParams struct {
	Pipe Pipe `json:"pipe"`
}

type JPLField struct {
	Key      Pipe `json:"key"`
	Value    Pipe `json:"value"`
	Optional bool `json:"optional,omitempty"`
}

type JPLInterpolation struct {
	Before string `json:"before"`
	Pipe   Pipe   `json:"pipe"`
}
