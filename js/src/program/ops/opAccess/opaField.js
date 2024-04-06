import { JPLTypeError } from '../../../library';
import { call } from '../utils';

export default {
  /** { pipe: [op], optional: boolean } */
  op(runtime, input, target, params, scope, next) {
    const value = runtime.unwrapValue(target);
    const tv = runtime.type(value);

    return runtime.executeInstructions(params.pipe, [input], scope, (output) => {
      const field = runtime.unwrapValue(output);
      const tf = runtime.type(field);
      switch (tv) {
        case 'null':
          return next(null);

        case 'object':
          if (tf === 'string') return next(value[field] ?? null);
          break;

        case 'array':
          if (tf === 'number') {
            const i = Math.trunc(field);
            return next(value[i >= 0 ? i : value.length + i] ?? null);
          }
          break;

        case 'string':
          if (tf === 'number') {
            const i = Math.trunc(field);
            const chars = [...value];
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
    });
  },

  /** { value: function, optional: boolean } */
  map(runtime, params) {
    return {
      pipe: call(params.value),
      optional: runtime.assertType(params.optional, 'boolean'),
    };
  },
};
