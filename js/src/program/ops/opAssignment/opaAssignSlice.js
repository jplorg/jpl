import { JPLTypeError } from '../../../library';
import { call } from '../utils';

function shallowCompareArrays(a, b) {
  if (a.length !== b.length) return false;
  return a.every((v, i) => v === b[i]);
}

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

      const v = runtime.unwrapValue(source);
      const tv = runtime.type(source);
      const [start, end] = ranges[from];
      const ts = runtime.type(start);
      const te = runtime.type(end);
      switch (tv) {
        case 'array':
          if (['number', 'null'].includes(ts) && ['number', 'null'].includes(te)) {
            return runtime.muxAll(
              [await next(v.slice(Math.trunc(start ?? 0), Math.trunc(end ?? v.length)))],
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
                        let e = Math.trunc(end ?? value.length);
                        s = s >= 0 ? Math.min(value.length, s) : Math.max(0, value.length + s);
                        e = Math.max(
                          s,
                          e >= 0 ? Math.min(value.length, e) : Math.max(0, value.length + e),
                        );
                        if (shallowCompareArrays(value.slice(s, e), result ?? [])) return value;
                        const c = [...value];
                        c.splice(s, e - s, ...(result ?? []));
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
          if (['number', 'null'].includes(ts) && ['number', 'null'].includes(te)) {
            const chars = [...v];
            return runtime.muxAll(
              [
                await next(
                  chars.slice(Math.trunc(start ?? 0), Math.trunc(end ?? chars.length)).join(''),
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
                        const chars = [...value];
                        let s = Math.trunc(start ?? 0);
                        let e = Math.trunc(end ?? chars.length);
                        s = s >= 0 ? Math.min(chars.length, s) : Math.max(0, chars.length + s);
                        e = Math.max(
                          s,
                          e >= 0 ? Math.min(chars.length, e) : Math.max(0, chars.length + e),
                        );
                        if (shallowCompareArrays(chars.slice(s, e), result ?? [])) return value;
                        chars.splice(s, e - s, ...(result ?? []));
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
  map(runtime, params) {
    return {
      from: call(params.from),
      to: call(params.to),
      optional: runtime.assertType(params.optional, 'boolean'),
    };
  },
};
