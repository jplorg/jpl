import { JPLTypeError, type JPLRuntimeSignal } from '@/library';
import type JPLRuntime from '@/runtime';

export default function builtin(
  runtime: JPLRuntime,
  _signal: JPLRuntimeSignal,
  next: (output: unknown) => Promise<unknown[]> | unknown[],
  input: unknown,
  arg0: unknown,
): Promise<unknown[]> | unknown[] {
  const value = runtime.unwrapValue(input);
  const t = runtime.type(value);

  const token = runtime.unwrapValue(arg0 ?? null);
  const tt = runtime.type(token);

  switch (t) {
    case 'array':
      return next(
        (value as unknown[]).some((item) => runtime.equals(item, token)),
      );

    case 'string':
      if (tt === 'string')
        return next((value as string).includes(token as string));
      break;

    case 'object':
      return next(
        Object.values(value as Record<string, unknown>).some((item) =>
          runtime.equals(item, token),
        ),
      );

    default:
  }

  throw new JPLTypeError(
    '%s (%*<100v) cannot contain %s (%*<100v)',
    t,
    value,
    tt,
    token,
  );
}
