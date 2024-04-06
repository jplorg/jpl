import { JPLTypeError } from '../library';

function builtin(runtime, signal, next, input, arg0) {
  const value = runtime.unwrapValue(input);
  const tv = runtime.type(value);

  const key = runtime.unwrapValue(arg0 ?? null);
  const tk = runtime.type(key);

  switch (tv) {
    case 'array':
      if (tk === 'number') return next(key >= 0 && key < value.length);
      break;

    case 'object':
      if (tk === 'string') return next(Object.hasOwn(value, key));
      break;

    default:
  }

  throw new JPLTypeError('%s (%*<100v) cannot have %s (%*<100v) as key', tv, value, tk, key);
}

export default builtin;
