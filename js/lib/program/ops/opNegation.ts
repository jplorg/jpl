import type { JPLInstructionParams } from '@/library/definition';
import type JPLRuntime from '@/runtime';
import { JPLTypeError, type JPLRuntimeScope } from '../../library';
import type { UserJPLInstructionParams } from '../params';

export default {
  /** {} */
  async op(
    runtime: JPLRuntime,
    input: unknown,
    _params: JPLInstructionParams,
    scope: JPLRuntimeScope,
    next: (
      output: unknown,
      scope: JPLRuntimeScope,
    ) => Promise<unknown[]> | unknown[],
  ): Promise<unknown[]> {
    return next(
      await runtime.alterValue(input, (value) => {
        const t = runtime.type(value);
        switch (t) {
          case 'number':
            return -(value as number);

          default:
        }

        throw new JPLTypeError('%s (%*<100v) cannot be negated', t, value);
      }),
      scope,
    );
  },

  /** {} */
  map(
    _runtime: JPLRuntime,
    _params: UserJPLInstructionParams,
  ): JPLInstructionParams {
    return {};
  },
};
