import { type JPLRuntimeScope } from '@/library';
import type { JPLInstructionParams } from '@/library/definition';
import type { UserJPLInstructionParams } from '@/program/params';
import type JPLRuntime from '@/runtime';

export default {
  /** { value: any } */
  op(
    _runtime: JPLRuntime,
    _input: unknown,
    _params: JPLInstructionParams,
    scope: JPLRuntimeScope,
    next: (
      output: unknown,
      scope: JPLRuntimeScope,
    ) => Promise<unknown[]> | unknown[],
  ): Promise<unknown[]> | unknown[] {
    return next(null, scope);
  },

  /** { value: any } */
  map(
    _runtime: JPLRuntime,
    _params: UserJPLInstructionParams,
  ): JPLInstructionParams {
    return {};
  },
};
