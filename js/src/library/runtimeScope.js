import { applyObject } from './apply';
import RuntimeSignal from './runtimeSignal';

/**
 * Runtime scope for managing contextual runtime information.
 *
 * Fields of the scope must not be modified directly.
 * Instead, `JPLRuntimeScope.next` should be used to create a sub scope with the specified modifications.
 */
class JPLRuntimeScope {
  constructor(presets = {}) {
    this._state = {
      signal: presets.signal ?? new RuntimeSignal(),
      vars: presets.vars ?? {},
    };
  }

  get signal() {
    return this._state.signal;
  }

  get vars() {
    return this._state.vars;
  }

  /** Inherit the next scope based on the specified modifications */
  next(modifications) {
    return new JPLRuntimeScope({
      signal: modifications.signal ?? this.signal,
      vars: applyObject(this.vars, Object.entries(modifications.vars ?? {})),
    });
  }
}

export default JPLRuntimeScope;
