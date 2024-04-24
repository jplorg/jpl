import { call } from './utils';

export default {
  /** { pipes: [[op]] } */
  op(runtime, input, params, scope, next) {
    return runtime.muxAll([params.pipes], (pipe) =>
      runtime.executeInstructions(pipe, [input], scope, (output) => next(output, scope)),
    );
  },

  /** { pipes: [function] } */
  map(runtime, params) {
    return {
      pipes: runtime.muxOne([params.pipes], (value) => call(value)),
    };
  },
};
