import type { JPLInstructionParams, JPLRuntimeScope } from '@/library';
import type { UserJPLInstructionParams } from '@/program/params';
import type JPLRuntime from '@/runtime';
import { call } from './utils';

export default {
  /** { pipe: [op] } */
  async op(
    runtime: JPLRuntime,
    input: unknown,
    params: JPLInstructionParams,
    scope: JPLRuntimeScope,
    next: (
      output: unknown,
      scope: JPLRuntimeScope,
    ) => Promise<unknown[]> | unknown[],
  ): Promise<unknown[]> {
    const outputs = await runtime.executeInstructions(
      params.pipe ?? [],
      [input],
      scope,
    );
    return next(outputs, scope);
  },

  /** { pipe: function } */
  map(
    _runtime: JPLRuntime,
    params: UserJPLInstructionParams,
  ): JPLInstructionParams {
    return {
      pipe: call(params.pipe),
    };
  },
};
