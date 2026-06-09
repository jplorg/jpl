import {
  JPLTypeError,
  type JPLOperationParams,
  type JPLRuntimeScope,
} from '@/library';
import { call } from '@/program/ops/utils';
import type { UserJPLOperationParams } from '@/program/params';
import type JPLRuntime from '@/runtime';

export default {
  /** { by: [op] } */
  op(
    runtime: JPLRuntime,
    input: unknown,
    target: unknown,
    params: JPLOperationParams,
    scope: JPLRuntimeScope,
    next: (output: unknown) => Promise<unknown[]> | unknown[],
  ): Promise<unknown[]> {
    return runtime.executeInstructions(
      params.by ?? [],
      [input],
      scope,
      async (by) =>
        next(
          await runtime.alterValue(target, (a) => {
            const b = runtime.unwrapValue(by);

            const ta = runtime.type(a);
            const tb = runtime.type(b);

            switch (ta) {
              case 'number':
                if (tb === 'number') return (a as number) * (b as number);
                break;

              case 'string':
                if (tb === 'number') {
                  if ((b as number) < 1) return null;
                  return (a as string).repeat(b as number);
                }
                break;

              case 'object':
                if (tb === 'object') return runtime.merge(a, b);
                break;

              default:
            }

            throw new JPLTypeError(
              '%s (%*<100v) and %s (%*<100v) cannot be multiplied together',
              ta,
              a,
              tb,
              b,
            );
          }),
        ),
    );
  },

  /** { by: function } */
  map(
    _runtime: JPLRuntime,
    params: UserJPLOperationParams,
  ): JPLOperationParams {
    return {
      by: call(params.by),
    };
  },
};
