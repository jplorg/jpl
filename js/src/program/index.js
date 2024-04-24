import applyDefaults from '../applyDefaults';
import { DEFINITION_VERSION, DEFINITION_VERSION_MAJOR, DEFINITION_VERSION_MINOR } from '../library';
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
    throw new Error('invalid program definition');
  }

  if (+major !== DEFINITION_VERSION_MAJOR || +minor > DEFINITION_VERSION_MINOR) {
    throw new Error(
      `unsupported program definition v${major}.${minor} - this version of JPL only supports v${DEFINITION_VERSION_MAJOR} (up to v${DEFINITION_VERSION})`,
    );
  }
}

/** JPL program */
class JPLProgram {
  constructor(programDefinition, options) {
    validateDefinition(programDefinition);

    this.options = applyProgramDefaults(options?.program, defaultOptions);

    this.definition = programDefinition;
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

    const outputs = await runtime.executeProgram(this.definition.instructions, normalizedInputs);
    return runtime.stripJSON(outputs);
  };
}

export default JPLProgram;
