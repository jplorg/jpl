import { JPLTypeError } from '../library';

function builtin(runtime, signal, next, input) {
  const value = runtime.unwrapValue(input);
  const t = runtime.type(value);

  switch (t) {
    case 'null':
      return next(0);

    case 'array':
      return next(value.length);

    case 'string':
      return next([...value].length);

    case 'object':
      return next(Object.keys(value).length);

    default:
  }

  throw new JPLTypeError('%s (%*<100v) has no length', t, value);
}

export default builtin;
