import { JPLTypeError, type JPLRuntimeScope } from '@/library';
import type { JPLSelectorParams } from '@/library/definition';
import type { UserJPLSelectorParams } from '@/program/params';
import type JPLRuntime from '@/runtime';
import { call } from '../utils';

export default {
  /** { pipe: [op], optional: boolean } */
  op(
    runtime: JPLRuntime,
    input: unknown,
    target: unknown,
    params: JPLSelectorParams,
    scope: JPLRuntimeScope,
    next: (output: unknown) => Promise<unknown[]> | unknown[],
  ): Promise<unknown[]> {
    const value = runtime.unwrapValue(target);
    const tv = runtime.type(value);

    return runtime.executeInstructions(
      params.pipe ?? [],
      [input],
      scope,
      (output) => {
        const field = runtime.unwrapValue(output);
        const tf = runtime.type(field);
        switch (tv) {
          case 'null':
            if (['string', 'number'].includes(tf)) return next(null);
            break;

          case 'object':
            if (tf === 'string')
              return next(
                (value as Record<string, unknown>)[field as string] ?? null,
              );
            break;

          case 'array':
            if (tf === 'number') {
              const i = Math.trunc(field as number);
              return next(
                (value as unknown[])[
                  i >= 0 ? i : (value as unknown[]).length + i
                ] ?? null,
              );
            }
            break;

          case 'string':
            if (tf === 'number') {
              const i = Math.trunc(field as number);
              const chars = [...(value as string)];
              return next(chars[i >= 0 ? i : chars.length + i] ?? null);
            }
            break;

          default:
        }

        if (params.optional) return [];
        throw new JPLTypeError(
          'cannot access field of %s (%*<100v) with %s (%*<100v)',
          tv,
          value,
          tf,
          field,
        );
      },
    );
  },

  /** { pipe: function, optional: boolean } */
  map(runtime: JPLRuntime, params: UserJPLSelectorParams): JPLSelectorParams {
    return {
      pipe: call(params.pipe),
      optional: runtime.assertType(params.optional, 'boolean'),
    };
  },
};
