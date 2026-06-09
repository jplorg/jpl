import { applyObject } from './apply';
import JPLRuntimeSignal from './runtimeSignal';

export type JPLRuntimeScopeConfig = {
  signal?: JPLRuntimeSignal;
  vars?: { [name: string]: unknown };
};

/**
 * Runtime scope for managing contextual runtime information.
 *
 * Fields of the scope must not be modified directly.
 * Instead, `JPLRuntimeScope.next` should be used to create a sub scope with the specified modifications.
 */
export default class JPLRuntimeScope {
  #signal;
  #vars;

  constructor(presets: JPLRuntimeScopeConfig = {}) {
    this.#signal = presets.signal ?? new JPLRuntimeSignal();
    this.#vars = presets.vars ?? {};
  }

  get signal() {
    return this.#signal;
  }

  get vars() {
    return this.#vars;
  }

  /** Inherit the next scope based on the specified modifications */
  next(modifications: JPLRuntimeScopeConfig) {
    return new JPLRuntimeScope({
      signal: modifications.signal ?? this.signal,
      vars: applyObject(
        this.vars,
        Object.entries(modifications.vars ?? {}),
      ) as { [name: string]: unknown },
    });
  }
}
