import applyDefaults from '@/applyDefaults';
import {
  DEFINITION_VERSION,
  DEFINITION_VERSION_MAJOR,
  DEFINITION_VERSION_MINOR,
  JPLFatalError,
  type JPLDefinition,
} from '@/library';
import JPLRuntime, {
  applyRuntimeDefaults,
  type JPLRuntimeOptions,
} from '@/runtime';
import type { JPLOPHandler } from './handler';
import ops from './ops';

export type JPLProgramConfig = {
  program?: JPLProgramOptions;
  runtime?: JPLRuntimeOptions;
};

export type JPLProgramOptions = Record<string, never>;

const defaultOptions: JPLProgramOptions = {};

export function applyProgramDefaults(
  options: JPLProgramOptions = {},
  defaults: JPLProgramOptions = {},
): JPLProgramOptions {
  return applyDefaults(options, defaults);
}

function validateDefinition(programDefinition: JPLDefinition): void {
  const version = programDefinition?.version;
  const [, major, minor] =
    typeof version === 'string' ? (/^(\d+)\.(\d+)?$/g.exec(version) ?? []) : [];

  if (!major || !minor || !Array.isArray(programDefinition.instructions)) {
    throw new JPLFatalError('invalid program definition');
  }

  if (
    +major !== DEFINITION_VERSION_MAJOR ||
    +minor > DEFINITION_VERSION_MINOR
  ) {
    throw new JPLFatalError(
      `unsupported program definition v${major}.${minor} - this version of JPL only supports v${DEFINITION_VERSION_MAJOR} (up to v${DEFINITION_VERSION})`,
    );
  }
}

/** JPL program */
export default class JPLProgram {
  #options;
  #runtimeOptions;
  #definition;
  #ops: { [op: string]: JPLOPHandler };

  constructor(programDefinition: JPLDefinition, options?: JPLProgramConfig) {
    validateDefinition(programDefinition);

    this.#options = applyProgramDefaults(options?.program, defaultOptions);
    this.#runtimeOptions = options?.runtime;

    this.#definition = programDefinition;
    this.#ops = ops;
  }

  /** Return the program's options */
  get options() {
    return this.#options;
  }

  /**
   * Return the program's definition.
   * The definition can be serialized as JSON to be reused in other JPL implementations.
   */
  get definition() {
    return this.#definition;
  }

  /** Return the program's OPs */
  get ops() {
    return this.#ops;
  }

  /**
   * Run the program with the provided inputs and runtime options.
   * The program throws a JPLExecutionError (see `JPLExecutionError.is`) for runtime failures.
   * Other errors may be thrown when execution fails.
   */
  run = async (
    inputs: unknown[],
    options?: JPLProgramConfig,
  ): Promise<unknown[]> => {
    const runtime = new JPLRuntime(this, {
      runtime: applyRuntimeDefaults(options?.runtime, this.#runtimeOptions),
    });

    const normalizedInputs = runtime.normalizeValues(inputs, 'program inputs');

    const outputs = await runtime.execute(normalizedInputs);
    return runtime.stripJSON(outputs) as unknown[];
  };

  /**
   * Return a new program with the provided definition's instructions prepended to the program.
   */
  prepend = (programDefinition: JPLDefinition): JPLProgram => {
    validateDefinition(programDefinition);

    const mergedDefinition: JPLDefinition = {
      version: DEFINITION_VERSION,
      instructions: [
        ...programDefinition.instructions,
        ...this.definition.instructions,
      ],
    };

    return new JPLProgram(mergedDefinition, {
      program: this.#options,
      runtime: this.#runtimeOptions,
    });
  };

  /**
   * Return a new program with the provided definition's instructions appended to the program.
   */
  append = (programDefinition: JPLDefinition): JPLProgram => {
    validateDefinition(programDefinition);

    const mergedDefinition: JPLDefinition = {
      version: DEFINITION_VERSION,
      instructions: [
        ...this.definition.instructions,
        ...programDefinition.instructions,
      ],
    };

    return new JPLProgram(mergedDefinition, {
      program: this.#options,
      runtime: this.#runtimeOptions,
    });
  };
}

export { type JPLOPHandler, type JPLOPSubHandler } from './handler';
