import type JPLRuntime from '@/runtime';
import type { Pipe } from './definition';
import JPLRuntimeScope, { type JPLRuntimeScopeConfig } from './runtimeScope';
import type JPLRuntimeSignal from './runtimeSignal';

export interface JPLFunc {
  (
    runtime: JPLRuntime,
    signal: JPLRuntimeSignal,
    next: (output: unknown) => Promise<unknown[]> | unknown[],
    input: unknown,
    ...args: unknown[]
  ): Promise<unknown[]> | unknown[];
}

function jplEnclosure(
  this: { argNames: string[]; pipe: Pipe; scope: JPLRuntimeScope },
  runtime: JPLRuntime,
  signal: JPLRuntimeSignal,
  next: (output: unknown) => Promise<unknown[]> | unknown[],
  input: unknown,
  ...args: unknown[]
): Promise<unknown[]> {
  const { argNames, pipe, scope } = this;

  return runtime.executeInstructions(
    pipe,
    [input],
    scope.next({
      signal,
      vars: Object.fromEntries(
        argNames.map((name, i) => [name, args[i] ?? null]),
      ),
    }),
    next,
  );
}

/**
 * Create a scoped JPL function from the specified instructions.
 *
 * The function is bound to the specified scope.
 *
 * @param {string[]} argNames Names to bind the function arguments to
 * @param {object[]} instructions Instructions to execute
 * @param {JPLRuntimeScope} scope Runtime scope to bind the function to
 */
export function scopedFunction(
  argNames: string[],
  instructions: Pipe,
  scope: JPLRuntimeScope,
): JPLFunc {
  return jplEnclosure.bind({ argNames, pipe: instructions, scope });
}

/**
 * Create an orphan JPL function from the specified instructions.
 *
 * Some optional scope presets may be specified, e.g. for allowing the function access to some specified variables.
 * Other than that, the function does not have access to any external variables.
 *
 * @param {string[]} argNames Names to bind the function arguments to
 * @param {object[]} instructions Instructions to execute
 * @param {{ vars: object? }} [presets] Optional scope presets
 */
export function orphanFunction(
  argNames: string[],
  instructions: Pipe,
  presets?: JPLRuntimeScopeConfig,
): JPLFunc {
  return scopedFunction(argNames, instructions, new JPLRuntimeScope(presets));
}

/**
 * Wrap the specified synchronous function to allow it to be used as a JPL function.
 *
 * If your function runs asynchronous routines, you should create the function yourself.
 * `next` must be called for each single result, which produces an array of results itself.
 * All resulting results must be returned as a single array.
 *
 * It is recommended to check the provided JPLRuntimeSignal in asynchronous routines to stop the routine when execution has been canceled:
 * ```
 * signal.checkHealth();
 * ```
 *
 * The following pattern can be used to merge the results of multiple asynchronous routines:
 * ```
 * const promises = ...;
 * return runtime.muxAll([promises], async (promise) => {
 *  const result = await promise;
 *  return next(result);
 * });
 * ```
 *
 * @param {(runtime: JPLRuntime, input: any, ...args: any[]) => (Promise<any[]> | any[])} fn The function to be wrapped
 */
export function nativeFunction(
  fn: (
    runtime: JPLRuntime,
    input: unknown,
    ...args: unknown[]
  ) => Promise<unknown[]> | unknown[],
): JPLFunc {
  return async (
    runtime: JPLRuntime,
    _signal: JPLRuntimeSignal,
    next: (output: unknown) => Promise<unknown[]> | unknown[],
    input: unknown,
    ...args: unknown[]
  ) => {
    const results = await fn(runtime, input, ...args);

    return runtime.muxAll([results], next);
  };
}
