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

          if (ta === 'null' || tb === 'null') {
            if (ta !== 'null') return a;
            if (tb !== 'null') return b;
            return null;
          }

          switch (ta) {
            case 'number':
              if (tb === 'number') return a + b;
              break;

            case 'array':
              if (tb === 'array') {
                if (b.length === 0) return a;
                if (a.length === 0) return b;
                return [...a, ...b];
              }
              break;

            case 'string':
              if (tb === 'string') return a + b;
              break;

            case 'object':
              if (tb === 'object') return applyObject(a, Object.entries(b));
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
  map(runtime, params) {
    return {
      by: call(params.by),
    };
  },
};
