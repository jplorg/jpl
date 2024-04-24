import { call } from '../utils';

export default {
  /** { pipe: [op] } */
  op(runtime, input, target, params, scope, next) {
    return runtime.executeInstructions(params.pipe, [target], scope, (output) => next(output));
  },

  /** { pipe: function } */
  map(runtime, params) {
    return {
      pipe: call(params.pipe),
    };
  },
};
