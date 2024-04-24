package definition

// Sub operator for OP_ACCESS
//
// { value: function, optional: boolean }
//
// { pipe: [op], optional: boolean }
const OPA_FIELD = "$.$"

// Sub operator for OP_ACCESS
//
// { args: [function], bound: boolean, optional: boolean }
//
// { args: [[op]], bound: boolean, optional: boolean }
const OPA_FUNCTION = "$()"

// Sub operator for OP_ACCESS
//
// { optional: boolean }
//
// { optional: boolean }
const OPA_ITER = "$[]"

// Sub operator for OP_ACCESS
//
// { from: function, to: function, optional: boolean }
//
// { from: [op], to: [op], optional: boolean }
const OPA_SLICE = "[:]"
