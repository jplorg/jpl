import { JPLRuntimeError } from '../library';

function builtin(runtime, signal, next, input) {
  throw new JPLRuntimeError(input);
}

export default builtin;
