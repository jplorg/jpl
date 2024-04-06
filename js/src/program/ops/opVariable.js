import { JPLReferenceError } from '../../library';

export default {
  /** { name: string } */
  op(runtime, input, params, scope, next) {
    if (!Object.hasOwn(scope.vars, params.name))
      throw new JPLReferenceError('%s is not defined', params.name);
    return next(scope.vars[params.name], scope);
  },

  /** { name: string } */
  map(runtime, params) {
    return {
      name: runtime.assertType(params.name, 'string'),
    };
  },
};
