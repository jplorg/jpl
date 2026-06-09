import {
  OPM_MULTIPLICATION,
  OP_CALCULATION,
  type JPLRuntimeScope,
} from '@/library';
import type { JPLAssignmentParams } from '@/library/definition';
import type { UserJPLAssignmentParams } from '@/program/params';
import type JPLRuntime from '@/runtime';
import { call, constant } from '../utils';

export default {
  /** { pipe: [op] } */
  op(
    runtime: JPLRuntime,
    input: unknown,
    target: unknown,
    params: JPLAssignmentParams,
    scope: JPLRuntimeScope,
    next: (output: unknown) => Promise<unknown[]> | unknown[],
  ): Promise<unknown[]> {
    return runtime.executeInstructions(
      [
        {
          op: OP_CALCULATION,
          params: {
            pipe: constant(target),
            operations: [
              { op: OPM_MULTIPLICATION, params: { by: params.pipe ?? [] } },
            ],
          },
        },
      ],
      [input],
      scope,
      next,
    );
  },

  /** { pipe: function } */
  map(
    _runtime: JPLRuntime,
    params: UserJPLAssignmentParams,
  ): JPLAssignmentParams {
    return {
      pipe: call(params.pipe),
    };
  },
};
