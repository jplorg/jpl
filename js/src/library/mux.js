/**
 * Multiplex the specified array of arguments by calling cb for all possible combinations of arguments.
 * If the specified `args` array is empty, `cb` is called once without any arguments.
 *
 * `mux([[1, 2], [3, 4]], cb)` for example calls:
 * - `cb(1, 3)`
 * - `cb(1, 4)`
 * - `cb(2, 3)`
 * - `cb(2, 4)`
 */
function mux(args, cb) {
  const argCount = args.length;
  if (argCount === 1) {
    args[0].forEach((arg) => cb(arg));
    return;
  }
  let execCount = 1;
  const indices = args.map((arg) => {
    const argLength = arg?.length ?? 0;
    execCount *= argLength;
    return [argLength, 0, arg];
  });
  if (execCount === 0) return;
  const buffer = args.map(([v]) => v);
  let execIndex = 0;
  for (;;) {
    cb(...buffer);
    execIndex += 1;
    if (execIndex >= execCount) break;
    // determine next combination
    for (let argIndex = argCount - 1; argIndex >= 0; argIndex -= 1) {
      const arg = indices[argIndex];
      const [max, current, values] = arg;
      const next = current + 1;
      if (next < max) {
        arg[1] = next;
        buffer[argIndex] = values[next];
        break;
      }
      arg[1] = 0;
      [buffer[argIndex]] = values;
    }
  }
}

export default mux;

/** Multiplex the specified array of arguments and return the results produced by the callbacks */
export function muxOne(args, cb) {
  const argCount = args.length;
  if (argCount === 1) return args[0].map((arg) => cb(arg));
  let execCount = 1;
  const indices = args.map((arg) => {
    const argLength = arg?.length ?? 0;
    execCount *= argLength;
    return [argLength, 0, arg];
  });
  if (execCount === 0) return [];
  const outputs = new Array(execCount);
  const buffer = args.map(([v]) => v);
  let execIndex = 0;
  for (;;) {
    outputs[execIndex] = cb(...buffer);
    execIndex += 1;
    if (execIndex >= execCount) break;
    // determine next combination
    for (let argIndex = argCount - 1; argIndex >= 0; argIndex -= 1) {
      const arg = indices[argIndex];
      const [max, current, values] = arg;
      const next = current + 1;
      if (next < max) {
        arg[1] = next;
        buffer[argIndex] = values[next];
        break;
      }
      arg[1] = 0;
      [buffer[argIndex]] = values;
    }
  }
  return outputs;
}

/** Multiplex the specified array of arguments asynchronously and return the results produced by the callbacks */
export function muxAsync(args, cb) {
  return Promise.all(muxOne(args, cb));
}

/** Multiplex the specified array of arguments asynchronously and return a single array of all merged result arrays produced by the callbacks */
export async function muxAll(args, cb) {
  return mergeSegments(await muxAsync(args, cb));
}

/** Create a single array from the specified array segments */
export function mergeSegments(segments) {
  if (segments.length === 0) return segments;
  if (segments.length === 1) return segments[0];
  return segments.flat(1);
}
