# [JPL - JSON Processing Language](index.md) - [Builtin functions](builtins.md) - Date and time

Dates are represented as a shadow type in JPL. Shadow types appear as one primitive type in JPL expressions while being treated as another primitive type when being returned from the program or interpolated into strings.

Dates are treated as numbers (milliseconds since epoch) in JPL expressions, while they are represented as strings when being returned. This means that you can simply add or subtract milliseconds to a date.

However, be aware, that a shadow type must always be at the leftmost side of an expression if you do not want it to be lost. For example, `date.now() + 1000` results in a date, whereas `1000 + date.now()` results in a simple number.

## Utilities

### `date.dates()`

Returns only those inputs that are dates.

### `date.isDate()`

Returns `true` if the input is a date, `false` otherwise.

## Constructing dates

### `date.now(options)`

Creates a date from the current instant, in the system's time zone.

`options` can be optionally provided:

- `options.zone`: any IANA zone (e.g. `Europe/Berlin`), a fixed-offset name (e.g. `UTC+2`), `local` or `utc`
- `options.locale`: Locale (e.g. `de-DE`)

### `date.parseISO(options)`

Creates a date by parsing the input as an ISO 8601 string.

`options` can be optionally provided:

- `options.zone`: any IANA zone (e.g. `Europe/Berlin`), a fixed-offset name (e.g. `UTC+2`), `local` or `utc`
- `options.setZone`: If set to `true`, the zone is overridden with a fixed-offset zone specified in the input string, if it specifies one
- `options.locale`: Locale (e.g. `de-DE`)

### `date.parse(format, options)`

Creates a date by parsing the input string with the specified format string.

See [here](https://moment.github.io/luxon/#/parsing?id=table-of-tokens) for a list of supported tokens.

`options` can be optionally provided:

- `options.zone`: any IANA zone (e.g. `Europe/Berlin`), a fixed-offset name (e.g. `UTC+2`), `local` or `utc`
- `options.setZone`: If set to `true`, the zone is overridden with a fixed-offset zone specified in the input string, if it specifies one
- `options.locale`: Locale (e.g. `de-DE`)

### `date.fromObject(options)`

Creates a date from the input object with keys like `year` and `hour`.

`options` can be optionally provided:

- `options.zone`: any IANA zone (e.g. `Europe/Berlin`), a fixed-offset name (e.g. `UTC+2`), `local` or `utc`
- `options.locale`: Locale (e.g. `de-DE`)

### `date.fromSeconds(options)`

Creates a date from a number of seconds since the epoch (1 January 1970 00:00:00 UTC), provided as input.

`options` can be optionally provided:

- `options.zone`: any IANA zone (e.g. `Europe/Berlin`), a fixed-offset name (e.g. `UTC+2`), `local` or `utc`
- `options.locale`: Locale (e.g. `de-DE`)

### `date.fromMillis(options)`

Creates a date from a number of milliseconds since the epoch (1 January 1970 00:00:00 UTC), provided as input.

`options` can be optionally provided:

- `options.zone`: any IANA zone (e.g. `Europe/Berlin`), a fixed-offset name (e.g. `UTC+2`), `local` or `utc`
- `options.locale`: Locale (e.g. `de-DE`)

## Acting on dates

## `date.formatISO()`

Formats the input date as an ISO 8601-compliant string.

## `date.formatISODate()`

Formats the input date as an ISO 8601-compliant date string.

A date may loose some of its details when formatting. E.g. `date.now() | date.formatISODate() | date.formatISO()` looses all information about the date's time.

## `date.formatISOTime()`

Formats the input date as an ISO 8601-compliant time string.

A date may loose some of its details when formatting. E.g. `date.now() | date.formatISODate() | date.formatISO()` looses all information about the date's time.

## `date.formatISOWeekDate()`

Formats the input date as an ISO 8601-compliant week date string.

A date may loose some of its details when formatting. E.g. `date.now() | date.formatISODate() | date.formatISO()` looses all information about the date's time.

## `date.format(format)`

Formates the input date with the specified format string.

See [here](https://moment.github.io/luxon/#/parsing?id=table-of-tokens) for a list of supported tokens.

A date may loose some of its details when formatting. E.g. `date.now() | date.formatISODate() | date.formatISO()` looses all information about the date's time.

## `date.locale()`

Returns the locale of the input date.

## `date.setLocale(locale)`

Sets the locale of the input date (e.g. `de-DE`).

## `date.zone()`

Returns the zone of the input date.

## `date.setZone(zone)`

Sets the zone of the input date.
`zone` can by any IANA zone (e.g. `Europe/Berlin`), a fixed-offset name (e.g. `UTC+2`), `local` or `utc`

## Unit selectors

The following unit selectors return the specified unit of the input date.

- `date.year()`
- `date.quarter()`
- `date.month()`
- `date.day()`
- `date.hour()`
- `date.minute()`
- `date.second()`
- `date.millisecond()`
- `date.weekYear()`
- `date.weekNumber()`
- `date.weekday()`
- `date.ordinal()`
- `date.monthShort()`
- `date.monthLong()`
- `date.weekdayShort()`
- `date.weekdayLong()`
- `date.offset()`
- `date.daysInMonth()`
- `date.daysInYear()`
- `date.weeksInWeekYear()`

## `date.set(values)`

Sets the input date's values of the units specified in the `values` object.

The following units can be used: `year`, `month`, `day`, `ordinal`, `weekYear`, `weekNumber`, `weekday`, `hour`, `minute`, `second` and `millisecond`

## `date.startOf(unit)`

Sets the input date to the beginning of the specified unit of time.

The following units can be used: `year`, `quarter`, `month`, `week`, `day`, `hour`, `minute`, `second` and `millisecond`

## `date.endOf(unit)`

Sets the input date to the end (meaning the last millisecond) of the specified unit of time.

The following units can be used: `year`, `quarter`, `month`, `week`, `day`, `hour`, `minute`, `second` and `millisecond`

## `date.hasSame(other, unit)`

Returns `true` if the input date is the same unit of time as `other`, `false` otherwise.
Higher-order units must also be identical for this function to return `true`.
Note that time zones are ignored in this comparison, which compares the local calendar time.
Use `date.setZone()` to convert one of the dates if needed.

The following units can be used: `year`, `month`, `day`, `hour`, `minute`, `second` and `millisecond`

## `date.toObject()`

Returns the input date as an object with the date's units `year`, `month`, `day`, `hour`, `minute`, `second` and `millisecond`.

## `date.toSeconds()`

Returns the epoch seconds of the input date as a (whole) number.

## `date.toMillis()`

Returns the epoch milliseconds of the input date as a number.
