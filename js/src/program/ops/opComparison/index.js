import {
  JPLFatalError,
  OPC_EQUAL,
  OPC_GREATER,
  OPC_GREATEREQUAL,
  OPC_LESS,
  OPC_LESSEQUAL,
  OPC_UNEQUAL,
} from '../../../library';
import { call } from '../utils';
import opcEqual from './opcEqual';
import opcGreater from './opcGreater';
import opcGreaterEqual from './opcGreaterEqual';
import opcLess from './opcLess';
import opcLessEqual from './opcLessEqual';
import opcUnequal from './opcUnequal';

export default {
  /** { pipe: [op], comparisons: [opc] } */
  op(runtime, input, params, scope, next) {
    const iter = async (from, value) => {
      // Call stack decoupling - This is necessary as some browsers (i.e. Safari) have very limited call stack sizes which result in stack overflow exceptions in certain situations.
      await undefined;

      scope.signal.checkHealth();

      if (from >= (params.comparisons?.length ?? 0)) return next(value, scope);

      const { op, params: opParams } = params.comparisons[from];
      const operator = opcs[op];
      if (!operator) throw new JPLFatalError(`invalid OPC '${op}'`);

      return operator.op(runtime, input, value, opParams ?? {}, scope, (output) =>
        iter(from + 1, output),
      );
    };

    return runtime.executeInstructions(params.pipe ?? [], [input], scope, (output) =>
      iter(0, output),
    );
  },

  /** { pipe: function, comparisons: [opc] } */
  map(runtime, params) {
    return {
      pipe: call(params.pipe),
      comparisons: runtime.muxOne([params.comparisons], ({ op, params: opParams }) => {
        const operator = opcs[op];
        if (!operator) throw new JPLFatalError(`invalid OPC '${op}'`);

        return {
          op: runtime.assertType(op, 'string'),
          params: operator.map(runtime, opParams),
        };
      }),
    };
  },
};

const opcs = {
  [OPC_EQUAL]: opcEqual,
  [OPC_GREATER]: opcGreater,
  [OPC_GREATEREQUAL]: opcGreaterEqual,
  [OPC_LESS]: opcLess,
  [OPC_LESSEQUAL]: opcLessEqual,
  [OPC_UNEQUAL]: opcUnequal,
};
