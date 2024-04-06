export default {
  /** {} */
  op(runtime, input, params, scope, next) {
    return next(true, scope);
  },

  /** {} */
  map() {
    return {};
  },
};
