import type { JPLFunc } from '@/library';
import {
  OPU_MULTIPLICATION,
  OPU_NULL_COALESCENCE,
  OPU_REMAINDER,
  OPU_SET,
  OPU_SUBTRACTION,
  OPU_UPDATE,
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
} from '@/library/ops';

export type UserJPLInstructionParams = {
  after?: string;
  argNames?: string[];
  assignment?: UserJPLAssignment;
  catch?: JPLFunc;
  comparisons?: UserJPLComparison[];
  else?: JPLFunc;
  fields?: UserJPLField[];
  ifs?: UserJPLIfThen[];
  interpolations?: UserJPLInterpolation[];
  name?: string;
  number?: number;
  operations?: UserJPLOperation[];
  pipe?: JPLFunc;
  pipes?: JPLFunc[];
  selectors?: UserJPLSelector[];
  string?: string;
  try?: JPLFunc;
  value?: unknown;
};

export type UserJPLComparison = {
  op:
    | typeof OPC_EQUAL
    | typeof OPC_GREATER
    | typeof OPC_GREATEREQUAL
    | typeof OPC_LESS
    | typeof OPC_LESSEQUAL
    | typeof OPC_UNEQUAL;
  params: UserJPLComparisonParams;
};

export type UserJPLComparisonParams = { by: JPLFunc };

export type UserJPLSelector = {
  op:
    | typeof OPA_FIELD
    | typeof OPA_FUNCTION
    | typeof OPA_ITER
    | typeof OPA_SLICE;
  params: UserJPLSelectorParams;
};

export type UserJPLSelectorParams = {
  args?: JPLFunc[];
  bound?: boolean;
  from?: JPLFunc;
  optional?: boolean;
  pipe?: JPLFunc;
  to?: JPLFunc;
};

export type UserJPLOperation = {
  op:
    | typeof OPM_ADDITION
    | typeof OPM_DIVISION
    | typeof OPM_MULTIPLICATION
    | typeof OPM_REMAINDER
    | typeof OPM_SUBTRACTION;
  params: UserJPLOperationParams;
};

export type UserJPLOperationParams = { by: JPLFunc };

export type UserJPLIfThen = { if: JPLFunc; then: JPLFunc };

export type UserJPLAssignment = {
  op:
    | typeof OPU_ADDITION
    | typeof OPU_DIVISION
    | typeof OPU_MULTIPLICATION
    | typeof OPU_NULL_COALESCENCE
    | typeof OPU_REMAINDER
    | typeof OPU_SET
    | typeof OPU_SUBTRACTION
    | typeof OPU_UPDATE;
  params: UserJPLAssignmentParams;
};

export type UserJPLAssignmentParams = { pipe: JPLFunc };

export type UserJPLField = { key: JPLFunc; value: JPLFunc; optional: boolean };

export type UserJPLInterpolation = { before: string; pipe: JPLFunc };
