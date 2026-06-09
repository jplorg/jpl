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
  const tv = runtime.type(value);

  const key = runtime.unwrapValue(arg0 ?? null);
  const tk = runtime.type(key);

  switch (tv) {
    case 'array':
      if (tk === 'number')
        return next(
          (key as number) >= 0 && (key as number) < (value as unknown[]).length,
        );
      break;

    case 'object':
      if (tk === 'string')
        return next(
          Object.hasOwn(value as Record<string, unknown>, key as string),
        );
      break;

    default:
  }

  throw new JPLTypeError(
    '%s (%*<100v) cannot have %s (%*<100v) as key',
    tv,
    value,
    tk,
    key,
  );
}
