import { type JPLRuntimeScope } from '@/library';
import type { JPLInstructionParams } from '@/library/definition';
import type { UserJPLInstructionParams } from '@/program/params';
import type JPLRuntime from '@/runtime';

export default {
  /** { string: string } */
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
    return next(params.string ?? '', scope);
  },

  /** { string: string } */
  map(
    runtime: JPLRuntime,
    params: UserJPLInstructionParams,
  ): JPLInstructionParams {
    return {
      string: runtime.assertType(params.string, 'string'),
    };
  },
};
