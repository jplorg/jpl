export default {
  /** { value: string } */
  op(runtime, input, params, scope, next) {
    return next(params.value, scope);
  },

  /** { value: string } */
  map(runtime, params) {
    return {
      value: runtime.assertType(params.value, 'string'),
    };
  },
};