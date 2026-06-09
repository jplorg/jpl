import type JPLRuntime from '@/runtime';
import { JPLTypeError, type JPLFunc } from '../library';

function unwrapNumber(runtime: JPLRuntime, v: unknown): number {
  const t = runtime.type(v);
  const u = runtime.unwrapValue(v);
  if (t !== 'number') {
    throw new JPLTypeError(
      '%s (%*<100v) cannot be used for mathematical operations',
      t,
      u,
    );
  }
  return u as number;
}

function funcMath(
  alter: (runtime: JPLRuntime, value: number, ...args: unknown[]) => number,
): JPLFunc {
  return async function builtin(runtime, _signal, next, input, ...args) {
    const t = runtime.type(input);
    if (t !== 'number') {
      throw new JPLTypeError(
        '%s (%*<100v) cannot be used for mathematical operations',
        t,
        runtime.unwrapValue(input),
      );
    }
    return next(
      await runtime.alterValue(input, (value) =>
        alter(runtime, value as number, ...args),
      ),
    );
  };
}

export const pow = funcMath(
  (runtime, value, arg0) => value ** unwrapNumber(runtime, arg0 ?? null),
);
export const sqrt = funcMath((_runtime, value) => Math.sqrt(value));
export const exp = funcMath((_runtime, value) => Math.exp(value));
export const log = funcMath((_runtime, value) => Math.log(value));
export const log10 = funcMath((_runtime, value) => Math.log10(value));

export const sin = funcMath((_runtime, value) => Math.sin(value));
export const cos = funcMath((_runtime, value) => Math.cos(value));
export const tan = funcMath((_runtime, value) => Math.tan(value));
export const asin = funcMath((_runtime, value) => Math.asin(value));
export const acos = funcMath((_runtime, value) => Math.acos(value));
export const atan = funcMath((_runtime, value) => Math.atan(value));

export const ceil = funcMath((_runtime, value) => Math.ceil(value));
export const floor = funcMath((_runtime, value) => Math.floor(value));
export const round = funcMath((_runtime, value) => Math.round(value));
export const trunc = funcMath((_runtime, value) => Math.trunc(value));

export const abs = funcMath((_runtime, value) => Math.abs(value));
