function builtin(runtime, signal, next, input) {
  return next(runtime.stringifyJSON(input, true));
}

export default builtin;
