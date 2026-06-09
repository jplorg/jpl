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
    const iter = async (from: number): Promise<unknown[]> => {
      // Call stack decoupling - This is necessary as some browsers (i.e. Safari) have very limited call stack sizes which result in stack overflow exceptions in certain situations.
      await undefined;

      scope.signal.checkHealth();

      if (from >= (params.pipes?.length ?? 0)) return next(false, scope);

      const pipe = params.pipes![from];

      return runtime.executeInstructions(
        pipe ?? [],
        [input],
        scope,
        (output) => {
          if (runtime.truthy(output)) return next(true, scope);

          return iter(from + 1);
        },
      );
    };

    return iter(0);
  },

  /** { pipes: [function] } */
  map(
    runtime: JPLRuntime,
    params: UserJPLInstructionParams,
  ): JPLInstructionParams {
    return {
      pipes: runtime.muxOne([params.pipes!], (pipe) => call(pipe)),
    };
  },
};
