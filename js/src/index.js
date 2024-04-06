import getBuiltins from './builtins';
import createSingleton from './createSingleton';
import JPLInterpreter, { applyInterpreterDefaults } from './interpreter';
import parseFunction from './interpreter/parseFunction';
import * as library from './library';
import { applyProgramDefaults } from './program';
import { applyRuntimeDefaults } from './runtime';

const getGlobalOptions = createSingleton(async () => ({
  interpreter: {},

  program: {},

  runtime: {
    vars: { ...(await getBuiltins()) },
  },
}));

const parse = async (source, options) => {
  const globalOptions = await getGlobalOptions();
  const interpreter = new JPLInterpreter({
    interpreter: applyInterpreterDefaults(options?.interpreter, globalOptions.interpreter),
    program: applyProgramDefaults(options?.program, globalOptions.program),
    runtime: applyRuntimeDefaults(options?.runtime, globalOptions.runtime),
  });
  return interpreter.parse(source);
};

const run = async (source, inputs, options) => {
  const program = await parse(source, options);
  return program.run(inputs);
};

export default {
  ...library,
  run,
  parse,
  parseFunction,
  getBuiltins,
  getOptions: getGlobalOptions,
};

export * from './library';
export { getBuiltins, getGlobalOptions as getOptions, parse, parseFunction, run };
