import { JPLSyntaxError } from '../library';

const setAZ = 'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ';
const setDigit = '0123456789';
const setVarFirst = `${setAZ}_$`;
const setVarRest = `${setVarFirst}${setDigit}`;
const setVarAll = setVarRest;
const setWhitespace = ' \xa0\r\n\t';
const setHex = `${setDigit}abcdefABCDEF`;

/** Walk whitespace at i */
export function walkWhitespace(src, i, c) {
  let n = i;
  let is = false;

  for (;;) {
    let m = match(src, n, c, { phrase: '#' });
    if (m.is) {
      ({ i: n } = m);
      for (;;) {
        const set = matchSet(src, n, c, { set: '\r\n' });
        if (set.is) {
          ({ i: n } = set);
          break;
        }

        const end = eot(src, n, c);
        if (end.is) {
          ({ i: n } = end);
          break;
        }

        n += 1;
      }
      is = true;
      continue;
    }

    m = match(src, n, c, { phrase: '/*' });
    if (m.is) {
      ({ i: n } = m);
      for (;;) {
        m = match(src, n, c, { phrase: '*/' });
        if (m.is) {
          ({ i: n } = m);
          break;
        }

        const end = eot(src, n, c);
        if (end.is) {
          ({ i: n } = end);
          return errorUnexpectedToken(src, n, c, {
            operator: 'comment',
            message: 'incomplete comment',
          });
        }

        n += 1;
      }
      is = true;
      continue;
    }

    const set = matchSet(src, n, c, { set: setWhitespace });
    ({ i: n } = set);
    if (!set.is) break;
    is = true;
  }

  return { i: n, is };
}

/**
 * Check if src matches phrase at i.
 * If notBeforeSet is provided, the phrase is only considered to be matching if it is not immediately followed by one of the symbols of the set.
 * If spaceBefore is true, the phrase is only considered to be matching if there is space (not to be confused with whitespace) before i (or i is 0).
 * If spaceAfter is true, the phrase is only considered to be matching if there is space (not to be confused with whitespace) after the phrase (or the end is reached).
 * The returned index is positioned directly after the phrase or at the first unmatched character, even if spaceAfter is true.
 * Space is defined as anything that cannot occur in a variable name.
 */
export function match(src, i, c, { phrase, notBeforeSet, spaceBefore, spaceAfter }) {
  let n = i;

  if (spaceBefore && i > 0) {
    if (!matchSet(src, n - 1, c, { set: setVarAll, exclusive: true }).is)
      return { i: n, is: false };
  }

  for (let j = 0; j < phrase.length; j += 1, n += 1) {
    const end = eot(src, n, c);
    if (end.is || src[n] !== phrase[j]) return { i: end.i, is: false };
  }

  if (!eot(src, n, c).is) {
    if (notBeforeSet && matchSet(src, n, c, { set: notBeforeSet }).is) return { i: n, is: false };

    if (spaceAfter && !matchSet(src, n, c, { set: setVarAll, exclusive: true }).is)
      return { i: n, is: false };
  }

  return { i: n, is: true };
}

/** The same as match, but also walk whitespace after the match if present */
export function matchWord(src, i, c, { phrase, notBeforeSet, spaceBefore, spaceAfter }) {
  let n = i;

  const m = match(src, n, c, { phrase, notBeforeSet, spaceBefore, spaceAfter });
  if (!m.is) return { i: m.i, is: false };
  ({ i: n } = m);

  ({ i: n } = walkWhitespace(src, n, c));

  return { i: n, is: true };
}

/** Check if i is a hexadecimal value */
export function hex(src, i, c) {
  return matchSet(src, i, c, { set: setHex });
}

/**
 * Check if i is one of the chars in set.
 * If exclusive is true, i must not be one of the chars in set for the check to succeed.
 */
export function matchSet(src, i, c, { set, exclusive }) {
  const end = eot(src, i, c);
  if (end.is) return { i: end.i, is: false };

  const char = src[i];
  if (set.includes(char) === !exclusive) return { i: i + 1, is: true, value: char };

  return { i, is: false, value: '' };
}

/** Parse variable selector at i */
export function variable(src, i, c) {
  let n = i;
  let is = false;
  let value = '';

  for (;;) {
    const set = matchSet(src, n, c, { set: !is ? setVarFirst : setVarRest });
    ({ i: n } = set);
    if (!set.is) break;
    is = true;
    value += set.value;
  }

  if (is) ({ i: n } = walkWhitespace(src, n, c));

  return { i: n, is, value };
}

/** The same as variable, but also check that the result is not a reserved term */
export function safeVariable(src, i, c) {
  const v = variable(src, i, c);

  if (!v.is) return v;

  switch (v.value) {
    case 'and':
    case 'catch':
    case 'func':
    case 'elif':
    case 'else':
    case 'end':
    case 'false':
    case 'if':
    case 'not':
    case 'null':
    case 'or':
    case 'then':
    case 'true':
    case 'try':
      return { i, is: false, value: v.value, reserved: true };

    default:
      return v;
  }
}

/** Check if i is at the end of src */
export function eot(src, i) {
  if (i >= src.length) return { i: src.length, is: true };
  return { i, is: false };
}

/** Get (zero based) line and column for i */
export function whereIs(src, i) {
  const lines = src.substring(0, i).split(/\r?\n|\r/);
  const line = lines[lines.length - 1];
  return { i, line: lines.length - 1, column: line.length };
}

/** Get a descriptive text highlighting i */
export function highlightLocation(src, i, c, { area = 25 } = {}) {
  const s = Math.max(Math.min(i, src.length - 1 - area), area) - area;
  const e = Math.min(s + area + 1 + area, src.length);
  const view = src
    .substring(s, e)
    .replaceAll(/\r?\n|\r/g, '⏎')
    .replaceAll('\t', '→');

  let prefix = ' > ';
  if (s > 0) prefix += '…';

  let suffix = '';
  if (e < src.length) suffix = '…';

  const description = `${prefix}${view}${suffix}\n${' '.repeat(prefix.length + (i - s))}^ here`;

  return { i, value: description };
}

/** Throw an error caused by an unexpected token at i */
export function errorUnexpectedToken(src, i, c, { operator, message } = {}) {
  let errorMessage;
  if (eot(src, i, c).is) errorMessage = 'unexpected EOT';
  else {
    const { line, column } = whereIs(src, i, c);
    errorMessage = `unexpected token '${src[i]}' at line ${line + 1}, column ${column + 1}`;
  }
  if (operator) errorMessage += ` while parsing ${operator}`;
  if (message) errorMessage += `: ${message}`;
  errorMessage += `\n${highlightLocation(src, i, c).value}`;
  throw new JPLSyntaxError(errorMessage);
}
