import { call } from './utils';

export default {
  /** { pipe: [op] } */
  async op(runtime, input, params, scope, next) {
    const outputs = await runtime.executeInstructions(params.pipe, [input], scope);
    return next(outputs, scope);
  },

  /** { pipe: function } */
  map(runtime, params) {
    return {
      pipe: call(params.pipe),
    };
  },
};
