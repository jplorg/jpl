import type { JPLRuntimeScope } from '@/library';
import type { JPLInstructionParams } from '@/library/definition';
import type JPLRuntime from '@/runtime';
import type { UserJPLInstructionParams } from '../params';
import { call } from './utils';

export default {
  /** { name: string, pipe: [op] } */
  op(
    runtime: JPLRuntime,
    input: unknown,
    params: JPLInstructionParams,
    scope: JPLRuntimeScope,
    next: (
      output: unknown,
      scope: JPLRuntimeScope,
    ) => Promise<unknown[]> | unknown[],
  ): Promise<unknown[]> {
    return runtime.executeInstructions(
      params.pipe ?? [],
      [input],
      scope,
      (output) =>
        next(input, scope.next({ vars: { [params.name ?? '']: output } })),
    );
  },

  /** { name: string, pipe: function } */
  map(
    runtime: JPLRuntime,
    params: UserJPLInstructionParams,
  ): JPLInstructionParams {
    return {
      name: runtime.assertType(params.name, 'string'),
      pipe: call(params.pipe),
    };
  },
};
