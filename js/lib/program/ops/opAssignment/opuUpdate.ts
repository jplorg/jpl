import type { JPLRuntimeScope } from '@/library';
import type { JPLAssignmentParams } from '@/library/definition';
import type { UserJPLAssignmentParams } from '@/program/params';
import type JPLRuntime from '@/runtime';
import { call } from '../utils';

export default {
  /** { pipe: [op] } */
  op(
    runtime: JPLRuntime,
    _input: unknown,
    target: unknown,
    params: JPLAssignmentParams,
    scope: JPLRuntimeScope,
    next: (output: unknown) => Promise<unknown[]> | unknown[],
  ): Promise<unknown[]> {
    return runtime.executeInstructions(
      params.pipe ?? [],
      [target],
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
