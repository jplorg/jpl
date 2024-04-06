import { JPLTypeError, applyArray } from '../../../library';
import { call } from '../utils';

export default {
  /** { pipe: [op], optional: boolean } */
  async op(runtime, input, target, params, scope, next) {
    let items;
    const vt = runtime.unwrapValue(target);
    const tt = runtime.type(target);
    switch (tt) {
      case 'null':
      case 'object':
      case 'array':
        items = vt;
        break;

      case 'string':
        items = [...vt];
        break;

      default:
        if (params.optional) return [undefined];
        throw new JPLTypeError(
          'cannot access fields of %s (%*<100v) (assignment)',
          tt,
          runtime.unwrapValue(vt),
        );
    }

    const fields = await runtime.executeInstructions(params.pipe, [input], scope, (output) => [
      output,
    ]);

    const iter = async (from, source) => {
      // Call stack decoupling - This is necessary as some browsers (i.e. Safari) have very limited call stack sizes which result in stack overflow exceptions in certain situations.
      await undefined;

      scope.signal.checkHealth();

      if (from >= fields.length) return [source];

      const t = tt === 'string' ? tt : runtime.type(source);
      const field = runtime.unwrapValue(fields[from]);
      const tf = runtime.type(field);
      switch (t) {
        case 'null': {
          switch (tf) {
            case 'string':
              return runtime.muxAll([await next(null)], (output) => {
                if (output === undefined) return iter(from + 1, source);
                return iter(from + 1, { [field]: output });
              });

            case 'number': {
              const i = Math.trunc(field);
              return runtime.muxAll([await next(null)], (output) => {
                if (output === undefined) return iter(from + 1, source);
                return iter(from + 1, applyArray([], [[i, output]]));
              });
            }

            default:
          }

          break;
        }

        case 'object':
          if (tf === 'string') {
            const item = source[field];
            return runtime.muxAll([await next(item ?? null)], (output) => {
              if (output === undefined || item === output) return iter(from + 1, source);
              return iter(from + 1, { ...source, [field]: output });
            });
          }
          break;

        case 'array':
          if (tf === 'number') {
            const i = Math.trunc(field);
            const item = source[i >= 0 ? i : source.length + i];
            return runtime.muxAll([await next(item ?? null)], (output) => {
              if (output === undefined || item === output) return iter(from + 1, source);
              return iter(from + 1, applyArray(source, [[i, output]]));
            });
          }
          break;

        case 'string':
          if (tf === 'number') {
            const i = Math.trunc(field);
            const item = source[i >= 0 ? i : source.length + i];
            return runtime.muxAll([await next(item ?? null)], (output) => {
              if (output === undefined || item === output) return iter(from + 1, source);
              const r = runtime.unwrapValue(output);
              const tr = runtime.type(r);
              switch (tr) {
                case 'null':
                case 'string':
                  return iter(from + 1, applyArray(source, [[i, output ?? ' ']], ' '));

                default:
              }

              throw new JPLTypeError(
                'cannot assign %s (%*<100v) to string (%*<100v)',
                tr,
                r,
                source.join(''),
              );
            });
          }
          break;

        default:
      }

      if (params.optional) return iter(from + 1, source);
      throw new JPLTypeError(
        'cannot access field of %s (%*<100v) with %s (%*<100v) (assignment)',
        t,
        t === 'string' ? source.join('') : source,
        tf,
        field,
      );
    };

    return runtime.muxAsync([await iter(0, items)], (results) => {
      if (items === results) return target;

      const t = tt === 'string' ? tt : runtime.type(results);
      switch (t) {
        case 'null':
          return undefined;

        case 'object':
        case 'array':
          return runtime.alterValue(target, () => results);

        case 'string':
          return runtime.alterValue(target, () => results.join(''));

        default:
      }

      return undefined;
    });
  },

  /** { value: function, optional: boolean } */
  map(runtime, params) {
    return {
      pipe: call(params.value),
      optional: runtime.assertType(params.optional, 'boolean'),
    };
  },
};
