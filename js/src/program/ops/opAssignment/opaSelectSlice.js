import { produce } from 'immer';
import { JPLTypeError } from '../../../library';
import { call } from '../utils';

export default {
  /** { from: [op], to: [op], optional: boolean } */
  async op(runtime, input, target, params, scope, next) {
    const [froms, tos] = await Promise.all([
      runtime.executeInstructions(params.from, [input], scope),
      runtime.executeInstructions(params.to, [input], scope),
    ]);

    const ranges = runtime.muxOne(
      [runtime.unwrapValues(froms), runtime.unwrapValues(tos)],
      (from, to) => [from, to],
    );

    const iter = async (from, source) => {
      // Call stack decoupling - This is necessary as some browsers (i.e. Safari) have very limited call stack sizes which result in stack overflow exceptions in certain situations.
      await undefined;

      scope.signal.checkHealth();

      if (from >= ranges.length) return [source];

      const tv = runtime.type(source);
      const [start, end] = ranges[from];
      const ts = runtime.type(start);
      const te = runtime.type(end);
      switch (tv) {
        case 'array':
          if (['number', 'null'].includes(ts) && ['number', 'null'].includes(te)) {
            const v = runtime.unwrapValue(source);
            return runtime.muxAll(
              [await next(v.slice(Math.trunc(start ?? 0), Math.trunc(end ?? v.length)))],
              async (output) =>
                iter(
                  from + 1,
                  output === undefined
                    ? source
                    : await runtime.alterValue(source, (value) => {
                        const result = runtime.unwrapValue(output);
                        const tr = runtime.type(result);
                        switch (tr) {
                          case 'null':
                          case 'array':
                            // TODO: remove immer dependency
                            return produce(value, (draft) => {
                              let s = Math.trunc(start ?? 0);
                              let e = Math.trunc(end ?? draft.length);
                              s =
                                s >= 0 ? Math.min(draft.length, s) : Math.max(0, draft.length + s);
                              e =
                                e >= 0 ? Math.min(draft.length, e) : Math.max(0, draft.length + e);
                              draft.splice(s, Math.max(0, e - s), ...(result ?? []));
                            });

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
                ),
            );
          }
          break;

        case 'string':
          if (['number', 'null'].includes(ts) && ['number', 'null'].includes(te)) {
            const chars = [...runtime.unwrapValue(source)];
            return runtime.muxAll(
              [
                await next(
                  chars.slice(Math.trunc(start ?? 0), Math.trunc(end ?? chars.length)).join(''),
                ),
              ],
              async (output) =>
                iter(
                  from + 1,
                  output === undefined
                    ? source
                    : await runtime.alterValue(source, (value) => {
                        const result = runtime.unwrapValue(output);
                        const tr = runtime.type(result);
                        switch (tr) {
                          case 'null':
                          case 'string':
                            // TODO: remove immer dependency
                            return produce([...value], (draft) => {
                              let s = Math.trunc(start ?? 0);
                              let e = Math.trunc(end ?? draft.length);
                              s =
                                s >= 0 ? Math.min(draft.length, s) : Math.max(0, draft.length + s);
                              e =
                                e >= 0 ? Math.min(draft.length, e) : Math.max(0, draft.length + e);
                              draft.splice(s, Math.max(0, e - s), ...(result ?? []));
                            }).join('');

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
                ),
            );
          }
          break;

        default:
      }

      if (params.optional) return iter(from + 1, source);
      throw new JPLTypeError(
        'cannot slice %s (%*<100v) with %s (%*<100v) and %s (%*<100v)',
        tv,
        runtime.unwrapValue(source),
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
  map(runtime, params) {
    return {
      from: call(params.from),
      to: call(params.to),
      optional: runtime.assertType(params.optional, 'boolean'),
    };
  },
};
