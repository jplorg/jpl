import applyDefaults from '../applyDefaults';
import { DEFINITION_VERSION } from '../library';
import JPLProgram, { applyProgramDefaults } from '../program';
import { applyRuntimeDefaults } from '../runtime';
import { parseEntrypoint } from './parse';

const defaultOptions = {};

export function applyInterpreterDefaults(options = {}, defaults = {}) {
  return applyDefaults(options, defaults);
}

/** JPL interpreter */
class JPLInterpreter {
  constructor(options) {
    this._options = applyInterpreterDefaults(options?.interpreter, defaultOptions);
    this._programOptions = options?.program;
    this._runtimeOptions = options?.runtime;
  }

  /** Return the interpreter's options */
  get options() {
    return this._options;
  }

  /** Parse the specified source program string into a reusable JPLProgram instance */
  parse = async (source, options) => {
    const instructions = await this.parseInstructions(source);

    const programDefinition = {
      version: DEFINITION_VERSION,
      instructions,
    };

    return new JPLProgram(programDefinition, {
      program: applyProgramDefaults(options?.program, this._programOptions),
      runtime: applyRuntimeDefaults(options?.runtime, this._runtimeOptions),
    });
  };

  /** Parse the specified source program string */
  parseInstructions = async (source) => {
    const { ops: instructions } = await parseEntrypoint(source, 0, { interpreter: this });
    return instructions;
  };
}

export default JPLInterpreter;
