import { JPLTypeError } from '../library';

function builtin(runtime, signal, next, input, arg0) {
  const value = runtime.unwrapValue(input);
  const t = runtime.type(value);

  const token = runtime.unwrapValue(arg0 ?? null);
  const tt = runtime.type(token);

  switch (t) {
    case 'string':
      if (tt === 'string') return next(value.endsWith(token));
      break;

    default:
  }

  throw new JPLTypeError('%s (%*<100v) cannot end with %s (%*<100v)', t, value, tt, token);
}

export default builtin;
