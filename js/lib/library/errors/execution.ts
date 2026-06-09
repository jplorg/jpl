import JPLError from './error';

/**
 * JPL error type for execution errors.
 *
 * All error types that infer this type can be caught in a program.
 */
export default class JPLExecutionError extends JPLError {
  #value: unknown;

  constructor(message?: string, name: string = 'JPLExecutionError') {
    super(message, name);
    this.#value = typeof message === 'string' && message ? message : null;
  }

  protected set value(value: unknown) {
    this.#value = value;
  }

  get value(): unknown {
    return this.#value;
  }
}
