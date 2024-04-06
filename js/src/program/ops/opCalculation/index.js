import {
  JPLFatalError,
  OPM_ADDITION,
  OPM_DIVISION,
  OPM_MULTIPLICATION,
  OPM_REMAINDER,
  OPM_SUBTRACTION,
} from '../../../library';
import { call } from '../utils';
import opmAddition from './opmAddition';
import opmDivision from './opmDivision';
import opmMultiplication from './opmMultiplication';
import opmRemainder from './opmRemainder';
import opmSubtraction from './opmSubtraction';

export default {
  /** { pipe: [op], operations: [opm] } */
  op(runtime, input, params, scope, next) {
    const iter = async (from, value) => {
      // Call stack decoupling - This is necessary as some browsers (i.e. Safari) have very limited call stack sizes which result in stack overflow exceptions in certain situations.
      await undefined;

      scope.signal.checkHealth();

      if (from >= params.operations.length) return next(value, scope);

      const { op, params: opParams } = params.operations[from];
      const operator = ops[op];
      if (!operator) throw new JPLFatalError(`invalid OPM '${op}'`);

      return operator.op(runtime, input, value, opParams, scope, (output) =>
        iter(from + 1, output),
      );
    };

    return runtime.executeInstructions(params.pipe, [input], scope, (output) => iter(0, output));
  },

  /** { value: function, operations: [opm] } */
  map(runtime, params) {
    return {
      pipe: call(params.value),
      operations: runtime.muxOne([params.operations], ({ op, params: opParams }) => {
        const operator = ops[op];
        if (!operator) throw new JPLFatalError(`invalid OPM '${op}'`);

        return {
          op: runtime.assertType(op, 'string'),
          params: operator.map(runtime, opParams),
        };
      }),
    };
  },
};

const ops = {
  [OPM_ADDITION]: opmAddition,
  [OPM_DIVISION]: opmDivision,
  [OPM_MULTIPLICATION]: opmMultiplication,
  [OPM_REMAINDER]: opmRemainder,
  [OPM_SUBTRACTION]: opmSubtraction,
};
