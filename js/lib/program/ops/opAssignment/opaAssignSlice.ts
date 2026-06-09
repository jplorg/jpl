import { JPLTypeError, type JPLRuntimeScope } from '@/library';
import type { JPLSelectorParams } from '@/library/definition';
import type { UserJPLSelectorParams } from '@/program/params';
import type JPLRuntime from '@/runtime';
import { call } from '../utils';

function shallowCompareArrays(a: unknown[], b: unknown[]): boolean {
  if (a.length !== b.length) return false;
  return a.every((v, i) => v === b[i]);
}

export default {
  /** { from: [op], to: [op], optional: boolean } */
  async op(
    runtime: JPLRuntime,
    input: unknown,
    target: unknown,
    params: JPLSelectorParams,
    scope: JPLRuntimeScope,
    next: (output: unknown) => Promise<unknown[]> | unknown[],
  ): Promise<unknown[]> {
    const [froms, tos] = await Promise.all([
      runtime.executeInstructions(params.from ?? [], [input], scope),
      runtime.executeInstructions(params.to ?? [], [input], scope),
    ]);

    const ranges = runtime.muxOne(
      [runtime.unwrapValues(froms), runtime.unwrapValues(tos)],
      (from, to) => [from, to],
    );

    const iter = async (from: number, source: unknown): Promise<unknown[]> => {
      // Call stack decoupling - This is necessary as some browsers (i.e. Safari) have very limited call stack sizes which result in stack overflow exceptions in certain situations.
      await undefined;

      scope.signal.checkHealth();

      if (from >= ranges.length) return [source];

      const v = runtime.unwrapValue(source);
      const tv = runtime.type(source);
      const [start, end] = ranges[from];
      const ts = runtime.type(start);
      const te = runtime.type(end);
      switch (tv) {
        case 'null':
          return iter(from + 1, source);

        case 'array':
          if (
            ['number', 'null'].includes(ts) &&
            ['number', 'null'].includes(te)
          ) {
            return runtime.muxAll(
              [
                await next(
                  (v as unknown[]).slice(
                    Math.trunc(start ?? 0),
                    Math.trunc(end ?? (v as unknown[]).length),
                  ),
                ),
              ],
              async (output) => {
                if (output === undefined) return iter(from + 1, source);
                return iter(
                  from + 1,
                  await runtime.alterValue(source, (value) => {
                    const result = runtime.unwrapValue(output);
                    const tr = runtime.type(result);
                    switch (tr) {
                      case 'null':
                      case 'array': {
                        let s = Math.trunc(start ?? 0);
                        let e = Math.trunc(end ?? (value as unknown[]).length);
                        s =
                          s >= 0
                            ? Math.min((value as unknown[]).length, s)
                            : Math.max(0, (value as unknown[]).length + s);
                        e = Math.max(
                          s,
                          e >= 0
                            ? Math.min((value as unknown[]).length, e)
                            : Math.max(0, (value as unknown[]).length + e),
                        );
                        if (
                          shallowCompareArrays(
                            (value as unknown[]).slice(s, e),
                            (result ?? []) as unknown[],
                          )
                        )
                          return value as unknown[];
                        const c = [...(value as unknown[])];
                        c.splice(s, e - s, ...((result ?? []) as unknown[]));
                        return c;
                      }

                      default:
                    }

                    throw new JPLTypeError(
                      'cannot assign %s (%*<100v) to slice of %s (%*<100v)',
                      tr,
                      result,
                      tv,
                      value,
                    );
                  }),
                );
              },
            );
          }
          break;

        case 'string':
          if (
            ['number', 'null'].includes(ts) &&
            ['number', 'null'].includes(te)
          ) {
            const chars = [...(v as string)];
            return runtime.muxAll(
              [
                await next(
                  chars
                    .slice(
                      Math.trunc(start ?? 0),
                      Math.trunc(end ?? chars.length),
                    )
                    .join(''),
                ),
              ],
              async (output) => {
                if (output === undefined) return iter(from + 1, source);
                return iter(
                  from + 1,
                  await runtime.alterValue(source, (value) => {
                    const result = runtime.unwrapValue(output);
                    const tr = runtime.type(result);
                    switch (tr) {
                      case 'null':
                      case 'string': {
                        const chars = [...(value as string)];
                        let s = Math.trunc(start ?? 0);
                        let e = Math.trunc(end ?? chars.length);
                        s =
                          s >= 0
                            ? Math.min(chars.length, s)
                            : Math.max(0, chars.length + s);
                        e = Math.max(
                          s,
                          e >= 0
                            ? Math.min(chars.length, e)
                            : Math.max(0, chars.length + e),
                        );
                        if (
                          shallowCompareArrays(
                            chars.slice(s, e),
                            (result ?? []) as unknown[],
                          )
                        )
                          return value;
                        chars.splice(s, e - s, ...((result ?? []) as string));
                        return chars.join('');
                      }

                      default:
                    }

                    throw new JPLTypeError(
                      'cannot assign %s (%*<100v) to slice of %s (%*<100v)',
                      tr,
                      result,
                      tv,
                      value,
                    );
                  }),
                );
              },
            );
          }
          break;

        default:
      }

      if (params.optional) return iter(from + 1, source);
      throw new JPLTypeError(
        'cannot slice %s (%*<100v) with %s (%*<100v) and %s (%*<100v) (assignment)',
        tv,
        v,
        ts,
        start,
        te,
        end,
      );
    };

    return runtime.muxOne([await iter(0, target)], (output) =>
      runtime.type(output) === 'null' ? undefined : output,
    );
  },

  /** { from: function, to: function, optional: boolean } */
  map(runtime: JPLRuntime, params: UserJPLSelectorParams): JPLSelectorParams {
    return {
      from: call(params.from),
      to: call(params.to),
      optional: runtime.assertType(params.optional, 'boolean'),
    };
  },
};
