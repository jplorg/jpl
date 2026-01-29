import { JPLTypeError } from '../../../library';
import { call } from '../utils';

export default {
  /** { from: [op], to: [op], optional: boolean } */
  async op(runtime, input, target, params, scope, next) {
    const value = runtime.unwrapValue(target);
    const tv = runtime.type(value);

    const [froms, tos] = await Promise.all([
      runtime.executeInstructions(params.from ?? [], [input], scope),
      runtime.executeInstructions(params.to ?? [], [input], scope),
    ]);

    return runtime.muxAll([runtime.unwrapValues(froms), runtime.unwrapValues(tos)], (from, to) => {
      const tf = runtime.type(from);
      const tt = runtime.type(to);
      switch (tv) {
        case 'null':
          return next(null);

        case 'array':
          if (['number', 'null'].includes(tf) && ['number', 'null'].includes(tt)) {
            return next(value.slice(Math.trunc(from ?? 0), Math.trunc(to ?? value.length)));
          }
          break;

        case 'string':
          if (['number', 'null'].includes(tf) && ['number', 'null'].includes(tt)) {
            const chars = [...value];
            return next(
              chars.slice(Math.trunc(from ?? 0), Math.trunc(to ?? chars.length)).join(''),
            );
          }
          break;

        default:
      }

      if (params.optional) return [];
      throw new JPLTypeError(
        'cannot slice %s (%*<100v) with %s (%*<100v) and %s (%*<100v)',
        tv,
        value,
        tf,
        from,
        tt,
        to,
      );
    });
  },

  /** { from: function, to: function, optional: boolean } */
  map(runtime, params) {
    return {
      from: call(params.from),
      to: call(params.to),
      optional: runtime.assertType(params.optional, 'boolean'),
    };
  },
};
