import applyDefaults from '../applyDefaults';
import JPLProgram, { applyProgramDefaults } from '../program';
import { applyRuntimeDefaults } from '../runtime';
import { entrypoint } from './parse';

const defaultOptions = {};

export function applyInterpreterDefaults(options = {}, defaults = {}) {
  return applyDefaults(options, defaults);
}

/** JPL interpreter */
class JPLInterpreter {
  constructor(options) {
    this.options = applyInterpreterDefaults(options?.interpreter, defaultOptions);

    this.programOptions = options?.program;
    this.runtimeOptions = options?.runtime;
  }

  /** Parse the specified source program string into a reusable JPLProgram instance */
  parse = async (source, options) => {
    const instructions = await this.parseInstructions(source);

    return new JPLProgram(instructions, {
      program: applyProgramDefaults(options?.program, this.programOptions),
      runtime: applyRuntimeDefaults(options?.runtime, this.runtimeOptions),
    });
  };

  /** Parse the specified source program string */
  parseInstructions = async (source) => {
    const { ops: instructions } = await entrypoint(source, 0, { interpreter: this });
    return instructions;
  };
}

export default JPLInterpreter;
