package definition

// JPL sub operator type for OP_COMPARISON
type JPLOPC string

// Sub operator for OP_COMPARISON
//
// { by: function }
//
// { by: [op] }
const OPC_EQUAL = JPLOPC("==")

// Sub operator for OP_COMPARISON
//
// { by: function }
//
// { by: [op] }
const OPC_GREATER = JPLOPC(">")

// Sub operator for OP_COMPARISON
//
// { by: function }
//
// { by: [op] }
const OPC_GREATEREQUAL = JPLOPC(">=")

// Sub operator for OP_COMPARISON
//
// { by: function }
//
// { by: [op] }
const OPC_LESS = JPLOPC("<")

// Sub operator for OP_COMPARISON
//
// { by: function }
//
// { by: [op] }
const OPC_LESSEQUAL = JPLOPC("<=")

// Sub operator for OP_COMPARISON
//
// { by: function }
//
// { by: [op] }
const OPC_UNEQUAL = JPLOPC("!=")
