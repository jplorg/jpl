package definition

// JPL sub operator type for OP_ACCESS
type JPLOPA string

// Sub operator for OP_ACCESS
//
// { value: function, optional: boolean }
//
// { pipe: [op], optional: boolean }
const OPA_FIELD = JPLOPA("$.$")

// Sub operator for OP_ACCESS
//
// { args: [function], bound: boolean, optional: boolean }
//
// { args: [[op]], bound: boolean, optional: boolean }
const OPA_FUNCTION = JPLOPA("$()")

// Sub operator for OP_ACCESS
//
// { optional: boolean }
//
// { optional: boolean }
const OPA_ITER = JPLOPA("$[]")

// Sub operator for OP_ACCESS
//
// { from: function, to: function, optional: boolean }
//
// { from: [op], to: [op], optional: boolean }
const OPA_SLICE = JPLOPA("[:]")
