import applyDefaults from '../applyDefaults';
import JPLRuntime, { applyRuntimeDefaults } from '../runtime';
import ops from './ops';

const defaultOptions = {};

export function applyProgramDefaults(options = {}, defaults = {}) {
  return applyDefaults(options, defaults);
}

/** JPL program */
class JPLProgram {
  constructor(instructions, options) {
    this.options = applyProgramDefaults(options?.program, defaultOptions);

    this.instructions = instructions;
    this.ops = ops;

    this.runtimeOptions = options?.runtime;
  }

  /**
   * Run the program with the provided inputs and runtime options.
   * The program throws a JPLExecutionError (see `JPLExecutionError.is`) for runtime failures.
   * Other errors may be thrown when execution fails.
   */
  run = async (inputs, options) => {
    const runtime = new JPLRuntime(this, {
      runtime: applyRuntimeDefaults(options?.runtime, this.runtimeOptions),
    });

    const normalizedInputs = runtime.normalizeValues(inputs, 'program inputs');

    const outputs = await runtime.executeProgram(this.instructions, normalizedInputs);
    return runtime.stripJSON(outputs);
  };
}

export default JPLProgram;
