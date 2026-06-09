import {
  OPU_MULTIPLICATION,
  OPU_NULL_COALESCENCE,
  OPU_REMAINDER,
  OPU_SET,
  OPU_SUBTRACTION,
  OPU_UPDATE,
  type OP_ACCESS,
  type OP_AND,
  type OP_ARRAY_CONSTRUCTOR,
  type OP_ASSIGNMENT,
  type OP_CALCULATION,
  type OP_COMPARISON,
  type OP_CONSTANT,
  type OP_CONSTANT_FALSE,
  type OP_CONSTANT_NULL,
  type OP_CONSTANT_TRUE,
  type OP_FUNCTION_DEFINITION,
  type OP_IF,
  type OP_INTERPOLATED_STRING,
  type OP_NEGATION,
  type OP_NOT,
  type OP_NULL_COALESCENCE,
  type OP_NUMBER,
  type OP_OBJECT_CONSTRUCTOR,
  type OP_OR,
  type OP_OUTPUT_CONCAT,
  type OP_STRING,
  type OP_TRY,
  type OP_VARIABLE,
  type OP_VARIABLE_DEFINITION,
  type OP_VOID,
  type OPA_FIELD,
  type OPA_FUNCTION,
  type OPA_ITER,
  type OPA_SLICE,
  type OPC_EQUAL,
  type OPC_GREATER,
  type OPC_GREATEREQUAL,
  type OPC_LESS,
  type OPC_LESSEQUAL,
  type OPC_UNEQUAL,
  type OPM_ADDITION,
  type OPM_DIVISION,
  type OPM_MULTIPLICATION,
  type OPM_REMAINDER,
  type OPM_SUBTRACTION,
  type OPU_ADDITION,
  type OPU_DIVISION,
} from './ops';

export const DEFINITION_VERSION_MAJOR = 1;
export const DEFINITION_VERSION_MINOR = 0;
export const DEFINITION_VERSION = `${DEFINITION_VERSION_MAJOR}.${DEFINITION_VERSION_MINOR}`;

export type JPLDefinition = {
  version: string;
  instructions: Pipe;
};

export type Pipe = JPLInstruction[];

export type JPLInstruction = {
  op:
    | typeof OP_ACCESS
    | typeof OP_AND
    | typeof OP_ARRAY_CONSTRUCTOR
    | typeof OP_ASSIGNMENT
    | typeof OP_CALCULATION
    | typeof OP_COMPARISON
    | typeof OP_CONSTANT
    | typeof OP_CONSTANT_FALSE
    | typeof OP_CONSTANT_NULL
    | typeof OP_CONSTANT_TRUE
    | typeof OP_FUNCTION_DEFINITION
    | typeof OP_IF
    | typeof OP_INTERPOLATED_STRING
    | typeof OP_NEGATION
    | typeof OP_NOT
    | typeof OP_NULL_COALESCENCE
    | typeof OP_NUMBER
    | typeof OP_OBJECT_CONSTRUCTOR
    | typeof OP_OR
    | typeof OP_OUTPUT_CONCAT
    | typeof OP_STRING
    | typeof OP_TRY
    | typeof OP_VARIABLE
    | typeof OP_VARIABLE_DEFINITION
    | typeof OP_VOID;
  params: JPLInstructionParams;
};

export type JPLInstructionParams = {
  after?: string;
  argNames?: string[];
  assignment?: JPLAssignment;
  catch?: Pipe;
  comparisons?: JPLComparison[];
  else?: Pipe;
  fields?: JPLField[];
  ifs?: JPLIfThen[];
  interpolations?: JPLInterpolation[];
  name?: string;
  number?: number;
  operations?: JPLOperation[];
  pipe?: Pipe;
  pipes?: Pipe[];
  selectors?: JPLSelector[];
  string?: string;
  try?: Pipe;
  value?: unknown;
};

export type JPLComparison = {
  op:
    | typeof OPC_EQUAL
    | typeof OPC_GREATER
    | typeof OPC_GREATEREQUAL
    | typeof OPC_LESS
    | typeof OPC_LESSEQUAL
    | typeof OPC_UNEQUAL;
  params: JPLComparisonParams;
};

export type JPLComparisonParams = { by: Pipe };

export type JPLSelector = {
  op:
    | typeof OPA_FIELD
    | typeof OPA_FUNCTION
    | typeof OPA_ITER
    | typeof OPA_SLICE;
  params: JPLSelectorParams;
};

export type JPLSelectorParams = {
  args?: Pipe[];
  bound?: boolean;
  from?: Pipe;
  optional?: boolean;
  pipe?: Pipe;
  to?: Pipe;
};

export type JPLOperation = {
  op:
    | typeof OPM_ADDITION
    | typeof OPM_DIVISION
    | typeof OPM_MULTIPLICATION
    | typeof OPM_REMAINDER
    | typeof OPM_SUBTRACTION;
  params: JPLOperationParams;
};

export type JPLOperationParams = { by: Pipe };

export type JPLIfThen = { if: Pipe; then: Pipe };

export type JPLAssignment = {
  op:
    | typeof OPU_ADDITION
    | typeof OPU_DIVISION
    | typeof OPU_MULTIPLICATION
    | typeof OPU_NULL_COALESCENCE
    | typeof OPU_REMAINDER
    | typeof OPU_SET
    | typeof OPU_SUBTRACTION
    | typeof OPU_UPDATE;
  params: JPLAssignmentParams;
};

export type JPLAssignmentParams = { pipe: Pipe };

export type JPLField = { key: Pipe; value: Pipe; optional: boolean };

export type JPLInterpolation = { before: string; pipe: Pipe };
