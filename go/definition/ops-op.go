package definition

// JPL operator type
type JPLOP string

// { value: function, selectors: [opa] }
//
// { pipe: [op], selectors: [opa] }
const OP_ACCESS = JPLOP("$.")

// { conditions: [function] }
//
// { pipes: [[op]] }
const OP_AND = JPLOP("and")

// { value: function }
//
// { pipe: [op] }
const OP_ARRAY_CONSTRUCTOR = JPLOP("[]")

// { value: function, selectors: [opa], assignment: opu }
//
// { pipe: [op], selectors: [opa], assignment: [opu] }
const OP_ASSIGNMENT = JPLOP("$=")

// { value: function, operations: [opm] }
//
// { pipe: [op], operations: [opm] }
const OP_CALCULATION = JPLOP("clc")

// { value: function, comparisons: [opc] }
//
// { pipe: [op], comparisons: [opc] }
const OP_COMPARISON = JPLOP("cmp")

// { value: any }
//
// { value: any }
const OP_CONSTANT = JPLOP("cnt")

// {}
//
// {}
const OP_CONSTANT_FALSE = JPLOP("fls")

// {}
//
// {}
const OP_CONSTANT_NULL = JPLOP("nul")

// {}
//
// {}
const OP_CONSTANT_TRUE = JPLOP("tru")

// { argNames: [string], value: [op] }
//
// { argNames: [string], pipe: [op] }
const OP_FUNCTION_DEFINITION = JPLOP("fun")

// { ifs: [{ if: function, then: function }], else: function }
//
// { ifs: [{ if: [op], then: [op] }], else: [op] }
const OP_IF = JPLOP("if")

// { interpolations: [{ before: string, value: function }], after: string }
//
// { interpolations: [{ before: string, pipe: [op] }], after: string }
const OP_INTERPOLATED_STRING = JPLOP(`"$"`)

// {}
//
// {}
const OP_NEGATION = JPLOP("neg")

// {}
//
// {}
const OP_NOT = JPLOP("not")

// { values: [function] }
//
// { pipes: [[op]] }
const OP_NULL_COALESCENCE = JPLOP("??")

// { value: number }
//
// { value: string }
const OP_NUMBER = JPLOP("nbr")

// { fields: [{ key: function, value: function, optional: boolean }] }
//
// { fields: [{ key: [op], value: [op], optional: boolean }] }
const OP_OBJECT_CONSTRUCTOR = JPLOP("{}")

// { conditions: [function] }
//
// { pipes: [[op]] }
const OP_OR = JPLOP("or")

// { values: [function] }
//
// { pipes: [[op]] }
const OP_OUTPUT_CONCAT = JPLOP(",")

// { value: string }
//
// { value: string }
const OP_STRING = JPLOP(`""`)

// { try: function, catch: function }
//
// { try: [op], catch: [op] }
const OP_TRY = JPLOP("try")

// { name: string }
//
// { name: string }
const OP_VARIABLE = JPLOP("var")

// { name: string, value: function }
//
// { name: string, pipe: [op] }
const OP_VARIABLE_DEFINITION = JPLOP("va=")

// {}
//
// {}
const OP_VOID = JPLOP("vod")
