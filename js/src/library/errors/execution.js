import JPLError from './error';

/**
 * JPL error type for execution errors.
 *
 * All error types that infer this type can be caught in a program.
 */
class JPLExecutionError extends JPLError {
  constructor(message, name = 'JPLExecutionError') {
    super(message, name);
    this.value = typeof message === 'string' && message ? message : null;
  }
}

export default JPLExecutionError;
