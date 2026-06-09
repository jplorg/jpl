import { JPLTypeError, type JPLRuntimeSignal } from '@/library';
import type JPLRuntime from '@/runtime';

export default async function builtin(
  runtime: JPLRuntime,
  _signal: JPLRuntimeSignal,
  next: (output: unknown) => Promise<unknown[]> | unknown[],
  input: unknown,
): Promise<unknown[]> {
  const t = runtime.type(input);

  switch (t) {
    case 'string':
      return next(
        await runtime.alterValue(input, (value) => (value as string).trimEnd()),
      );

    default:
  }

  throw new JPLTypeError(
    '%s (%*<100v) cannot be trimmed',
    t,
    runtime.unwrapValue(input),
  );
}
