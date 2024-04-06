import { DateTime } from 'luxon';
import { JPLTypeError } from '../../library';
import JPLDate from './typeDate';

function builtin(runtime, signal, next, input, arg0) {
  const options = runtime.stripValue(arg0 ?? null);
  const to = runtime.type(options);

  if (!['null', 'object'].includes(to)) {
    throw new JPLTypeError('expected an object for date options but got %s (%*<100v)', to, options);
  }

  let date = DateTime.now();
  if (options) {
    if (Object.hasOwn(options, 'zone')) date = date.setZone(options.zone);
    if (Object.hasOwn(options, 'locale')) date = date.setLocale(options.locale);
  }
  return next(new JPLDate(date));
}

export default builtin;
