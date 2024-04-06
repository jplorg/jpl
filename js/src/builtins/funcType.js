function builtin(runtime, signal, next, input) {
  return next(runtime.type(input));
}

export default builtin;
