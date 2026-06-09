import type { JPLInstructionParams } from '@/library/definition';
import type JPLRuntime from '@/runtime';
import { JPLReferenceError, type JPLRuntimeScope } from '../../library';
import type { UserJPLInstructionParams } from '../params';

export default {
  /** { name: string } */
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
    if (!Object.hasOwn(scope.vars, params.name ?? ''))
      throw new JPLReferenceError('%s is not defined', params.name ?? '');
    return next(scope.vars[params.name ?? ''], scope);
  },

  /** { name: string } */
  map(
    runtime: JPLRuntime,
    params: UserJPLInstructionParams,
  ): JPLInstructionParams {
    return {
      name: runtime.assertType(params.name, 'string'),
    };
  },
};
