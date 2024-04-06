import JPLDate from './typeDate';

function builtin(runtime, signal, next, input) {
  return next(JPLDate.is(input));
}

export default builtin;
