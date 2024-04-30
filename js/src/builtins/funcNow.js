function builtin(runtime, signal, next) {
  return next(Date.now());
}

export default builtin;
