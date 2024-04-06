import { JPLTypeError } from '../library';

function builtin(runtime, signal, next, input, arg0) {
  const value = runtime.unwrapValue(input);
  const t = runtime.type(value);

  const token = runtime.unwrapValue(arg0 ?? null);
  const tt = runtime.type(token);

  switch (t) {
    case 'array':
      return next(value.some((item) => runtime.equals(item, token)));

    case 'string':
      if (tt === 'string') return next(value.includes(token));
      break;

    case 'object':
      return next(Object.values(value).some((item) => runtime.equals(item, token)));

    default:
  }

  throw new JPLTypeError('%s (%*<100v) cannot contain %s (%*<100v)', t, value, tt, token);
}

export default builtin;
