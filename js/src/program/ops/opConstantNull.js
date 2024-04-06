export default {
  /** {} */
  op(runtime, input, params, scope, next) {
    return next(null, scope);
  },

  /** {} */
  map() {
    return {};
  },
};
