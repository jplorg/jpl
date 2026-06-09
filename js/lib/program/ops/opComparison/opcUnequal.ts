import type { JPLComparisonParams, JPLRuntimeScope } from '@/library';
import { call } from '@/program/ops/utils';
import type { UserJPLComparisonParams } from '@/program/params';
import type JPLRuntime from '@/runtime';

export default {
  /** { by: [op] } */
  op(
    runtime: JPLRuntime,
    input: unknown,
    target: unknown,
    params: JPLComparisonParams,
    scope: JPLRuntimeScope,
    next: (output: unknown) => Promise<unknown[]> | unknown[],
  ): Promise<unknown[]> {
    return runtime.executeInstructions(params.by ?? [], [input], scope, (by) =>
      next(!runtime.equals(target, by)),
    );
  },

  /** { by: function } */
  map(
    _runtime: JPLRuntime,
    params: UserJPLComparisonParams,
  ): JPLComparisonParams {
    return {
      by: call(params.by),
    };
  },
};
