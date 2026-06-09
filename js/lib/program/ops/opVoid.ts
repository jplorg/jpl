import type { JPLRuntimeScope } from '@/library';
import type { JPLInstructionParams } from '@/library/definition';
import type JPLRuntime from '@/runtime';
import type { UserJPLInstructionParams } from '../params';

export default {
  /** {} */
  op(
    _runtime: JPLRuntime,
    _input: unknown,
    _params: JPLInstructionParams,
    _scope: JPLRuntimeScope,
    _next: (
      output: unknown,
      scope: JPLRuntimeScope,
    ) => Promise<unknown[]> | unknown[],
  ): Promise<unknown[]> | unknown[] {
    return [];
  },

  /** {} */
  map(
    _runtime: JPLRuntime,
    _params: UserJPLInstructionParams,
  ): JPLInstructionParams {
    return {};
  },
};
