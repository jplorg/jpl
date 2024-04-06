import { call } from '../utils';

export default {
  /** { pipe: [op] } */
  op(runtime, input, target, params, scope, next) {
    return runtime.executeInstructions(params.pipe, [input], scope, next);
  },

  /** { value: function } */
  map(runtime, params) {
    return {
      pipe: call(params.value),
    };
  },
};