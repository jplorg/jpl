import {
  OP_NULL_COALESCENCE,
  type JPLAssignmentParams,
  type JPLRuntimeScope,
} from '@/library';
import { call, constant } from '@/program/ops/utils';
import type { UserJPLAssignmentParams } from '@/program/params';
import type JPLRuntime from '@/runtime';

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
          op: OP_NULL_COALESCENCE,
          params: { pipes: [constant(target), params.pipe ?? []] },
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
