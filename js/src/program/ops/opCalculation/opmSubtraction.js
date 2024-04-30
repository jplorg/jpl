import { JPLTypeError, applyObject } from '../../../library';
import { call } from '../utils';

export default {
  /** { by: [op] } */
  op(runtime, input, target, params, scope, next) {
    return runtime.executeInstructions(params.by ?? [], [input], scope, async (by) =>
      next(
        await runtime.alterValue(target, (a) => {
          const b = runtime.unwrapValue(by);

          const ta = runtime.type(a);
          const tb = runtime.type(b);

          switch (ta) {
            case 'number':
              if (tb === 'number') return a - b;
              break;

            case 'array':
              if (tb === 'array') {
                if (a.length === 0 || b.length === 0) return a;
                const filtered = a.filter((v) => !b.some((entry) => runtime.equals(v, entry)));
                return a.length === filtered.length ? a : filtered;
              }
              break;

            case 'string':
              if (tb === 'string') return a.replaceAll(b, '');
              break;

            case 'object':
              switch (tb) {
                case 'array': {
                  if (b.length === 0) return a;
                  return applyObject(
                    a,
                    Object.entries(a)
                      .filter(([, v]) => b.some((entry) => runtime.equals(v, entry)))
                      .map(([k]) => [k]),
                  );
                }

                case 'string':
                  return applyObject(a, [[b]]);

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
  map(runtime, params) {
    return {
      by: call(params.by),
    };
  },
};
