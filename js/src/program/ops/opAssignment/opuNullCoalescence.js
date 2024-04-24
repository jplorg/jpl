import { OP_CONSTANT, OP_NULL_COALESCENCE } from '../../../library';
import { call } from '../utils';

export default {
  /** { pipe: [op] } */
  op(runtime, input, target, params, scope, next) {
    return runtime.executeInstructions(
      [
        {
          op: OP_NULL_COALESCENCE,
          params: { pipes: [[{ op: OP_CONSTANT, params: { value: target } }], params.pipe] },
        },
      ],
      [input],
      scope,
      next,
    );
  },

  /** { pipe: function } */
  map(runtime, params) {
    return {
      pipe: call(params.pipe),
    };
  },
};
