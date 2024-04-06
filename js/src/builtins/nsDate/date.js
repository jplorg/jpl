import { JPLTypeError } from '../../library';
import JPLDate from './typeDate';

function unwrapType(runtime, type, purpose, v) {
  const t = runtime.type(v);
  if (t !== type) {
    throw new JPLTypeError('%s (%*<100v) cannot be used as %s', t, runtime.unwrapValue(v), purpose);
  }
  return runtime.unwrapValue(v);
}

function funcGet(getter) {
  return async function builtin(runtime, signal, next, input, ...args) {
    if (!JPLDate.is(input)) throw new JPLTypeError('%*<100v is not a date', input);

    return next(await getter(runtime, input, ...args));
  };
}

function funcSet(setter) {
  return async function builtin(runtime, signal, next, input, ...args) {
    if (!JPLDate.is(input)) throw new JPLTypeError('%*<100v is not a date', input);

    return next(await input.alterDate((date, format) => setter(runtime, date, format, ...args)));
  };
}

export const locale = funcGet((runtime, input) => input.date.locale);
export const setLocale = funcSet((runtime, date, format, arg0) => [
  date.setLocale(unwrapType(runtime, 'string', 'a locale string', arg0 ?? null)),
  format,
  true,
]);

export const zone = funcGet((runtime, input) => input.date.zoneName);
export const setZone = funcSet((runtime, date, format, arg0) => [
  date.setZone(unwrapType(runtime, 'string', 'a zone string', arg0 ?? null)),
  format,
  false,
]);

export const year = funcGet((runtime, input) => input.date.year);
export const quarter = funcGet((runtime, input) => input.date.quarter);
export const month = funcGet((runtime, input) => input.date.month);
export const day = funcGet((runtime, input) => input.date.day);
export const hour = funcGet((runtime, input) => input.date.hour);
export const minute = funcGet((runtime, input) => input.date.minute);
export const second = funcGet((runtime, input) => input.date.second);
export const millisecond = funcGet((runtime, input) => input.date.millisecond);
export const weekYear = funcGet((runtime, input) => input.date.weekYear);
export const weekNumber = funcGet((runtime, input) => input.date.weekNumber);
export const weekday = funcGet((runtime, input) => input.date.weekday);
export const ordinal = funcGet((runtime, input) => input.date.ordinal);
export const monthShort = funcGet((runtime, input) => input.date.monthShort);
export const monthLong = funcGet((runtime, input) => input.date.monthLong);
export const weekdayShort = funcGet((runtime, input) => input.date.weekdayShort);
export const weekdayLong = funcGet((runtime, input) => input.date.weekdayLong);
export const offset = funcGet((runtime, input) => input.date.offset);
export const daysInMonth = funcGet((runtime, input) => input.date.daysInMonth);
export const daysInYear = funcGet((runtime, input) => input.date.daysInYear);
export const weeksInWeekYear = funcGet((runtime, input) => input.date.weeksInWeekYear);

// units year, month, day, ordinal, weekYear, weekNumber, weekday, hour, minute, second, millisecond
export const set = funcSet((runtime, date, format, arg0) => [
  date.set(runtime.stripValue(unwrapType(runtime, 'object', 'an object', arg0 ?? null))),
  format,
  false,
]);

// units year, quarter, month, week, day, hour, minute, second, millisecond
export const startOf = funcSet((runtime, date, format, arg0) => [
  date.startOf(runtime.stripValue(unwrapType(runtime, 'string', 'a unit', arg0 ?? null))),
  format,
  false,
]);

// units year, quarter, month, week, day, hour, minute, second, millisecond
export const endOf = funcSet((runtime, date, format, arg0) => [
  date.endOf(runtime.stripValue(unwrapType(runtime, 'string', 'a unit', arg0 ?? null))),
  format,
  false,
]);

// units year, month, day, hour, minute, second, millisecond
// this function ignores zones, which compares the local calendar date
export const hasSame = funcGet((runtime, input, arg0, arg1) => {
  if (!JPLDate.is(arg0)) throw new JPLTypeError('%*<100v is not a date', arg0 ?? null);
  return input.date.hasSame(arg0.date, unwrapType(runtime, 'string', 'a unit', arg1 ?? null));
});

export const toSeconds = funcGet((runtime, input) => input.date.toUnixInteger());
export const toMillis = funcGet((runtime, input) => input.date.toMillis());
export const toObject = funcGet((runtime, input) => input.date.toObject());
