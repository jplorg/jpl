import {
  JPLFatalError,
  OPA_FIELD,
  OPA_FUNCTION,
  OPA_ITER,
  OPA_SLICE,
  type JPLInstructionParams,
  type JPLRuntimeScope,
  type JPLSelectorParams,
} from '@/library';
import type { JPLOPSubHandler } from '@/program/handler';
import { call } from '@/program/ops/utils';
import type {
  UserJPLInstructionParams,
  UserJPLSelectorParams,
} from '@/program/params';
import type JPLRuntime from '@/runtime';
import opaField from './opaField';
import opaFunction from './opaFunction';
import opaIter from './opaIter';
import opaSlice from './opaSlice';

export default {
  /** { pipe: [op], selectors: [opa] } */
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

      if (from >= (params.selectors?.length ?? 0)) return next(value, scope);

      const { op, params: opParams } = params.selectors![from];
      const operator = opas[op];
      if (!operator) throw new JPLFatalError(`invalid OPA '${op}'`);

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

  /** { pipe: function, selectors: [opa] } */
  map(
    runtime: JPLRuntime,
    params: UserJPLInstructionParams,
  ): JPLInstructionParams {
    return {
      pipe: call(params.pipe),
      selectors: runtime.muxOne(
        [params.selectors!],
        ({ op, params: opParams }) => {
          const operator = opas[op];
          if (!operator) throw new JPLFatalError(`invalid OPA '${op}'`);

          return {
            op: runtime.assertType(op, 'string'),
            params: operator.map(runtime, opParams),
          };
        },
      ),
    };
  },
};

const opas: {
  [op: string]: JPLOPSubHandler<JPLSelectorParams, UserJPLSelectorParams>;
} = {
  [OPA_FIELD]: opaField,
  [OPA_FUNCTION]: opaFunction,
  [OPA_ITER]: opaIter,
  [OPA_SLICE]: opaSlice,
};
