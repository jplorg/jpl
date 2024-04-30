import { JPLErrorEnclosure, JPLExecutionError } from '../../library';
import { call } from './utils';

export default {
  /** { try: [op], catch: [op] } */
  async op(runtime, input, params, scope, next) {
    const nextScope = scope.next({ signal: scope.signal.next() });
    try {
      return await runtime.executeInstructions(
        params.try ?? [],
        [input],
        nextScope,
        async (output) => {
          try {
            return await next(output, scope);
          } catch (err) {
            throw new JPLErrorEnclosure(err);
          }
        },
      );
    } catch (err) {
      if (JPLErrorEnclosure.is(err)) throw err.inner;
      if (!JPLExecutionError.is(err)) throw err;
      nextScope.signal.exit();
      return runtime.executeInstructions(params.catch ?? [], [err.value], scope, (output) =>
        next(output, scope),
      );
    }
  },

  /** { try: function, catch: function } */
  map(runtime, params) {
    return {
      try: call(params.try),
      catch: call(params.catch),
    };
  },
};
