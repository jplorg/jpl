import { orphanFunction } from '../library';
import systemInterpreter from './system';

/**
 * Create an orphan JPL function from the specified source program string.
 *
 * Some optional scope presets may be specified, e.g. for allowing the function access to some specified variables.
 * Other than that, the function does not have access to any external variables.
 *
 * @param {string[]} argNames Names to bind the function arguments to
 * @param {string} source Source program string to execute
 * @param {{ vars: object? }} [presets] Optional scope presets
 */
async function parseFunction(argNames, source, presets) {
  const instructions = await systemInterpreter.parseInstructions(source);

  return orphanFunction(argNames, instructions, presets);
}

export default parseFunction;
