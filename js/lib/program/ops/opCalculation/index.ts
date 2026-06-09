import {
  JPLFatalError,
  OPM_ADDITION,
  OPM_DIVISION,
  OPM_MULTIPLICATION,
  OPM_REMAINDER,
  OPM_SUBTRACTION,
  type JPLInstructionParams,
  type JPLOperationParams,
  type JPLRuntimeScope,
} from '@/library';
import type { JPLOPSubHandler } from '@/program/handler';
import { call } from '@/program/ops/utils';
import type {
  UserJPLInstructionParams,
  UserJPLOperationParams,
} from '@/program/params';
import type JPLRuntime from '@/runtime';
import opmAddition from './opmAddition';
import opmDivision from './opmDivision';
import opmMultiplication from './opmMultiplication';
import opmRemainder from './opmRemainder';
import opmSubtraction from './opmSubtraction';

export default {
  /** { pipe: [op], operations: [opm] } */
  op(
    runtime: JPLRuntime,
    input: unknown,
    params: JPLInstructionParams,
    scope: JPLRuntimeScope,
    next: (
      output: unknown,
      scope: JPLRuntimeScope,
    ) => Promise<unknown[]> | unknown[],
  ): Promise<unknown[]> {
    const iter = async (from: number, value: unknown): Promise<unknown[]> => {
      // Call stack decoupling - This is necessary as some browsers (i.e. Safari) have very limited call stack sizes which result in stack overflow exceptions in certain situations.
      await undefined;

      scope.signal.checkHealth();

      if (from >= (params.operations?.length ?? 0)) return next(value, scope);

      const { op, params: opParams } = params.operations![from];
      const operator = opms[op];
      if (!operator) throw new JPLFatalError(`invalid OPM '${op}'`);

      return operator.op(
        runtime,
        input,
        value,
        opParams ?? {},
        scope,
        (output: unknown) => iter(from + 1, output),
      );
    };

    return runtime.executeInstructions(
      params.pipe ?? [],
      [input],
      scope,
      (output) => iter(0, output),
    );
  },

  /** { pipe: function, operations: [opm] } */
  map(
    runtime: JPLRuntime,
    params: UserJPLInstructionParams,
  ): JPLInstructionParams {
    return {
      pipe: call(params.pipe),
      operations: runtime.muxOne(
        [params.operations!],
        ({ op, params: opParams }) => {
          const operator = opms[op];
          if (!operator) throw new JPLFatalError(`invalid OPM '${op}'`);

          return {
            op: runtime.assertType(op, 'string'),
            params: operator.map(runtime, opParams),
          };
        },
      ),
    };
  },
};

const opms: {
  [op: string]: JPLOPSubHandler<JPLOperationParams, UserJPLOperationParams>;
} = {
  [OPM_ADDITION]: opmAddition,
  [OPM_DIVISION]: opmDivision,
  [OPM_MULTIPLICATION]: opmMultiplication,
  [OPM_REMAINDER]: opmRemainder,
  [OPM_SUBTRACTION]: opmSubtraction,
};
