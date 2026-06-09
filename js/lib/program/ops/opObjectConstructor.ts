import type { JPLInstructionParams } from '@/library/definition';
import type JPLRuntime from '@/runtime';
import { JPLTypeError, type JPLRuntimeScope } from '../../library';
import type { UserJPLInstructionParams } from '../params';
import { call } from './utils';

export default {
  /** { fields: [{ key: [op], value: [op], optional: boolean }] } */
  async op(
    runtime: JPLRuntime,
    input: unknown,
    params: JPLInstructionParams,
    scope: JPLRuntimeScope,
    next: (
      output: unknown,
      scope: JPLRuntimeScope,
    ) => Promise<unknown[]> | unknown[],
  ): Promise<unknown[]> {
    const fields = await runtime.muxAsync(
      [params.fields ?? []],
      async (field) => {
        const [keys, values] = await Promise.all([
          runtime.executeInstructions(field.key ?? [], [input], scope),
          runtime.executeInstructions(field.value ?? [], [input], scope),
        ]);

        return runtime.muxAll(
          [runtime.unwrapValues(keys), values],
          (key, value) => {
            const t = runtime.type(key);
            switch (t) {
              case 'string':
                return [[key, value]];

              default:
            }

            if (field.optional) return [];
            throw new JPLTypeError(
              'cannot use %s (%*<100v) as object key',
              t,
              key,
            );
          },
        );
      },
    );

    return runtime.muxAll(fields, (...entries) =>
      next(Object.fromEntries(entries), scope),
    );
  },

  /** { fields: [{ key: function, value: function, optional: boolean }] } */
  map(
    runtime: JPLRuntime,
    params: UserJPLInstructionParams,
  ): JPLInstructionParams {
    return {
      fields: runtime.muxOne([params.fields!], (entry) => ({
        key: call(entry.key),
        value: call(entry.value),
        optional: runtime.assertType(entry.optional, 'boolean'),
      })),
    };
  },
};
