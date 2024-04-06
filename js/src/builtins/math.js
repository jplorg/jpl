import { JPLTypeError } from '../library';

function unwrapNumber(runtime, v) {
  const t = runtime.type(v);
  if (t !== 'number') {
    throw new JPLTypeError(
      '%s (%*<100v) cannot be used for mathematical operations',
      t,
      runtime.unwrapValue(v),
    );
  }
  return runtime.unwrapValue(v);
}

function funcMath(alter) {
  return async function builtin(runtime, signal, next, input, ...args) {
    const t = runtime.type(input);
    if (t !== 'number') {
      throw new JPLTypeError(
        '%s (%*<100v) cannot be used for mathematical operations',
        t,
        runtime.unwrapValue(input),
      );
    }
    return next(await runtime.alterValue(input, (value) => alter(runtime, value, ...args)));
  };
}

export const pow = funcMath((runtime, value, arg0) => value ** unwrapNumber(runtime, arg0 ?? null));
export const sqrt = funcMath((runtime, value) => Math.sqrt(value));
export const exp = funcMath((runtime, value) => Math.exp(value));
export const log = funcMath((runtime, value) => Math.log(value));
export const log10 = funcMath((runtime, value) => Math.log10(value));

export const sin = funcMath((runtime, value) => Math.sin(value));
export const cos = funcMath((runtime, value) => Math.cos(value));
export const tan = funcMath((runtime, value) => Math.tan(value));
export const asin = funcMath((runtime, value) => Math.asin(value));
export const acos = funcMath((runtime, value) => Math.acos(value));
export const atan = funcMath((runtime, value) => Math.atan(value));

export const ceil = funcMath((runtime, value) => Math.ceil(value));
export const floor = funcMath((runtime, value) => Math.floor(value));
export const round = funcMath((runtime, value) => Math.round(value));
export const trunc = funcMath((runtime, value) => Math.trunc(value));

export const abs = funcMath((runtime, value) => Math.abs(value));
