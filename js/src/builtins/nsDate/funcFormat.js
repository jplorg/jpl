import { JPLTypeError } from '../../library';
import JPLDate from './typeDate';

async function builtin(runtime, signal, next, input, arg0) {
  if (!JPLDate.is(input)) throw new JPLTypeError('%*<100v is not a date', input);

  const format = runtime.unwrapValue(arg0 ?? null);
  const tf = runtime.type(format);
  if (tf !== 'string') {
    throw new JPLTypeError('cannot use %s (%*<100v) as a date format string', tf, format);
  }

  return next(await input.alterDate((date) => [date, { type: 'layout', layout: format }, false]));
}

export default builtin;
