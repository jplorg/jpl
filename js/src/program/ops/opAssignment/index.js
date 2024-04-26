import {
  JPLFatalError,
  OPA_FIELD,
  OPA_ITER,
  OPA_SLICE,
  OPU_ADDITION,
  OPU_DIVISION,
  OPU_MULTIPLICATION,
  OPU_NULL_COALESCENCE,
  OPU_REMAINDER,
  OPU_SET,
  OPU_SUBTRACTION,
  OPU_UPDATE,
} from '../../../library';
import { call } from '../utils';
import opaSelectField from './opaSelectField';
import opaSelectIter from './opaSelectIter';
import opaSelectSlice from './opaSelectSlice';
import opuAddition from './opuAddition';
import opuDivision from './opuDivision';
import opuMultiplication from './opuMultiplication';
import opuNullCoalescence from './opuNullCoalescence';
import opuRemainder from './opuRemainder';
import opuSet from './opuSet';
import opuSubtraction from './opuSubtraction';
import opuUpdate from './opuUpdate';

export default {
  /** { pipe: [op], selectors: [opa], assignment: [opu] } */
  op(runtime, input, params, scope, next) {
    const iter = async (from, value) => {
      // Call stack decoupling - This is necessary as some browsers (i.e. Safari) have very limited call stack sizes which result in stack overflow exceptions in certain situations.
      await undefined;

      scope.signal.checkHealth();

      if (from >= params.selectors.length) {
        const { op, params: opParams } = params.assignment;
        const operator = opus[op];
        if (!operator) throw new JPLFatalError(`invalid OPU '${op}'`);

        return operator.op(runtime, input, value, opParams, scope, (output) => [output]);
      }

      const { op, params: opParams } = params.selectors[from];
      const operator = opasSelect[op];
      if (!operator) throw new JPLFatalError(`invalid OPA '${op}' (assignment)`);

      return operator.op(runtime, input, value, opParams, scope, (output) =>
        iter(from + 1, output),
      );
    };

    return runtime.executeInstructions(params.pipe, [input], scope, async (output) =>
      runtime.muxAll([await iter(0, output)], (result) =>
        next(result === undefined ? output : result, scope),
      ),
    );
  },

  /** { pipe: function, selectors: [opa], assignment: opu } */
  map(runtime, params) {
    return {
      pipe: call(params.pipe),
      selectors: runtime.muxOne([params.selectors], ({ op, params: opParams }) => {
        const operator = opasSelect[op];
        if (!operator) throw new JPLFatalError(`invalid OPA '${op}' (assignment)`);

        return {
          op: runtime.assertType(op, 'string'),
          params: operator.map(runtime, opParams),
        };
      }),
      assignment: (({ op, params: opParams }) => {
        const operator = opus[op];
        if (!operator) throw new JPLFatalError(`invalid OPU '${op}'`);

        return {
          op: runtime.assertType(op, 'string'),
          params: operator.map(runtime, opParams),
        };
      })(params.assignment),
    };
  },
};

const opasSelect = {
  [OPA_FIELD]: opaSelectField,
  [OPA_ITER]: opaSelectIter,
  [OPA_SLICE]: opaSelectSlice,
};

const opus = {
  [OPU_ADDITION]: opuAddition,
  [OPU_DIVISION]: opuDivision,
  [OPU_MULTIPLICATION]: opuMultiplication,
  [OPU_NULL_COALESCENCE]: opuNullCoalescence,
  [OPU_REMAINDER]: opuRemainder,
  [OPU_SET]: opuSet,
  [OPU_SUBTRACTION]: opuSubtraction,
  [OPU_UPDATE]: opuUpdate,
};
