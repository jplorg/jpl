import type { JPLInstructionParams, JPLRuntimeScope } from '@/library';
import type { UserJPLInstructionParams } from '@/program/params';
import type JPLRuntime from '@/runtime';

export default {
  /** { value: any } */
  op(
    _runtime: JPLRuntime,
    _input: unknown,
    params: JPLInstructionParams,
    scope: JPLRuntimeScope,
    next: (
      output: unknown,
      scope: JPLRuntimeScope,
    ) => Promise<unknown[]> | unknown[],
  ): Promise<unknown[]> | unknown[] {
    return next(params.value ?? null, scope);
  },

  /** { value: any } */
  map(
    runtime: JPLRuntime,
    params: UserJPLInstructionParams,
  ): JPLInstructionParams {
    return {
      value: runtime.normalizeValue(params.value),
    };
  },
};
