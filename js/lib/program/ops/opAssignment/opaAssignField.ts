import {
  JPLTypeError,
  applyArray,
  applyObject,
  type JPLRuntimeScope,
  type JPLSelectorParams,
} from '@/library';
import { call } from '@/program/ops/utils';
import type { UserJPLSelectorParams } from '@/program/params';
import type JPLRuntime from '@/runtime';

export default {
  /** { pipe: [op], optional: boolean } */
  async op(
    runtime: JPLRuntime,
    input: unknown,
    target: unknown,
    params: JPLSelectorParams,
    scope: JPLRuntimeScope,
    next: (output: unknown) => Promise<unknown[]> | unknown[],
  ): Promise<unknown[]> {
    let items: unknown;
    const vt = runtime.unwrapValue(target);
    const tt = runtime.type(target);
    switch (tt) {
      case 'null':
      case 'object':
      case 'array':
        items = vt;
        break;

      case 'string':
        items = [...(vt as string)];
        break;

      default:
        if (params.optional) return [undefined];
        throw new JPLTypeError(
          'cannot access fields of %s (%*<100v) (assignment)',
          tt,
          vt,
        );
    }

    const fields = await runtime.executeInstructions(
      params.pipe ?? [],
      [input],
      scope,
      (output) => [output],
    );

    const iter = async (from: number, source: unknown): Promise<unknown[]> => {
      // Call stack decoupling - This is necessary as some browsers (i.e. Safari) have very limited call stack sizes which result in stack overflow exceptions in certain situations.
      await undefined;

      scope.signal.checkHealth();

      if (from >= fields.length) return [source];

      const t = tt === 'string' ? tt : runtime.type(source);
      const field = runtime.unwrapValue(fields[from]);
      const tf = runtime.type(field);
      switch (t) {
        case 'null': {
          switch (tf) {
            case 'string':
              return runtime.muxAll([await next(null)], (output) => {
                if (output === undefined) return iter(from + 1, source);
                return iter(from + 1, { [field as string]: output });
              });

            case 'number': {
              const i = Math.trunc(field as number);
              return runtime.muxAll([await next(null)], (output) => {
                if (output === undefined) return iter(from + 1, source);
                return iter(from + 1, applyArray([], [[i, output]]));
              });
            }

            default:
          }

          break;
        }

        case 'object':
          if (tf === 'string') {
            const item = (source as Record<string, unknown>)[field as string];
            return runtime.muxAll([await next(item ?? null)], (output) => {
              if (output === undefined) return iter(from + 1, source);
              return iter(
                from + 1,
                applyObject(source, [[field as string, output]]),
              );
            });
          }
          break;

        case 'array':
          if (tf === 'number') {
            const i = Math.trunc(field as number);
            const item = (source as unknown[])[
              i >= 0 ? i : (source as unknown[]).length + i
            ];
            return runtime.muxAll([await next(item ?? null)], (output) => {
              if (output === undefined) return iter(from + 1, source);
              return iter(
                from + 1,
                applyArray(source as unknown[], [[i, output]]),
              );
            });
          }
          break;

        case 'string':
          if (tf === 'number') {
            const i = Math.trunc(field as number);
            const item = (source as string[])[
              i >= 0 ? i : (source as string[]).length + i
            ];
            return runtime.muxAll([await next(item ?? null)], (output) => {
              if (output === undefined) return iter(from + 1, source);
              const r = runtime.unwrapValue(output);
              const tr = runtime.type(r);
              switch (tr) {
                case 'null':
                case 'string':
                  return iter(
                    from + 1,
                    applyArray(
                      source as string[],
                      [[i, (output as string | null) ?? ' ']],
                      ' ',
                    ),
                  );

                default:
              }

              throw new JPLTypeError(
                'cannot assign %s (%*<100v) to string (%*<100v)',
                tr,
                r,
                (source as string[]).join(''),
              );
            });
          }
          break;

        default:
      }

      if (params.optional) return iter(from + 1, source);
      throw new JPLTypeError(
        'cannot access field of %s (%*<100v) with %s (%*<100v) (assignment)',
        t,
        t === 'string' ? (source as string[]).join('') : source,
        tf,
        field,
      );
    };

    return runtime.muxAsync([await iter(0, items)], (results) => {
      if (items === results) return target;

      const t = tt === 'string' ? tt : runtime.type(results);
      switch (t) {
        case 'null':
          return undefined;

        case 'object':
        case 'array':
          return runtime.alterValue(target, () => results);

        case 'string':
          return runtime.alterValue(target, () => results.join(''));

        default:
      }

      return undefined;
    });
  },

  /** { pipe: function, optional: boolean } */
  map(runtime: JPLRuntime, params: UserJPLSelectorParams): JPLSelectorParams {
    return {
      pipe: call(params.pipe),
      optional: runtime.assertType(params.optional, 'boolean'),
    };
  },
};
