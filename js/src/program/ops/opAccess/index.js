import { JPLFatalError, OPA_FIELD, OPA_FUNCTION, OPA_ITER, OPA_SLICE } from '../../../library';
import { call } from '../utils';
import opaField from './opaField';
import opaFunction from './opaFunction';
import opaIter from './opaIter';
import opaSlice from './opaSlice';

export default {
  /** { pipe: [op], operations: [opa] } */
  op(runtime, input, params, scope, next) {
    const iter = async (from, value) => {
      // Call stack decoupling - This is necessary as some browsers (i.e. Safari) have very limited call stack sizes which result in stack overflow exceptions in certain situations.
      await undefined;

      scope.signal.checkHealth();

      if (from >= params.operations.length) return next(value, scope);

      const { op, params: opParams } = params.operations[from];
      const operator = ops[op];
      if (!operator) throw new JPLFatalError(`invalid OPA '${op}'`);

      return operator.op(runtime, input, value, opParams, scope, (output) =>
        iter(from + 1, output),
      );
    };

    return runtime.executeInstructions(params.pipe, [input], scope, (output) => iter(0, output));
  },

  /** { value: function, operations: [opa] } */
  map(runtime, params) {
    return {
      pipe: call(params.value),
      operations: runtime.muxOne([params.operations], ({ op, params: opParams }) => {
        const operator = ops[op];
        if (!operator) throw new JPLFatalError(`invalid OPA '${op}'`);

        return {
          op: runtime.assertType(op, 'string'),
          params: operator.map(runtime, opParams),
        };
      }),
    };
  },
};

const ops = {
  [OPA_FIELD]: opaField,
  [OPA_FUNCTION]: opaFunction,
  [OPA_ITER]: opaIter,
  [OPA_SLICE]: opaSlice,
};
