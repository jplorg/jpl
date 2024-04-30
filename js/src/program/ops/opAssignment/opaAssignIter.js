import { JPLTypeError, applyCombinations } from '../../../library';

export default {
  /** { optional: boolean } */
  async op(runtime, input, target, params, scope, next) {
    const vt = runtime.unwrapValue(target);
    const tt = runtime.type(vt);
    switch (tt) {
      case 'object': {
        const items = Object.entries(vt);
        const outputs = await runtime.muxAsync([items], async (item) =>
          runtime.muxAsync([await next(item[1] ?? null)], (output) => {
            if (output === undefined) return item;
            return [item[0], output];
          }),
        );
        return runtime.muxAsync([applyCombinations(items, outputs)], (results) => {
          if (items === results) return target;
          return runtime.alterValue(target, () => Object.fromEntries(results));
        });
      }

      case 'array': {
        const items = vt;
        const outputs = await runtime.muxAsync([items], async (item) =>
          runtime.muxAsync([await next(item ?? null)], (output) => {
            if (output === undefined) return item;
            return output;
          }),
        );
        return runtime.muxAsync([applyCombinations(items, outputs)], (results) => {
          if (items === results) return target;
          return runtime.alterValue(target, () => results);
        });
      }

      case 'string': {
        const items = [...vt];
        const outputs = await runtime.muxAsync([items], async (item) =>
          runtime.muxAsync([await next(item)], (output) => {
            if (output === undefined) return item;
            const r = runtime.unwrapValue(output);
            const tr = runtime.type(r);
            switch (tr) {
              case 'null':
              case 'string':
                return r ?? ' ';

              default:
            }

            throw new JPLTypeError('cannot assign %s (%*<100v) to string (%*<100v)', tr, r, vt);
          }),
        );
        return runtime.muxAsync([applyCombinations(items, outputs)], (results) => {
          if (items === results) return target;
          return runtime.alterValue(target, () => results.join(''));
        });
      }

      default:
    }

    if (params.optional) return [undefined];
    throw new JPLTypeError('cannot iterate over %s (%*<100v) (assignment)', tt, vt);
  },

  /** { optional: boolean } */
  map(runtime, params) {
    return {
      optional: runtime.assertType(params.optional, 'boolean'),
    };
  },
};
