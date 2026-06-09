import { JPLTypeError, applyObject, type JPLRuntimeScope } from '@/library';
import type { JPLOperationParams } from '@/library/definition';
import type { UserJPLOperationParams } from '@/program/params';
import type JPLRuntime from '@/runtime';
import { call } from '../utils';

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

            if (ta === 'null' || tb === 'null') {
              if (ta !== 'null') return a;
              if (tb !== 'null') return b;
              return null;
            }

            switch (ta) {
              case 'number':
                if (tb === 'number') return (a as number) + (b as number);
                break;

              case 'array':
                if (tb === 'array') {
                  if ((b as unknown[]).length === 0) return a;
                  if ((a as unknown[]).length === 0) return b;
                  return [...(a as unknown[]), ...(b as unknown[])];
                }
                break;

              case 'string':
                if (tb === 'string') return (a as string) + (b as string);
                break;

              case 'object':
                if (tb === 'object')
                  return applyObject(
                    a,
                    Object.entries(b as Record<string, unknown>),
                  );
                break;

              default:
            }

            throw new JPLTypeError(
              '%s (%*<100v) and %s (%*<100v) cannot be added together',
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
