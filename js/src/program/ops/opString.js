export default {
  /** { string: string } */
  op(runtime, input, params, scope, next) {
    return next(params.string ?? '', scope);
  },

  /** { string: string } */
  map(runtime, params) {
    return {
      string: runtime.assertType(params.string, 'string'),
    };
  },
};
