import type JPLRuntimeSignal from '@/library/runtimeSignal';
import type JPLRuntime from '@/runtime';
import { JPLTypeError } from '../library';

export default function builtin(
  runtime: JPLRuntime,
  _signal: JPLRuntimeSignal,
  next: (output: unknown) => Promise<unknown[]> | unknown[],
  input: unknown,
) {
  const value = runtime.unwrapValue(input);
  const t = runtime.type(value);

  switch (t) {
    case 'null':
      return next(0);

    case 'array':
      return next((value as unknown[]).length);

    case 'string':
      return next([...(value as string)].length);

    case 'object':
      return next(Object.keys(value as Record<string, unknown>).length);

    default:
  }

  throw new JPLTypeError('%s (%*<100v) has no length', t, value);
}
