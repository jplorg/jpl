import { JPLTypeConversionError, JPLTypeError } from '../library';

function builtin(runtime, signal, next, input) {
  const value = runtime.unwrapValue(input);
  const t = runtime.type(value);

  switch (t) {
    case 'number':
      return next(value);

    case 'string': {
      const result = runtime.normalizeValue(+value);
      if (value === '' || runtime.type(result) !== 'number') {
        throw new JPLTypeConversionError('%s (%*<100v) does not contain a valid number', t, value);
      }
      return next(result);
    }

    default:
  }

  throw new JPLTypeError('%s (%*<100v) cannot be parsed as a number', t, value);
}

export default builtin;
