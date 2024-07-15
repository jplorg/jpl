import applyDefaults from '../applyDefaults';
import {
  DEFINITION_VERSION,
  DEFINITION_VERSION_MAJOR,
  DEFINITION_VERSION_MINOR,
  JPLFatalError,
} from '../library';
import JPLRuntime, { applyRuntimeDefaults } from '../runtime';
import ops from './ops';

const defaultOptions = {};

export function applyProgramDefaults(options = {}, defaults = {}) {
  return applyDefaults(options, defaults);
}

function validateDefinition(programDefinition) {
  const version = programDefinition?.version;
  const [, major, minor] =
    typeof version === 'string' ? /^(\d+)\.(\d+)?$/g.exec(version) ?? [] : [];

  if (!major || !minor || !Array.isArray(programDefinition.instructions)) {
    throw new JPLFatalError('invalid program definition');
  }

  if (+major !== DEFINITION_VERSION_MAJOR || +minor > DEFINITION_VERSION_MINOR) {
    throw new JPLFatalError(
      `unsupported program definition v${major}.${minor} - this version of JPL only supports v${DEFINITION_VERSION_MAJOR} (up to v${DEFINITION_VERSION})`,
    );
  }
}

/** JPL program */
class JPLProgram {
  constructor(programDefinition, options) {
    validateDefinition(programDefinition);

    this._options = applyProgramDefaults(options?.program, defaultOptions);
    this._runtimeOptions = options?.runtime;

    this._definition = programDefinition;
    this._ops = ops;
  }

  /** Return the program's options */
  get options() {
    return this._options;
  }

  /**
   * Return the program's definition.
   * The definition can be serialized as JSON to be reused in other JPL implementations.
   */
  get definition() {
    return this._definition;
  }

  /** Return the program's OPs */
  get ops() {
    return this._ops;
  }

  /**
   * Run the program with the provided inputs and runtime options.
   * The program throws a JPLExecutionError (see `JPLExecutionError.is`) for runtime failures.
   * Other errors may be thrown when execution fails.
   */
  run = async (inputs, options) => {
    const runtime = new JPLRuntime(this, {
      runtime: applyRuntimeDefaults(options?.runtime, this._runtimeOptions),
    });

    const normalizedInputs = runtime.normalizeValues(inputs, 'program inputs');

    const outputs = await runtime.execute(normalizedInputs);
    return runtime.stripJSON(outputs);
  };

  /**
   * Return a new program with the provided definition's instructions prepended to the program.
   */
  prepend = (programDefinition) => {
    validateDefinition(programDefinition);

    const mergedDefinition = {
      version: DEFINITION_VERSION,
      instructions: [...programDefinition.instructions, ...this.definition.instructions],
    };

    return new JPLProgram(mergedDefinition, {
      program: this._options,
      runtime: this._runtimeOptions,
    });
  };

  /**
   * Return a new program with the provided definition's instructions appended to the program.
   */
  append = (programDefinition) => {
    validateDefinition(programDefinition);

    const mergedDefinition = {
      version: DEFINITION_VERSION,
      instructions: [...this.definition.instructions, ...programDefinition.instructions],
    };

    return new JPLProgram(mergedDefinition, {
      program: this._options,
      runtime: this._runtimeOptions,
    });
  };
}

export default JPLProgram;
