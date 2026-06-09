import {
  type JPLInstructionParams,
  type JPLRuntimeScope,
  scopedFunction,
} from '@/library';
import type { UserJPLInstructionParams } from '@/program/params';
import type JPLRuntime from '@/runtime';
import { call } from './utils';

export default {
  /** { argNames: [string], pipe: [op] } */
  op(
    _runtime: JPLRuntime,
    _input: unknown,
    params: JPLInstructionParams,
    scope: JPLRuntimeScope,
    next: (
      output: unknown,
      scope: JPLRuntimeScope,
    ) => Promise<unknown[]> | unknown[],
  ): Promise<unknown[]> | unknown[] {
    return next(
      scopedFunction(params.argNames ?? [], params.pipe ?? [], scope),
      scope,
    );
  },

  /** { argNames: [string], pipe: function } */
  map(
    runtime: JPLRuntime,
    params: UserJPLInstructionParams,
  ): JPLInstructionParams {
    return {
      argNames: runtime.muxOne([params.argNames!], (entry) =>
        runtime.assertType(entry, 'string'),
      ),
      pipe: call(params.pipe),
    };
  },
};
