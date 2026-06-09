import {
  JPLErrorEnclosure,
  JPLTypeError,
  adaptError,
  type JPLFunc,
  type JPLRuntimeScope,
  type JPLSelectorParams,
} from '@/library';
import { call } from '@/program/ops/utils';
import type { UserJPLSelectorParams } from '@/program/params';
import type JPLRuntime from '@/runtime';

export default {
  /** { args: [[op]], bound: boolean, optional: boolean } */
  async op(
    runtime: JPLRuntime,
    input: unknown,
    target: unknown,
    params: JPLSelectorParams,
    scope: JPLRuntimeScope,
    next: (output: unknown) => Promise<unknown[]> | unknown[],
  ): Promise<unknown[]> {
    const value = runtime.unwrapValue(target);
    const t = runtime.type(value);
    switch (t) {
      case 'function': {
        const args = await runtime.muxAsync([params.args ?? []], (arg) =>
          runtime.executeInstructions(arg ?? [], [input], scope),
        );

        return runtime.muxAll(args, async (...a) => {
          try {
            const fnNext = async (output: unknown): Promise<unknown[]> => {
              try {
                return await next(output);
              } catch (err) {
                throw new JPLErrorEnclosure(err);
              }
            };

            if (params.bound) {
              return await (value as JPLFunc)(
                runtime,
                scope.signal,
                fnNext,
                a[0] ?? null,
                ...a.slice(1),
              );
            }

            return await (value as JPLFunc)(
              runtime,
              scope.signal,
              fnNext,
              input,
              ...a,
            );
          } catch (err) {
            if (JPLErrorEnclosure.is(err))
              throw (err as JPLErrorEnclosure).inner;
            throw adaptError(err);
          }
        });
      }

      default:
    }

    if (params.optional) return [];
    throw new JPLTypeError('cannot execute %s (%*<100v)', t, value);
  },

  /** { args: [function], bound: boolean, optional: boolean } */
  map(runtime: JPLRuntime, params: UserJPLSelectorParams): JPLSelectorParams {
    return {
      args: runtime.muxOne([params.args!], (entry) => call(entry)),
      bound: runtime.assertType(params.bound, 'boolean'),
      optional: runtime.assertType(params.optional, 'boolean'),
    };
  },
};
