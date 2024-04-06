export default {
  /** { value: any } */
  op(runtime, input, params, scope, next) {
    return next(params.value, scope);
  },

  /** { value: any } */
  map(runtime, params) {
    return {
      value: runtime.normalizeValue(params.value),
    };
  },
};
