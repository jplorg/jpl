import { DateTime } from 'luxon';
import { JPLTypeError } from '../../library';
import JPLDate from './typeDate';

function builtin(runtime, signal, next, input, arg0, arg1) {
  const value = runtime.unwrapValue(input);
  const t = runtime.type(value);
  if (t !== 'string') throw new JPLTypeError('cannot use %s (%*<100v) as a date string', t, value);

  const format = runtime.unwrapValue(arg0 ?? null);
  const tf = runtime.type(format);
  if (tf !== 'string') {
    throw new JPLTypeError('cannot use %s (%*<100v) as a date format string', tf, format);
  }

  const options = runtime.stripValue(arg1 ?? null);
  const to = runtime.type(options);
  if (!['null', 'object'].includes(to)) {
    throw new JPLTypeError('expected an object for date options but got %s (%*<100v)', to, options);
  }

  const dateOpts = {};
  if (options) {
    if (Object.hasOwn(options, 'zone')) dateOpts.zone = options.zone;
    if (Object.hasOwn(options, 'setZone')) dateOpts.setZone = options.setZone;
    if (Object.hasOwn(options, 'locale')) dateOpts.locale = options.locale;
  }
  return next(
    new JPLDate(
      DateTime.fromFormat(value, format, dateOpts),
      { type: 'layout', layout: format },
      true,
    ),
  );
}

export default builtin;
