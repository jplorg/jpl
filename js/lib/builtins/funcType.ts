import type JPLRuntimeSignal from '@/library/runtimeSignal';
import type JPLRuntime from '@/runtime';

export default function builtin(
  runtime: JPLRuntime,
  _signal: JPLRuntimeSignal,
  next: (output: unknown) => Promise<unknown[]> | unknown[],
  input: unknown,
): Promise<unknown[]> | unknown[] {
  return next(runtime.type(input));
}
