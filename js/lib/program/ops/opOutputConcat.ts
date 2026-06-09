import type { JPLRuntimeScope } from '@/library';
import type { JPLInstructionParams } from '@/library/definition';
import type JPLRuntime from '@/runtime';
import type { UserJPLInstructionParams } from '../params';
import { call } from './utils';

export default {
  /** { pipes: [[op]] } */
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
    return runtime.muxAll([params.pipes ?? []], (pipe) =>
      runtime.executeInstructions(pipe ?? [], [input], scope, (output) =>
        next(output, scope),
      ),
    );
  },

  /** { pipes: [function] } */
  map(
    runtime: JPLRuntime,
    params: UserJPLInstructionParams,
  ): JPLInstructionParams {
    return {
      pipes: runtime.muxOne([params.pipes!], (value) => call(value)),
    };
  },
};
