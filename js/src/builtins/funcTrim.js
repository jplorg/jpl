import { JPLTypeError } from '../library';

async function builtin(runtime, signal, next, input) {
  const t = runtime.type(input);

  switch (t) {
    case 'string':
      return next(await runtime.alterValue(input, (value) => value.trim()));

    default:
  }

  throw new JPLTypeError('%s (%*<100v) cannot be trimmed', t, runtime.unwrapValue(input));
}

export default builtin;
