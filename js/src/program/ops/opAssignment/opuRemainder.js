import { OPM_REMAINDER, OP_CALCULATION } from '../../../library';
import { call, constant } from '../utils';

export default {
  /** { pipe: [op] } */
  op(runtime, input, target, params, scope, next) {
    return runtime.executeInstructions(
      [
        {
          op: OP_CALCULATION,
          params: {
            pipe: constant(target),
            operations: [{ op: OPM_REMAINDER, params: { by: params.pipe ?? [] } }],
          },
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
