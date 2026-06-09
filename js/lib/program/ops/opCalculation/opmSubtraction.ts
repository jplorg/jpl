import {
  JPLTypeError,
  applyObject,
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
                if (tb === 'number') return (a as number) - (b as number);
                break;

              case 'array':
                if (tb === 'array') {
                  if (
                    (a as unknown[]).length === 0 ||
                    (b as unknown[]).length === 0
                  )
                    return a;
                  const filtered = (a as unknown[]).filter(
                    (v) =>
                      !(b as unknown[]).some((entry) =>
                        runtime.equals(v, entry),
                      ),
                  );
                  return (a as unknown[]).length === filtered.length
                    ? a
                    : filtered;
                }
                break;

              case 'string':
                if (tb === 'string')
                  return (a as string).replaceAll(b as string, '');
                break;

              case 'object':
                switch (tb) {
                  case 'array': {
                    if ((b as unknown[]).length === 0) return a;
                    return applyObject(
                      a,
                      Object.entries(a as Record<string, unknown>)
                        .filter(([, v]) =>
                          (b as unknown[]).some((entry) =>
                            runtime.equals(v, entry),
                          ),
                        )
                        .map(([k]) => [k]),
                    );
                  }

                  case 'string':
                    return applyObject(a, [[b as string]]);

                  default:
                }
                break;

              default:
            }

            throw new JPLTypeError(
              '%s (%*<100v) and %s (%*<100v) cannot be subtracted',
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
