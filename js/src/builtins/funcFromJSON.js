import { JPLRuntimeError, JPLTypeError } from '../library';

function builtin(runtime, signal, next, input) {
  const value = runtime.unwrapValue(input);
  const t = runtime.type(value);

  switch (t) {
    case 'string': {
      let result;
      try {
        result = JSON.parse(value);
      } catch (err) {
        throw new JPLRuntimeError(err.message);
      }
      return next(result);
    }

    default:
  }

  throw new JPLTypeError('%s (%*<100v) cannot be parsed as JSON', t, value);
}

export default builtin;
