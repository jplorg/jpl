import { scopedFunction } from '../../library';
import { call } from './utils';

export default {
  /** { argNames: [string], pipe: [op] } */
  op(runtime, input, params, scope, next) {
    return next(scopedFunction(params.argNames ?? [], params.pipe ?? [], scope), scope);
  },

  /** { argNames: [string], pipe: function } */
  map(runtime, params) {
    return {
      argNames: runtime.muxOne([params.argNames], (entry) => runtime.assertType(entry, 'string')),
      pipe: call(params.pipe),
    };
  },
};
