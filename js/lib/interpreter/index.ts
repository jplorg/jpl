import applyDefaults from '@/applyDefaults';
import { DEFINITION_VERSION, type JPLDefinition, type Pipe } from '@/library';
import JPLProgram, {
  applyProgramDefaults,
  type JPLProgramOptions,
} from '@/program';
import { applyRuntimeDefaults, type JPLRuntimeOptions } from '@/runtime';
import { parseEntrypoint } from './parse';

export type JPLInterpreterConfig = {
  interpreter?: JPLInterpreterOptions;
  program?: JPLProgramOptions;
  runtime?: JPLRuntimeOptions;
};

export type JPLInterpreterOptions = Record<string, never>;

const defaultOptions: JPLInterpreterOptions = {};

export function applyInterpreterDefaults(
  options: JPLInterpreterOptions = {},
  defaults: JPLInterpreterOptions = {},
): JPLInterpreterOptions {
  return applyDefaults(options, defaults);
}

/** JPL interpreter */
export default class JPLInterpreter {
  #options;
  #programOptions;
  #runtimeOptions;

  constructor(options?: JPLInterpreterConfig) {
    this.#options = applyInterpreterDefaults(
      options?.interpreter,
      defaultOptions,
    );
    this.#programOptions = options?.program;
    this.#runtimeOptions = options?.runtime;
  }

  /** Return the interpreter's options */
  get options() {
    return this.#options;
  }

  /** Parse the specified source program string into a reusable JPLProgram instance */
  parse = async (
    source: string,
    options?: JPLInterpreterConfig,
  ): Promise<JPLProgram> => {
    const instructions = await this.parseInstructions(source);

    const programDefinition: JPLDefinition = {
      version: DEFINITION_VERSION,
      instructions,
    };

    return new JPLProgram(programDefinition, {
      program: applyProgramDefaults(options?.program, this.#programOptions),
      runtime: applyRuntimeDefaults(options?.runtime, this.#runtimeOptions),
    });
  };

  /** Parse the specified source program string */
  parseInstructions = async (source: string): Promise<Pipe> => {
    const { ops: instructions } = await parseEntrypoint(source, 0, {
      interpreter: this,
    });
    return instructions;
  };
}
