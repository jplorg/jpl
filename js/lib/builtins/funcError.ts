import { JPLRuntimeError, type JPLRuntimeSignal } from '@/library';
import type JPLRuntime from '@/runtime';

export default function builtin(
  _runtime: JPLRuntime,
  _signal: JPLRuntimeSignal,
  _next: (output: unknown) => Promise<unknown[]> | unknown[],
  input: unknown,
): never {
  throw new JPLRuntimeError(input);
}
