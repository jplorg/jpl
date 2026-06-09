import {
  JPLRuntimeError,
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
    case 'string': {
      let result;
      try {
        result = JSON.parse(value as string);
      } catch (err) {
        throw new JPLRuntimeError((err as Error).message);
      }
      return next(result);
    }

    default:
  }

  throw new JPLTypeError('%s (%*<100v) cannot be parsed as JSON', t, value);
}
