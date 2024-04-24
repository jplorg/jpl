package definition

// JPL sub operator type for OP_ASSIGNMENT
type JPLOPU string

// Sub operator for OP_ASSIGNMENT
//
// { pipe: function }
//
// { pipe: [op] }
const OPU_ADDITION = JPLOPU("+=")

// Sub operator for OP_ASSIGNMENT
//
// { pipe: function }
//
// { pipe: [op] }
const OPU_DIVISION = JPLOPU("/=")

// Sub operator for OP_ASSIGNMENT
//
// { pipe: function }
//
// { pipe: [op] }
const OPU_MULTIPLICATION = JPLOPU("*=")

// Sub operator for OP_ASSIGNMENT
//
// { pipe: function }
//
// { pipe: [op] }
const OPU_NULL_COALESCENCE = JPLOPU("?=")

// Sub operator for OP_ASSIGNMENT
//
// { pipe: function }
//
// { pipe: [op] }
const OPU_REMAINDER = JPLOPU("%=")

// Sub operator for OP_ASSIGNMENT
//
// { pipe: function }
//
// { pipe: [op] }
const OPU_SET = JPLOPU("=")

// Sub operator for OP_ASSIGNMENT
//
// { pipe: function }
//
// { pipe: [op] }
const OPU_SUBTRACTION = JPLOPU("-=")

// Sub operator for OP_ASSIGNMENT
//
// { pipe: function }
//
// { pipe: [op] }
const OPU_UPDATE = JPLOPU("|=")
