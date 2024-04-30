export default {
  /** { value: any } */
  op(runtime, input, params, scope, next) {
    return next(params.value ?? null, scope);
  },

  /** { value: any } */
  map(runtime, params) {
    return {
      value: runtime.normalizeValue(params.value),
    };
  },
};
