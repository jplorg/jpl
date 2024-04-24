import { call } from './utils';

export default {
  /** { name: string, pipe: [op] } */
  op(runtime, input, params, scope, next) {
    return runtime.executeInstructions(params.pipe, [input], scope, (output) =>
      next(input, scope.next({ vars: { [params.name]: output } })),
    );
  },

  /** { name: string, pipe: function } */
  map(runtime, params) {
    return {
      name: runtime.assertType(params.name, 'string'),
      pipe: call(params.pipe),
    };
  },
};
