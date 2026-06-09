import type JPLRuntimeSignal from '@/library/runtimeSignal';
import type JPLRuntime from '@/runtime';

export default function builtin(
  _runtime: JPLRuntime,
  _signal: JPLRuntimeSignal,
  next: (output: unknown) => Promise<unknown[]> | unknown[],
  _input: unknown,
): Promise<unknown[]> | unknown[] {
  return next(Date.now());
}
