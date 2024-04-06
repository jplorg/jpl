import { JPLTypeError } from '../library';

function builtin(runtime, signal, next, input, arg0) {
  const key = runtime.unwrapValue(input);
  const tk = runtime.type(key);

  const value = runtime.unwrapValue(arg0 ?? null);
  const tv = runtime.type(value);

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
