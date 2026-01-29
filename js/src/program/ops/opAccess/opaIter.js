import { JPLTypeError } from '../../../library';

export default {
  /** { optional: boolean } */
  op(runtime, input, target, params, scope, next) {
    const value = runtime.unwrapValue(target);
    const t = runtime.type(value);
    switch (t) {
      case 'null':
        return [];

      case 'object':
        return runtime.muxAll([Object.values(value)], next);

      case 'array':
        return runtime.muxAll([value], next);

      case 'string':
        return runtime.muxAll([[...value]], next);

      default:
    }

    if (params.optional) return [];
    throw new JPLTypeError('cannot iterate over %s (%*<100v)', t, value);
  },

  /** { optional: boolean } */
  map(runtime, params) {
    return {
      optional: runtime.assertType(params.optional, 'boolean'),
    };
  },
};
