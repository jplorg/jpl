export default {
  /** {} */
  op(runtime, input, params, scope, next) {
    return next(!runtime.truthy(input), scope);
  },

  /** {} */
  map() {
    return {};
  },
};
