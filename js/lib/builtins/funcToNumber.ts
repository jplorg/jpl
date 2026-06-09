import {
  JPLTypeConversionError,
  JPLTypeError,
  type JPLRuntimeSignal,
} from '@/library';
import type JPLRuntime from '@/runtime';

export default function builtin(
  runtime: JPLRuntime,
  _signal: JPLRuntimeSignal,
  next: (output: unknown) => Promise<unknown[]> | unknown[],
  input: unknown,
): Promise<unknown[]> | unknown[] {
  const value = runtime.unwrapValue(input);
  const t = runtime.type(value);

  switch (t) {
    case 'number':
      return next(value);

    case 'string': {
      const result = runtime.normalizeValue(+(value as string));
      if (value === '' || runtime.type(result) !== 'number') {
        throw new JPLTypeConversionError(
          '%s (%*<100v) does not contain a valid number',
          t,
          value,
        );
      }
      return next(result);
    }

    default:
  }

  throw new JPLTypeError('%s (%*<100v) cannot be parsed as a number', t, value);
}
