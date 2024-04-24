/**
 * { value: function, selectors: [opa] }
 *
 * { pipe: [op], selectors: [opa] }
 */
export const OP_ACCESS = '$.';

/**
 * { conditions: [function] }
 *
 * { pipes: [[op]] }
 */
export const OP_AND = 'and';

/**
 * { value: function }
 *
 * { pipe: [op] }
 */
export const OP_ARRAY_CONSTRUCTOR = '[]';

/**
 * { value: function, selectors: [opa], assignment: opu }
 *
 * { pipe: [op], selectors: [opa], assignment: [opu] }
 */
export const OP_ASSIGNMENT = '$=';

/**
 * { value: function, operations: [opm] }
 *
 * { pipe: [op], operations: [opm] }
 */
export const OP_CALCULATION = 'clc';

/**
 * { value: function, comparisons: [opc] }
 *
 * { pipe: [op], comparisons: [opc] }
 */
export const OP_COMPARISON = 'cmp';

/**
 * { value: any }
 *
 * { value: any }
 */
export const OP_CONSTANT = 'cnt';

/**
 * {}
 *
 * {}
 */
export const OP_CONSTANT_FALSE = 'fls';

/**
 * {}
 *
 * {}
 */
export const OP_CONSTANT_NULL = 'nul';

/**
 * {}
 *
 * {}
 */
export const OP_CONSTANT_TRUE = 'tru';

/**
 * { argNames: [string], value: [op] }
 *
 * { argNames: [string], pipe: [op] }
 */
export const OP_FUNCTION_DEFINITION = 'fun';

/**
 * { ifs: [{ if: function, then: function }], else: function }
 *
 * { ifs: [{ if: [op], then: [op] }], else: [op] }
 */
export const OP_IF = 'if';

/**
 * { interpolations: [{ before: string, value: function }], after: string }
 *
 * { interpolations: [{ before: string, pipe: [op] }], after: string }
 */
export const OP_INTERPOLATED_STRING = '"$"';

/**
 * {}
 *
 * {}
 */
export const OP_NEGATION = 'neg';

/**
 * {}
 *
 * {}
 */
export const OP_NOT = 'not';

/**
 * { values: [function] }
 *
 * { pipes: [[op]] }
 */
export const OP_NULL_COALESCENCE = '??';

/**
 * { value: number }
 *
 * { value: string }
 */
export const OP_NUMBER = 'nbr';

/**
 * { fields: [{ key: function, value: function, optional: boolean }] }
 *
 * { fields: [{ key: [op], value: [op], optional: boolean }] }
 */
export const OP_OBJECT_CONSTRUCTOR = '{}';

/**
 * { conditions: [function] }
 *
 * { pipes: [[op]] }
 */
export const OP_OR = 'or';

/**
 * { values: [function] }
 *
 * { pipes: [[op]] }
 */
export const OP_OUTPUT_CONCAT = ',';

/**
 * { value: string }
 *
 * { value: string }
 */
export const OP_STRING = '""';

/**
 * { try: function, catch: function }
 *
 * { try: [op], catch: [op] }
 */
export const OP_TRY = 'try';

/**
 * { name: string }
 *
 * { name: string }
 */
export const OP_VARIABLE = 'var';

/**
 * { name: string, value: function }
 *
 * { name: string, pipe: [op] }
 */
export const OP_VARIABLE_DEFINITION = 'va=';

/**
 * {}
 *
 * {}
 */
export const OP_VOID = 'vod';
