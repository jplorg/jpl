import { JPLErrorEnclosure, JPLTypeError, adaptError } from '../../../library';
import { call } from '../utils';

export default {
  /** { args: [[op]], bound: boolean, optional: boolean } */
  async op(runtime, input, target, params, scope, next) {
    const value = runtime.unwrapValue(target);
    const t = runtime.type(value);
    switch (t) {
      case 'function': {
        const args = await runtime.muxAsync([params.args ?? []], (arg) =>
          runtime.executeInstructions(arg ?? [], [input], scope),
        );

        return runtime.muxAll(args, async (...a) => {
          try {
            const fnNext = async (output) => {
              try {
                return await next(output);
              } catch (err) {
                throw new JPLErrorEnclosure(err);
              }
            };

            if (params.bound) {
              return await value(runtime, scope.signal, fnNext, a[0] ?? null, ...a.slice(1));
            }

            return await value(runtime, scope.signal, fnNext, input, ...a);
          } catch (err) {
            if (JPLErrorEnclosure.is(err)) throw err.inner;
            throw adaptError(err);
          }
        });
      }

      default:
    }

    if (params.optional) return [];
    throw new JPLTypeError('cannot execute %s (%*<100v)', t, value);
  },

  /** { args: [function], bound: boolean, optional: boolean } */
  map(runtime, params) {
    return {
      args: runtime.muxOne([params.args], (entry) => call(entry)),
      bound: runtime.assertType(params.bound, 'boolean'),
      optional: runtime.assertType(params.optional, 'boolean'),
    };
  },
};
