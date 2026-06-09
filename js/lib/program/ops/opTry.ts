import {
  JPLErrorEnclosure,
  JPLExecutionError,
  type JPLInstructionParams,
  type JPLRuntimeScope,
} from '@/library';
import type { UserJPLInstructionParams } from '@/program/params';
import type JPLRuntime from '@/runtime';
import { call } from './utils';

export default {
  /** { try: [op], catch: [op] } */
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
    const nextScope = scope.next({ signal: scope.signal.next() });
    try {
      return await runtime.executeInstructions(
        params.try ?? [],
        [input],
        nextScope,
        async (output) => {
          try {
            return await next(output, scope);
          } catch (err) {
            throw new JPLErrorEnclosure(err);
          }
        },
      );
    } catch (err) {
      if (JPLErrorEnclosure.is(err)) throw (err as JPLErrorEnclosure).inner;
      if (!JPLExecutionError.is(err)) throw err;
      nextScope.signal.exit();
      return runtime.executeInstructions(
        params.catch ?? [],
        [(err as JPLExecutionError).value],
        scope,
        (output) => next(output, scope),
      );
    }
  },

  /** { try: function, catch: function } */
  map(
    _runtime: JPLRuntime,
    params: UserJPLInstructionParams,
  ): JPLInstructionParams {
    return {
      try: call(params.try),
      catch: call(params.catch),
    };
  },
};
