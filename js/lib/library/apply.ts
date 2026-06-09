/* eslint-disable @typescript-eslint/no-explicit-any */

/**
 * Apply all changes immutably to the source object.
 *
 * @param {[key: string, value: any][]} changes The changes to be applied
 */
export function applyObject(
  source: unknown,
  changes: [key: string, value?: unknown][],
): unknown {
  let result: any = source;
  let unchanged = true;

  changes.forEach((entry) => {
    const [key, value] = entry;

    if (entry.length === 1) {
      if (!Object.hasOwn(result, key)) return;
      if (unchanged) {
        const { [key]: _, ...rest } = result;
        result = rest;
        unchanged = false;
      } else {
        delete result[key];
      }
      return;
    }

    if (Object.hasOwn(result, key) && result[key] === value) return;
    if (unchanged) {
      result = { ...result, [key]: value };
      unchanged = false;
    } else {
      result[key] = value;
    }
  });

  return result;
}

/**
 * Apply all changes immutably to the source array.
 * Indices can be negative to be applied from the end of the array.
 *
 * @param {[index: number, value: any][]} changes The changes to be applied
 * @param filler Value to be used as padding
 */
export function applyArray(
  source: unknown[],
  changes: [index: number, value: unknown][],
  filler: unknown = null,
): unknown[] {
  let result = source;
  let unchanged = true;

  changes.forEach((entry) => {
    const [index, value] = entry;
    let i = index >= 0 ? index : result.length + index;

    if (i >= 0 && i < result.length && result[i] === value) return;
    if (index >= 0) {
      const suf = i + 1 - result.length;
      if (suf > 0) {
        result = [...result, ...new Array(suf).fill(filler)];
        unchanged = false;
      }
    } else {
      const pre = -i;
      if (pre > 0) {
        result = [...new Array(pre).fill(filler), ...result];
        unchanged = false;
        i = 0;
      }
    }
    if (unchanged) {
      result = [...result];
      unchanged = false;
    }
    result[i] = value;
  });

  return result;
}

/**
 * Create all possible combinations immutably based on the specified source array.
 * If the specified `combinations` array is empty, the resulting array contains a single empty array.
 * This function has essentially the same base functionality as the `mux` function, but uses a more performant approach for generating immutable arrays as it reduces the number of necessary array copies.
 *
 * `applyCombinations([], [[1, 2], [3, 4]])` for example produces:
 * - `[1, 3]`
 * - `[1, 4]`
 * - `[2, 3]`
 * - `[2, 4]`
 *
 * If the values of `source` are equal to the values of one of the combinations, it is used instead of a copy in the output array, e.g.:
 * `let i = [1, 2]; applyCombinations(i, [[1], [2]])[0] == i`
 * - `true`
 *
 * @param {any[][]} combinations The combinations
 */
export function applyCombinations<T>(source: T[], combinations: T[][]): T[][] {
  const l = combinations.length;
  let total = 1;
  const indices = combinations.map((entry): [number, number, T[]] => {
    const count = entry.length;
    total *= count;
    return [count, 0, entry];
  });
  if (total === 0) return [];
  let s = source;
  const sl = source.length;
  if (sl !== l) s = sl > l ? s.slice(0, l) : [...s, ...new Array(l - sl)];
  const out = new Array<T[]>(total);
  let c = 0;
  for (;;) {
    let result = s;
    let unchanged = true;
    indices.forEach(([, j, values], i) => {
      const v = values[j];
      if (result[i] === v) return;
      if (unchanged) {
        result = [...result];
        unchanged = false;
      }
      result[i] = v;
    });
    out[c] = result;
    c += 1;
    if (c >= total) break;
    for (let n = l - 1; n >= 0; n -= 1) {
      const i = indices[n];
      const [max, current] = i;
      const next = current + 1;
      if (next < max) {
        i[1] = next;
        break;
      }
      i[1] = 0;
    }
  }
  return out;
}
