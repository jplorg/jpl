import {
  OPA_FUNCTION,
  OP_ACCESS,
  OP_CONSTANT,
  assertType,
  type Pipe,
} from '@/library';

/** Create OP_CONSTANT operation for the specified normalized value */
export function constant(value: unknown): Pipe {
  return [{ op: OP_CONSTANT, params: { value: value ?? null } }];
}

/** Create OPA_FUNCTION operation for the specified function */
export function call(value: unknown): Pipe {
  return [
    {
      op: OP_ACCESS,
      params: {
        pipe: constant(assertType(value, 'function')),
        selectors: [
          {
            op: OPA_FUNCTION,
            params: { args: [], bound: false, optional: false },
          },
        ],
      },
    },
  ];
}
