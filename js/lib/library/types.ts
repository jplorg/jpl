import { applyArray, applyObject } from './apply';
import adaptErrors, { adaptError } from './errors/adaptErrors';
import JPLFatalError from './errors/fatal';

/**
 * Generic type for handling special formatting on values.
 *
 * The value accessor should be used in JPL expressions, where type safety may be considered important.
 *
 * The json accessor should be used in program outputs, where the format of the value should be applied and the original type may be ignored.
 *
 * The alter function is called when the value is being changed by an operation.
 * The JPLType is responsible to decide whether the resulting value should be wrapped in a new JPLType or not.
 * It may also return itself if the value did not change.
 * Note that the resulting value may have another type than the original value.
 * Commonly, if multiple values wrapped in JPLTypes are involved in an operation, the leftmost operand should be altered whereas all other operands should be unwrapped using their value accessors, effectively loosing their formatting.
 *
 * However, both accessors should represent the same essential value.
 * For example, a JPLType that applies rounding to numbers with a fixed number of decimal digits, may return the rounded numeric value for its value accessor (e.g. `1`), whereas it may return a formatted string for its JSON accessor (e.g. `"1.00"`).
 * This allows this JPLType to be processed in JPL operations like generic numbers but resolves to formatted strings in the program output.
 */
export class JPLType {
  /** Detect whether the provided value is an instance of this JPLType or one of its subclasses */
  static is(value: unknown): boolean {
    return value instanceof this;
  }

  #v;

  constructor(value: unknown) {
    this.#v = normalizeInternalValue(value);
  }

  /** Resolve the internal value for usage in JPL operations */
  get value() {
    return this.#v;
  }

  /** Resolve the JSON value for usage in program outputs */
  get json() {
    return this.#v;
  }

  /**
   * Alter the internal value using the specified updater.
   * The result may or may not be wrapped into a new JPLType if desired.
   * The JPLType may return itself if the value did not change.
   */
  async alter(
    updater: (value: unknown) => Promise<unknown> | unknown,
  ): Promise<unknown> {
    return normalize(await updater(this.value));
  }

  /**
   * Support for `JSON.stringify` like procedures (see `JPLType.json`).
   *
   * It is generally recommended to override the `get json` member access function instead of this method.
   */
  toJSON(): unknown {
    try {
      return this.json;
    } catch (err) {
      throw adaptError(err);
    }
  }
}

/** Normalize the specified external value to be used in a JPLType */
function normalizeInternalValue(value: unknown): unknown {
  if (JPLType.is(value)) return (value as JPLType).value;
  return normalize(value);
}

/** Normalize the specified external value */
export function normalize(value: unknown): unknown {
  return strip(value, undefined, jplTypedStripper);
}

/** Unwrap the specified normalized value */
export function unwrap(value: unknown): unknown {
  if (JPLType.is(value)) return adaptErrors(() => (value as JPLType).value);
  return value;
}

/** Order which applies when comparing values with different types */
export const typeOrder = [
  'null',
  'function',
  'boolean',
  'number',
  'string',
  'array',
  'object',
];

/** Resolve the type of the specified normalized value */
export function typeOf(value: unknown): string {
  const v = unwrap(value);
  if (v === null) return 'null';
  if (Array.isArray(v)) return 'array';
  const t = typeof v;
  switch (t) {
    case 'function':
    case 'boolean':
    case 'number':
    case 'string':
    case 'object':
      return t;

    default:
      throw new JPLFatalError(`invalid type ${t} (${value})`);
  }
}

/** Assert the type for the specified unwrapped value */
export function assertType<T>(value: T, type: string): T {
  if (JPLType.is(value)) {
    throw new JPLFatalError(
      template(
        'unexpected type: JPLTypes (%*<100v) are not allowed here',
        value,
      ),
    );
  }
  const t = typeOf(value);
  if (t !== type) {
    throw new JPLFatalError(
      template(
        'unexpected type: %s (%*<100v) cannot be used as %s',
        t,
        value,
        type,
      ),
    );
  }
  return value;
}

/** Stringify the specified normalized value */
export function stringify(
  value: unknown,
  unescapeString?: boolean,
  escapeFunctions?: boolean,
): string {
  const rawValue = escapeFunctions
    ? strip(value, (_k, v) =>
        typeof unwrap(v) === 'function' ? '<function>' : v,
      )
    : strip(value);
  if (unescapeString && typeof rawValue === 'string') return rawValue;
  return JSON.stringify(rawValue);
}

/**
 * Format the specified template string.
 * The general form of a format is a percent sign, followed by optional flags, an optional width and a verb.
 *
 * Examples:
 * - `%10s`: Format the next replacement as a string and pads the result at the left with spaces to be at least 10 unicode codepoints long.
 * - `%*<10v`: Format the next replacement as a JSON value and truncates it, if it is longer then 10 unicode codepoints.
 *
 * The width specifies the desired field width and defaults to whatever is necessary to display the full replacement.
 * If the width is specified without any other corresponding flags, it is used for padding the field if necessary.
 *
 * Valid flags:
 *
 * - `*`: Do not pad the value even if it is shorter than the specified width
 * - `-`: Pad the value at the right rather than the left
 * - `<`: Truncate the value at the right if it is too long for the specified width
 *
 * Valid verbs:
 *
 * - `%`: Returns a literal `%`
 * - `s`: Format the next replacement as a string (like JSON, but does not escape strings)
 * - `v`: Format the next replacement as a JSON value
 */
export function template(tmpl: unknown, ...replacements: unknown[]): string {
  let i = 0;
  return displayValue(tmpl).replace(
    /%([*\-<]+)?([1-9][0-9]*)?(.)/g,
    (match, flags, width, verb) => {
      // verbs without replacement
      switch (verb) {
        case '%':
          return '%';
        default:
      }
      // verbs with replacement
      const value = replacements[i] ?? null;
      i += 1;
      let result;
      switch (verb) {
        case 's':
          result = displayValue(value);
          break;
        case 'v':
          result = strictDisplayValue(value);
          break;
        default:
          throw new JPLFatalError(`format ${match} has unknown verb ${verb}`);
      }
      let pad = true;
      let padRight = false;
      let trunc = false;
      [...(flags ?? [])].forEach((flag) => {
        switch (flag) {
          case '*':
            pad = false;
            break;
          case '-':
            padRight = true;
            break;
          case '<':
            trunc = true;
            break;
          default:
            throw new JPLFatalError(`format ${match} has unknown flag ${flag}`);
        }
      });
      const w = +(width ?? 0);
      if (w > 0) {
        const rl = [...result].length;
        if (pad && rl < w) {
          const padding = ' '.repeat(w - rl);
          result = padRight ? `${result}${padding}` : `${padding}${result}`;
        } else if (trunc && rl > w) {
          result = `${result.substring(0, w - 1)}…`;
        }
      }
      return result;
    },
  );
}

/** Format the specified normalized value as a string */
export function displayValue(value: unknown): string {
  return stringify(value, true, true);
}

/** Format the specified normalized value as a string, without removing escaping */
export function strictDisplayValue(value: unknown): string {
  return stringify(value, false, true);
}

/**
 * Unwrap the specified value similar to `JSON.stringify`.
 * However, unlike with `JSON.stringify`, object member functions `toJSON` are not supported as they could interfere with user defined members.
 * Instead, JPLTypes are being unwrapped by default.
 *
 * A custom stripper can be provided to customize the behavior.
 */
export function strip(
  value: unknown,
  replacer?: (k: string, v: unknown) => unknown,
  stripper: (
    k: string | null,
    v: unknown,
    iter: (k: string, v: unknown) => unknown,
  ) => unknown = jplJSONStripper,
) {
  const iter = (k: string | null, v: unknown): unknown => {
    let r = v;
    if (replacer) r = replacer(k ?? '', r);
    return stripper(k, r, iter);
  };
  return iter(null, value);
}

/** Stripper that allows JPLTypes and normalized values */
export function jplTypedStripper(
  k: string | null,
  v: unknown,
  iter: (k: string, v: unknown) => unknown,
): unknown {
  if (JPLType.is(v)) return v;
  if (typeof v === 'function') return v;
  return rawStripper(k, v, iter);
}

/** Stripper that only allows normalized values and unwraps JPLTypes */
export function jplStripper(
  k: string | null,
  v: unknown,
  iter: (k: string, v: unknown) => unknown,
): unknown {
  let r = v;
  if (JPLType.is(r)) r = (r as JPLType).toJSON();
  if (typeof r === 'function') return r;
  return rawStripper(k, r, iter);
}

/** Stripper that allows JSON like values and unwraps JPLTypes */
export function jplJSONStripper(
  k: string | null,
  v: unknown,
  iter: (k: string, v: unknown) => unknown,
): unknown {
  let r = v;
  if (JPLType.is(r)) r = (r as JPLType).toJSON();
  return rawStripper(k, r, iter);
}

/** Stripper that allows JSON like values and calls `toJSON` similar to `JSON.stringify` */
export function jsonStripper(
  k: string | null,
  v: unknown,
  iter: (k: string, v: unknown) => unknown,
): unknown {
  let r = v as unknown;
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  if (typeof (r as any)?.toJSON === 'function') r = (r as any).toJSON();
  return rawStripper(k, r, iter);
}

/** Stripper that allows JSON like values */
export function rawStripper(
  k: string | null,
  v: unknown,
  iter: (k: string, v: unknown) => unknown,
): unknown {
  const top = k == null;

  switch (typeof v) {
    case 'string':
    case 'boolean':
      return v;
    case 'number':
      if (!Number.isFinite(v)) return null;
      return v;
    case 'undefined':
    case 'function':
      return top ? null : undefined;
    case 'object':
      if (v === null) return null;
      if (Array.isArray(v)) {
        return applyArray(
          v,
          Array.from(new Array(v.length), (_, i) => [
            i,
            iter(`${i}`, v[i]) ?? null,
          ]),
        );
      }
      return applyObject(
        v,
        Object.entries(v).map(([i, s]) => {
          const entry = iter(i, s);
          if (entry === undefined) return [i];
          return [i, entry];
        }),
      );
    default:
      throw new JPLFatalError(`unexpected ${typeof v}`);
  }
}
