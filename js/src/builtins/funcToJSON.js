function builtin(runtime, signal, next, input) {
  return next(runtime.stringifyJSON(input, false));
}

export default builtin;
