import getBuiltins from './builtins';
import createSingleton from './createSingleton';
import JPLInterpreter, {
  applyInterpreterDefaults,
  type JPLInterpreterConfig,
  type JPLInterpreterOptions,
} from './interpreter';
import parseFunction from './interpreter/parseFunction';
import JPLProgram, {
  applyProgramDefaults,
  type JPLOPHandler,
  type JPLOPSubHandler,
  type JPLProgramConfig,
  type JPLProgramOptions,
} from './program';
import {
  applyRuntimeDefaults,
  type JPLRuntimeConfig,
  type JPLRuntimeOptions,
} from './runtime';

const getGlobalOptions = createSingleton(async () => ({
  interpreter: {},

  program: {},

  runtime: {
    vars: { ...(await getBuiltins()) },
  },
}));

const parse = async (
  source: string,
  options?: JPLInterpreterConfig,
): Promise<JPLProgram> => {
  const globalOptions = await getGlobalOptions();
  const interpreter = new JPLInterpreter({
    interpreter: applyInterpreterDefaults(
      options?.interpreter,
      globalOptions.interpreter,
    ),
    program: applyProgramDefaults(options?.program, globalOptions.program),
    runtime: applyRuntimeDefaults(options?.runtime, globalOptions.runtime),
  });
  return interpreter.parse(source);
};

const run = async (
  source: string,
  inputs: unknown[],
  options?: JPLInterpreterConfig,
): Promise<unknown[]> => {
  const program = await parse(source, options);
  return program.run(inputs);
};

export * from './library';
export * from './program/params';
export {
  getBuiltins,
  getGlobalOptions as getOptions,
  parse,
  parseFunction,
  run,
  type JPLInterpreterConfig,
  type JPLInterpreterOptions,
  type JPLOPHandler,
  type JPLOPSubHandler,
  type JPLProgramConfig,
  type JPLProgramOptions,
  type JPLRuntimeConfig,
  type JPLRuntimeOptions,
};
