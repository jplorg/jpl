import type { JPLInstructionParams, JPLRuntimeScope } from '@/library';
import type { UserJPLInstructionParams } from '@/program/params';
import type JPLRuntime from '@/runtime';

export default {
  /** {} */
  op(
    runtime: JPLRuntime,
    input: unknown,
    _params: JPLInstructionParams,
    scope: JPLRuntimeScope,
    next: (
      output: unknown,
      scope: JPLRuntimeScope,
    ) => Promise<unknown[]> | unknown[],
  ): Promise<unknown[]> | unknown[] {
    return next(!runtime.truthy(input), scope);
  },

  /** {} */
  map(
    _runtime: JPLRuntime,
    _params: UserJPLInstructionParams,
  ): JPLInstructionParams {
    return {};
  },
};
