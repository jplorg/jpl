package definition

// { value: function, operations: [opa] }
//
// { pipe: [op], operations: [opa] }
const OP_ACCESS = "$."

// { conditions: [function] }
//
// { pipes: [[op]] }
const OP_AND = "and"

// { value: function }
//
// { pipe: [op] }
const OP_ARRAY_CONSTRUCTOR = "[]"

// { value: function, operations: [opa], assignment: opu }
//
// { pipe: [op], operations: [opa], assignment: [opu] }
const OP_ASSIGNMENT = "$="

// { value: function, operations: [opm] }
//
// { pipe: [op], operations: [opm] }
const OP_CALCULATION = "clc"

// { value: function, comparisons: [opc] }
//
// { pipe: [op], comparisons: [opc] }
const OP_COMPARISON = "cmp"

// { value: any }
//
// { value: any }
const OP_CONSTANT = "cnt"

// {}
//
// {}
const OP_CONSTANT_FALSE = "fls"

// {}
//
// {}
const OP_CONSTANT_NULL = "nul"

// {}
//
// {}
const OP_CONSTANT_TRUE = "tru"

// { argNames: [string], value: [op] }
//
// { argNames: [string], pipe: [op] }
const OP_FUNCTION_DEFINITION = "fun"

// { ifs: [{ if: function, then: function }], else: function }
//
// { ifs: [{ if: [op], then: [op] }], else: [op] }
const OP_IF = "if"

// { interpolations: [{ before: string, value: function }], after: string }
//
// { interpolations: [{ before: string, pipe: [op] }], after: string }
const OP_INTERPOLATED_STRING = `"$"`

// {}
//
// {}
const OP_NEGATION = "neg"

// {}
//
// {}
const OP_NOT = "not"

// { values: [function] }
//
// { pipes: [[op]] }
const OP_NULL_COALESCENCE = "??"

// { value: number }
//
// { value: string }
const OP_NUMBER = "nbr"

// { fields: [{ key: function, value: function, optional: boolean }] }
//
// { fields: [{ key: [op], value: [op], optional: boolean }] }
const OP_OBJECT_CONSTRUCTOR = "{}"

// { conditions: [function] }
//
// { pipes: [[op]] }
const OP_OR = "or"

// { values: [function] }
//
// { pipes: [[op]] }
const OP_OUTPUT_CONCAT = ","

// { value: string }
//
// { value: string }
const OP_STRING = `""`

// { try: function, catch: function }
//
// { try: [op], catch: [op] }
const OP_TRY = "try"

// { name: string }
//
// { name: string }
const OP_VARIABLE = "var"

// { name: string, value: function }
//
// { name: string, pipe: [op] }
const OP_VARIABLE_DEFINITION = "va="

// {}
//
// {}
const OP_VOID = "vod"
