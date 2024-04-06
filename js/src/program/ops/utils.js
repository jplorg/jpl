import { OPA_FUNCTION, OP_ACCESS, OP_CONSTANT, assertType } from '../../library';

/** Create OP_CONSTANT operation for the specified normalized value */
export function constant(value) {
  return [{ op: OP_CONSTANT, params: { value } }];
}

/** Create OPA_FUNCTION operation for the specified function */
export function call(value) {
  return [
    {
      op: OP_ACCESS,
      params: {
        pipe: constant(assertType(value, 'function')),
        operations: [{ op: OPA_FUNCTION, params: { args: [], bound: false, optional: false } }],
      },
    },
  ];
}
