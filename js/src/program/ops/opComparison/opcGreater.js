import { call } from '../utils';

export default {
  /** { by: [op] } */
  op(runtime, input, target, params, scope, next) {
    return runtime.executeInstructions(params.by ?? [], [input], scope, (by) =>
      next(runtime.compare(target, by) > 0),
    );
  },

  /** { by: function } */
  map(runtime, params) {
    return {
      by: call(params.by),
    };
  },
};
