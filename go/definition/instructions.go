package definition

type Pipe = []JPLInstruction

type JPLInstruction struct {
	OP     string               `json:"op"`
	Params JPLInstructionParams `json:"params"`
}

type JPLInstructionParams struct {
	Pipes          []Pipe             `json:"params,omitempty"`
	Try            Pipe               `json:"try,omitempty"`
	Catch          Pipe               `json:"catch,omitempty"`
	Pipe           Pipe               `json:"pipe,omitempty"`
	Comparisons    []JPLComparison    `json:"comparisons,omitempty"`
	Operations     []JPLOperation     `json:"operations,omitempty"`
	Ifs            []JPLIfThen        `json:"ifs,omitempty"`
	Else           Pipe               `json:"else,omitempty"`
	Name           string             `json:"name,omitempty"`
	ArgNames       []string           `json:"argNames,omitempty"`
	Assignment     *JPLAssignment     `json:"assignment,omitempty"`
	Fields         []JPLField         `json:"fields,omitempty"`
	Value          string             `json:"value,omitempty"`
	Interpolations []JPLInterpolation `json:"interpolations,omitempty"`
	After          string             `json:"after,omitempty"`
}

type JPLComparison struct {
	OP     string              `json:"op"`
	Params JPLComparisonParams `json:"params"`
}

type JPLComparisonParams struct {
	By Pipe `json:"by"`
}

type JPLOperation struct {
	OP     string             `json:"op"`
	Params JPLOperationParams `json:"params"`
}

type JPLOperationParams struct {
	By       Pipe   `json:"by,omitempty"`
	Pipe     Pipe   `json:"pipe,omitempty"`
	Optional bool   `json:"optional"`
	From     Pipe   `json:"from,omitempty"`
	To       Pipe   `json:"to,omitempty"`
	Args     []Pipe `json:"args,omitempty"`
	Bound    bool   `json:"bound"`
}

type JPLIfThen struct {
	If   Pipe `json:"if"`
	Then Pipe `json:"then"`
}

type JPLAssignment struct {
	OP     string              `json:"op"`
	Params JPLAssignmentParams `json:"params"`
}

type JPLAssignmentParams struct {
	Pipe Pipe `json:"pipe"`
}

type JPLField struct {
	Key      Pipe `json:"key"`
	Value    Pipe `json:"value"`
	Optional bool `json:"optional"`
}

type JPLInterpolation struct {
	Before string `json:"before"`
	Pipe   Pipe   `json:"pipe"`
}
