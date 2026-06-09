import {
  JPLTypeError,
  type JPLRuntimeScope,
  type JPLSelectorParams,
} from '@/library';
import type { UserJPLSelectorParams } from '@/program/params';
import type JPLRuntime from '@/runtime';

export default {
  /** { optional: boolean } */
  op(
    runtime: JPLRuntime,
    _input: unknown,
    target: unknown,
    params: JPLSelectorParams,
    _scope: JPLRuntimeScope,
    next: (output: unknown) => Promise<unknown[]> | unknown[],
  ): Promise<unknown[]> | unknown[] {
    const value = runtime.unwrapValue(target);
    const t = runtime.type(value);
    switch (t) {
      case 'null':
        return [];

      case 'object':
        return runtime.muxAll(
          [Object.values(value as Record<string, unknown>)],
          next,
        );

      case 'array':
        return runtime.muxAll([value as unknown[]], next);

      case 'string':
        return runtime.muxAll([[...(value as string)]], next);

      default:
    }

    if (params.optional) return [];
    throw new JPLTypeError('cannot iterate over %s (%*<100v)', t, value);
  },

  /** { optional: boolean } */
  map(runtime: JPLRuntime, params: UserJPLSelectorParams): JPLSelectorParams {
    return {
      optional: runtime.assertType(params.optional, 'boolean'),
    };
  },
};
