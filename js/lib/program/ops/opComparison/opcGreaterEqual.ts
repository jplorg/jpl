import type { JPLRuntimeScope } from '@/library';
import type { JPLComparisonParams } from '@/library/definition';
import type { UserJPLComparisonParams } from '@/program/params';
import type JPLRuntime from '@/runtime';
import { call } from '../utils';

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
      next(runtime.compare(target, by) >= 0),
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
