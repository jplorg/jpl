import { scopedFunction } from '../../library';

export default {
  /** { argNames: [string], pipe: [op] } */
  op(runtime, input, params, scope, next) {
    return next(scopedFunction(params.argNames, params.pipe, scope), scope);
  },

  /** { argNames: [string], value: [op] } */
  map(runtime, params) {
    return {
      argNames: runtime.muxOne([params.argNames], (entry) => runtime.assertType(entry, 'string')),
      value: runtime.assertType(params.value, 'array'),
    };
  },
};
