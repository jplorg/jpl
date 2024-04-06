import { call } from './utils';

export default {
  /** { pipes: [[op]] } */
  op(runtime, input, params, scope, next) {
    const iter = async (from) => {
      // Call stack decoupling - This is necessary as some browsers (i.e. Safari) have very limited call stack sizes which result in stack overflow exceptions in certain situations.
      await undefined;

      scope.signal.checkHealth();

      if (from >= params.pipes.length) return next(true, scope);

      const pipe = params.pipes[from];

      return runtime.executeInstructions(pipe, [input], scope, (output) => {
        if (!runtime.truthy(output)) return next(false, scope);

        return iter(from + 1);
      });
    };

    return iter(0);
  },

  /** { conditions: [function] } */
  map(runtime, params) {
    return {
      pipes: runtime.muxOne([params.conditions], (condition) => call(condition)),
    };
  },
};
