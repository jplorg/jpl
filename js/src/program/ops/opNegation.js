import { JPLTypeError } from '../../library';

export default {
  /** {} */
  async op(runtime, input, params, scope, next) {
    return next(
      await runtime.alterValue(input, (value) => {
        const t = runtime.type(value);
        switch (t) {
          case 'number':
            return -value;

          default:
        }

        throw new JPLTypeError('%s (%*<100v) cannot be negated', t, value);
      }),
      scope,
    );
  },

  /** {} */
  map() {
    return {};
  },
};
