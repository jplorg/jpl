package definition

// JPL sub operator type for OP_CALCULATION
type JPLOPM string

// Sub operator for OP_CALCULATION
//
// { by: function }
//
// { by: [op] }
const OPM_ADDITION = JPLOPM("+")

// Sub operator for OP_CALCULATION
//
// { by: function }
//
// { by: [op] }
const OPM_DIVISION = JPLOPM("/")

// Sub operator for OP_CALCULATION
//
// { by: function }
//
// { by: [op] }
const OPM_MULTIPLICATION = JPLOPM("*")

// Sub operator for OP_CALCULATION
//
// { by: function }
//
// { by: [op] }
const OPM_REMAINDER = JPLOPM("%")

// Sub operator for OP_CALCULATION
//
// { by: function }
//
// { by: [op] }
const OPM_SUBTRACTION = JPLOPM("-")
