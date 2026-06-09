import type { JPLInstructionParams, JPLRuntimeScope } from '@/library';
import type { UserJPLInstructionParams } from '@/program/params';
import type JPLRuntime from '@/runtime';

export default {
  /** { number: number } */
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
    return next(params.number ?? 0, scope);
  },

  /** { number: number } */
  map(
    runtime: JPLRuntime,
    params: UserJPLInstructionParams,
  ): JPLInstructionParams {
    return {
      number: runtime.assertType(params.number, 'number'),
    };
  },
};
