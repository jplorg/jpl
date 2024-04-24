package ops

import (
	"github.com/2manyvcos/jpl/go/definition"
	"github.com/2manyvcos/jpl/go/jpl"
)

var OPs = map[definition.JPLOP]jpl.JPLOPHandler{
	// TODO:
	definition.OP_VARIABLE_DEFINITION: opVariableDefinition{},
	definition.OP_VOID:                opVoid{},
}
