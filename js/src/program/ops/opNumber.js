export default {
  /** { number: number } */
  op(runtime, input, params, scope, next) {
    return next(params.number ?? 0, scope);
  },

  /** { number: number } */
  map(runtime, params) {
    return {
      number: runtime.assertType(params.number, 'number'),
    };
  },
};
