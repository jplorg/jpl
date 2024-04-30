import { JPLTypeError } from '../../../library';
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
              if (tb === 'number') return a * b;
              break;

            case 'string':
              if (tb === 'number') {
                if (b < 1) return null;
                return a.repeat(b);
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
  map(runtime, params) {
    return {
      by: call(params.by),
    };
  },
};
