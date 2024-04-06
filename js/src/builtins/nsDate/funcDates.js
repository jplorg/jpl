import JPLDate from './typeDate';

function builtin(runtime, signal, next, input) {
  return JPLDate.is(input) ? next(input) : [];
}

export default builtin;
