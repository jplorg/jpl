import { call } from '../utils';

export default {
  /** { by: [op] } */
  op(runtime, input, target, params, scope, next) {
    return runtime.executeInstructions(params.by ?? [], [input], scope, (by) =>
      next(!runtime.equals(target, by)),
    );
  },

  /** { by: function } */
  map(runtime, params) {
    return {
      by: call(params.by),
    };
  },
};
