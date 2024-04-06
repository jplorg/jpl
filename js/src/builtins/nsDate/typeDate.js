import { DateTime } from 'luxon';
import { JPLFatalError, JPLRuntimeError, JPLType, typeOf, unwrap } from '../../library';

const defaultFormat = { type: 'iso' };

class JPLDate extends JPLType {
  constructor(value, format, trimmed) {
    super();

    this._date = assertDate(value);
    if (format) {
      this.format = format;
      this.trimmed = !!trimmed;
    } else {
      this.format = defaultFormat;
      this.trimmed = true;
    }
  }

  get date() {
    if (!this.trimmed) {
      this._date = trim(this._date, this.format);
    }
    return this._date;
  }

  get value() {
    if (!Object.hasOwn(this, '_iv')) {
      this._iv = this.date.toMillis();
    }
    return this._iv;
  }

  get json() {
    if (!Object.hasOwn(this, '_ij')) {
      this._ij = toFormat(this.date, this.format);
    }
    return this._ij;
  }

  async alter(updater) {
    const nextValue = await super.alter(updater);
    if (JPLDate.is(nextValue)) return nextValue;
    if (typeOf(nextValue) !== 'number') return nextValue;
    return new JPLDate(
      DateTime.fromMillis(unwrap(nextValue), copyOptions(this.date)),
      this.format,
      false,
    );
  }

  async alterDate(updater) {
    const [nextDate, nextFormat, trimmed] = await updater(this.date, this.format);
    if (nextDate === this.date && compareFormat(this.format, nextFormat)) return this;
    return new JPLDate(nextDate, nextFormat, trimmed);
  }
}

export default JPLDate;

function trim(date, format) {
  switch (format.type) {
    case 'iso':
      return date;

    default:
      return assertDate(fromFormat(toFormat(date, format), format, copyOptions(date)));
  }
}

function assertDate(date) {
  if (!DateTime.isDateTime(date)) throw new JPLFatalError(`${date} is not a valid date`);

  if (!date.isValid) {
    throw new JPLRuntimeError(
      `invalid date${date.invalidReason ? ` (${date.invalidReason})` : ''}${
        date.invalidExplanation ? `: ${date.invalidExplanation}` : ''
      }`,
    );
  }

  return date;
}

function compareFormat(a, b) {
  const format = a ?? defaultFormat;
  const other = b ?? defaultFormat;
  if (format.type !== other.type) return false;

  switch (format.type) {
    case 'iso':
    case 'iso-date':
    case 'iso-week-date':
    case 'iso-time':
      return true;

    case 'layout':
      return format.layout === other.layout;

    default:
      throw new JPLRuntimeError('invalid format type %s', format.type);
  }
}

function toFormat(date, format) {
  switch (format.type) {
    case 'iso':
      return date.toISO();

    case 'iso-date':
      return date.toISODate();

    case 'iso-week-date':
      return date.toISOWeekDate();

    case 'iso-time':
      return date.toISOTime();

    case 'layout':
      return date.toFormat(format.layout, { locale: date.locale });

    default:
      throw new JPLRuntimeError('invalid format type %s', format.type);
  }
}

function fromFormat(value, format, options) {
  switch (format.type) {
    case 'iso':
    case 'iso-date':
    case 'iso-week-date':
    case 'iso-time':
      return DateTime.fromISO(value, options);

    case 'layout':
      return DateTime.fromFormat(value, format.layout, options);

    default:
      throw new JPLRuntimeError('invalid format type %s', format.type);
  }
}

function copyOptions(date) {
  return {
    zone: date.zone,
    locale: date.locale,
    outputCalendar: date.outputCalendar,
    numberingSystem: date.numberingSystem,
  };
}
