import { displayValue, normalize, template } from '../types';
import JPLExecutionError from './execution';

function format(value, replacements) {
  return replacements.length > 0 ? template(value, ...replacements) : displayValue(value);
}

/**
 * JPL error type for generic runtime errors.
 *
 * `value` can by of any type.
 * If at least one replacement is specified, the value is formatted as a template.
 */
class JPLRuntimeError extends JPLExecutionError {
  constructor(value, ...replacements) {
    const message = format(value, replacements);

    super(message, 'JPLRuntimeError');
    this.value = replacements.length > 0 ? message : normalize(value);
  }
}

export default JPLRuntimeError;

/**
 * JPL runtime error type for type errors.
 *
 * `value` can by of any type.
 * If at least one replacement is specified, the value is formatted as a template.
 */
export class JPLTypeError extends JPLRuntimeError {
  constructor(value, ...replacements) {
    super(`TypeError - ${format(value, replacements)}`);
  }
}

/**
 * JPL runtime error type for reference errors.
 *
 * `value` can by of any type.
 * If at least one replacement is specified, the value is formatted as a template.
 */
export class JPLReferenceError extends JPLRuntimeError {
  constructor(value, ...replacements) {
    super(`ReferenceError - ${format(value, replacements)}`);
  }
}

/**
 * JPL runtime error type for zero division errors.
 *
 * `value` can by of any type.
 * If at least one replacement is specified, the value is formatted as a template.
 */
export class JPLZeroDivisionError extends JPLRuntimeError {
  constructor(value, ...replacements) {
    super(`ZeroDivisionError - ${format(value, replacements)}`);
  }
}

/**
 * JPL runtime error type for type conversion errors.
 *
 * `value` can by of any type.
 * If at least one replacement is specified, the value is formatted as a template.
 */
export class JPLTypeConversionError extends JPLRuntimeError {
  constructor(value, ...replacements) {
    super(`TypeConversionError - ${format(value, replacements)}`);
  }
}
