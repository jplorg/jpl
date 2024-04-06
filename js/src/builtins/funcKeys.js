import { JPLTypeError } from '../library';

function builtin(runtime, signal, next, input) {
  const value = runtime.unwrapValue(input);
  const t = runtime.type(value);

  switch (t) {
    case 'array':
      return next(value.map((_, i) => i));

    case 'object':
      return next(Object.keys(value));

    default:
  }

  throw new JPLTypeError('%s (%*<100v) has no keys', t, value);
}

export default builtin;
