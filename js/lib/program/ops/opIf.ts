import type { JPLRuntimeScope } from '@/library';
import type { JPLInstructionParams } from '@/library/definition';
import type JPLRuntime from '@/runtime';
import type { UserJPLInstructionParams } from '../params';
import { call } from './utils';

export default {
  /** { ifs: [{ if: [op], then: [op] }], else: [op] } */
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

      if (from >= (params.ifs?.length ?? 0)) {
        return runtime.executeInstructions(
          params.else ?? [],
          [input],
          scope,
          (output) => next(output, scope),
        );
      }

      const condition = params.ifs![from];

      return runtime.executeInstructions(
        condition.if ?? [],
        [input],
        scope,
        (result) => {
          if (runtime.truthy(result)) {
            return runtime.executeInstructions(
              condition.then ?? [],
              [input],
              scope,
              (output) => next(output, scope),
            );
          }

          return iter(from + 1);
        },
      );
    };

    return iter(0);
  },

  /** { ifs: [{ if: function, then: function }], else: function } */
  map(
    runtime: JPLRuntime,
    params: UserJPLInstructionParams,
  ): JPLInstructionParams {
    return {
      ifs: runtime.muxOne([params.ifs!], (entry) => ({
        if: call(entry.if),
        then: call(entry.then),
      })),
      else: call(params.else),
    };
  },
};
