export default {
  /** {} */
  op(runtime, input, params, scope, next) {
    return next(false, scope);
  },

  /** {} */
  map() {
    return {};
  },
};
