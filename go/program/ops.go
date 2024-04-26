package program

import (
	"github.com/2manyvcos/jpl/go/definition"
	"github.com/2manyvcos/jpl/go/jpl"
)

var ops = map[definition.JPLOP]jpl.JPLOPHandler{
	definition.OP_ACCESS:              opAccess{},
	definition.OP_AND:                 opAnd{},
	definition.OP_ARRAY_CONSTRUCTOR:   opArrayConstructor{},
	definition.OP_ASSIGNMENT:          opAssignment{},
	definition.OP_CALCULATION:         opCalculation{},
	definition.OP_COMPARISON:          opComparison{},
	definition.OP_CONSTANT:            opConstant{},
	definition.OP_CONSTANT_FALSE:      opConstantFalse{},
	definition.OP_CONSTANT_NULL:       opConstantNull{},
	definition.OP_CONSTANT_TRUE:       opConstantTrue{},
	definition.OP_FUNCTION_DEFINITION: opFunctionDefinition{},
	definition.OP_IF:                  opIf{},
	definition.OP_INTERPOLATED_STRING: opInterpolatedString{},
	definition.OP_NEGATION:            opNegation{},
	definition.OP_NOT:                 opNot{},
	definition.OP_NULL_COALESCENCE:    opNullCoalescence{},
	definition.OP_NUMBER:              opNumber{},
	definition.OP_OBJECT_CONSTRUCTOR:  opObjectConstructor{},
	definition.OP_OR:                  opOr{},
	definition.OP_OUTPUT_CONCAT:       opOutputConcat{},
	definition.OP_STRING:              opString{},
	definition.OP_TRY:                 opTry{},
	definition.OP_VARIABLE:            opVariable{},
	definition.OP_VARIABLE_DEFINITION: opVariableDefinition{},
	definition.OP_VOID:                opVoid{},
}
