import type JPLRuntimeSignal from '@/library/runtimeSignal';
import type JPLRuntime from '@/runtime';
import { JPLTypeError } from '../library';

export default function builtin(
  runtime: JPLRuntime,
  _signal: JPLRuntimeSignal,
  next: (output: unknown) => Promise<unknown[]> | unknown[],
  input: unknown,
): Promise<unknown[]> | unknown[] {
  const value = runtime.unwrapValue(input);
  const t = runtime.type(value);

  switch (t) {
    case 'array':
      return next((value as unknown[]).map((_, i) => i));

    case 'object':
      return next(Object.keys(value as Record<string, unknown>));

    default:
  }

  throw new JPLTypeError('%s (%*<100v) has no keys', t, value);
}
