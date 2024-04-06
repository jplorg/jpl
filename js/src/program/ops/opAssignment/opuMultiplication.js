import { OPM_MULTIPLICATION, OP_CALCULATION, OP_CONSTANT } from '../../../library';
import { call } from '../utils';

export default {
  /** { pipe: [op] } */
  op(runtime, input, target, params, scope, next) {
    return runtime.executeInstructions(
      [
        {
          op: OP_CALCULATION,
          params: {
            pipe: [{ op: OP_CONSTANT, params: { value: target } }],
            operations: [{ op: OPM_MULTIPLICATION, params: { by: params.pipe } }],
          },
        },
      ],
      [input],
      scope,
      next,
    );
  },

  /** { value: function } */
  map(runtime, params) {
    return {
      pipe: call(params.value),
    };
  },
};
