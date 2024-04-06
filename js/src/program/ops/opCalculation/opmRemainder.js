import { JPLTypeError, JPLZeroDivisionError } from '../../../library';
import { call } from '../utils';

export default {
  /** { by: [op] } */
  op(runtime, input, target, params, scope, next) {
    return runtime.executeInstructions(params.by, [input], scope, async (by) =>
      next(
        await runtime.alterValue(target, (a) => {
          const b = runtime.unwrapValue(by);

          const ta = runtime.type(a);
          const tb = runtime.type(b);

          switch (ta) {
            case 'number':
              if (tb === 'number') {
                if (b === 0)
                  throw new JPLZeroDivisionError(
                    '%s (%*<100v) cannot be divided by zero (remainder)',
                    ta,
                    a,
                  );
                return a % b;
              }
              break;

            default:
          }

          throw new JPLTypeError(
            '%s (%*<100v) cannot be divided by %s (%*<100v) (remainder)',
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
