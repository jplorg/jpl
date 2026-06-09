import type { Pipe } from '@/library/definition';
import type JPLProgram from '@/program';
import type { UserJPLInstructionParams } from '@/program/params';
import applyDefaults from '../applyDefaults';
import {
  JPLFatalError,
  JPLRuntimeScope,
  JPLType,
  adaptErrorsAsync,
  applyObject,
  assertType,
  mux,
  muxAll,
  muxAsync,
  muxOne,
  normalize,
  stringify,
  strip,
  typeOf,
  typeOrder,
  unwrap,
  type JPLRuntimeScopeConfig,
} from '../library';

export type JPLRuntimeConfig = {
  runtime?: JPLRuntimeOptions;
};

export type JPLRuntimeOptions = {
  vars?: { [name: string]: unknown };
  adjustResult?: (
    result: unknown,
    scope: JPLRuntimeScope,
  ) => Promise<unknown[]> | unknown[];
};

const defaultOptions: JPLRuntimeOptions = {
  vars: {},
};

export function applyRuntimeDefaults(
  options: JPLRuntimeOptions = {},
  defaults: JPLRuntimeOptions = {},
): JPLRuntimeOptions {
  return applyDefaults(options, defaults, 'vars');
}

/** JPL runtime */
export default class JPLRuntime {
  #options;
  #program;

  constructor(program: JPLProgram, options?: JPLRuntimeConfig) {
    this.#options = applyRuntimeDefaults(options?.runtime, defaultOptions);

    this.#program = program;
  }

  /** Return the runtime's options */
  get options() {
    return this.#options;
  }

  /** Return the runtime's program */
  get program() {
    return this.#program;
  }

  /** Create a new orphan scope */
  createScope = (presets?: JPLRuntimeScopeConfig): JPLRuntimeScope =>
    new JPLRuntimeScope(presets);

  /** Execute a new dedicated program */
  execute = async (inputs: unknown[]): Promise<unknown[]> => {
    const scope = this.createScope({
      vars: Object.fromEntries(
        this.muxOne(
          [Object.entries(this.options.vars!)],
          ([name, value]: [string, unknown]) => [
            name,
            this.normalizeValue(value),
          ],
        ),
      ),
    });

    try {
      return await this.executeInstructions(
        this.program.definition.instructions ?? [],
        inputs,
        scope,
        this.options.adjustResult,
      );
    } finally {
      scope.signal.exit();
    }
  };

  /** Execute the specified instructions */
  executeInstructions = (
    instructions: Pipe,
    inputs: unknown[],
    scope: JPLRuntimeScope,
    next: (
      output: unknown,
      scope: JPLRuntimeScope,
    ) => Promise<unknown[]> | unknown[] = (output) => [output],
  ): Promise<unknown[]> => {
    const iter = async (
      from: number,
      input: unknown,
      currentScope: JPLRuntimeScope,
    ): Promise<unknown[]> => {
      // Call stack decoupling - This is necessary as some browsers (i.e. Safari) have very limited call stack sizes which result in stack overflow exceptions in certain situations.
      await undefined;

      currentScope.signal.checkHealth();

      if (from >= instructions.length) return next(input, currentScope);

      const { op, params } = instructions[from];
      const operator = this.program.ops[op];
      if (!operator) throw new JPLFatalError(`invalid OP '${op}'`);

      return operator.op(
        this,
        input,
        params ?? {},
        currentScope,
        (output: unknown, nextScope: JPLRuntimeScope): Promise<unknown[]> =>
          iter(from + 1, output, nextScope),
      );
    };

    return this.muxAll([inputs], (input: unknown) => iter(0, input, scope));
  };

  /** Execute the specified OP */
  op(
    op: string,
    params: UserJPLInstructionParams,
    inputs: unknown[],
    scope: JPLRuntimeScope,
    next: (output: unknown, scope: JPLRuntimeScope) => unknown[] = (output) => [
      output,
    ],
  ): Promise<unknown[]> {
    const operator = this.program.ops[op];
    if (!operator) throw new JPLFatalError(`invalid OP '${op}'`);

    const opParams = operator.map(this, params);
    return this.muxAll([inputs], (input: unknown) =>
      operator.op(this, input, opParams ?? {}, scope, next),
    );
  }

  /** Normalize the specified external value */
  normalizeValue = normalize;

  /** Normalize the specified array of external values */
  normalizeValues = (
    values: unknown[] = [],
    name: string = 'values',
  ): unknown[] => {
    if (!Array.isArray(values))
      throw new JPLFatalError(`expected ${name} to be an array`);
    return this.normalizeValue(values) as unknown[];
  };

  /** Unwrap the specified normalized value for usage in JPL operations */
  unwrapValue = unwrap;

  /** Unwrap the specified array of normalized values for usage in JPL operations */
  unwrapValues = (
    values: unknown[] = [],
    name: string = 'values',
  ): unknown[] => {
    if (!Array.isArray(values))
      throw new JPLFatalError(`expected ${name} to be an array`);
    return this.muxOne([values], (value) => this.unwrapValue(value));
  };

  /** Strip the specified normalized value for usage in JPL operations */
  stripValue = (value: unknown): unknown =>
    strip(value, (_k, v) => this.unwrapValue(v));

  /** Strip the specified array of normalized values for usage in JPL operations */
  stripValues = (
    values: unknown[] = [],
    name: string = 'values',
  ): unknown[] => {
    if (!Array.isArray(values))
      throw new JPLFatalError(`expected ${name} to be an array`);
    return this.muxOne([values], (value) => this.stripValue(value));
  };

  /** Alter the specified normalized value using the specified updater */
  alterValue = async (
    value: unknown,
    updater: (value: unknown) => Promise<unknown> | unknown,
  ): Promise<unknown> =>
    JPLType.is(value)
      ? adaptErrorsAsync(() => (value as JPLType).alter(updater))
      : this.normalizeValue(await updater(value));

  /** Resolve the type of the specified normalized value for JPL operations */
  type = typeOf;

  /** Assert the type for the specified unwrapped value for JPL operations */
  assertType = assertType;

  /** Determine whether the specified normalized value should be considered as truthy in JPL operations */
  truthy = (value: unknown): boolean => {
    const raw = this.unwrapValue(value);
    return raw !== null && raw !== false;
  };

  /** Compare the specified normalized values */
  compare = compare.bind(this);

  /** Compare the specified normalized strings based on their unicode code points */
  compareStrings = (a: unknown, b: unknown): number => {
    const ta = this.type(a);
    if (ta !== 'string') throw new JPLFatalError(`unexpected type ${ta}`);
    const tb = this.type(b);
    if (tb !== 'string') throw new JPLFatalError(`unexpected type ${tb}`);
    return compareStrings(
      this.unwrapValue(a) as string,
      this.unwrapValue(b) as string,
    );
  };

  /** Compare the specified normalized arrays based on their lexical order */
  compareArrays = (a: unknown, b: unknown): number => {
    const ta = this.type(a);
    if (ta !== 'array') throw new JPLFatalError(`unexpected type ${ta}`);
    const tb = this.type(b);
    if (tb !== 'array') throw new JPLFatalError(`unexpected type ${tb}`);

    return compareArrays.call(
      this,
      this.unwrapValue(a) as unknown[],
      this.unwrapValue(b) as unknown[],
    );
  };

  /** Compare the specified normalized objects */
  compareObjects = (a: unknown, b: unknown): number => {
    const ta = this.type(a);
    if (ta !== 'object') throw new JPLFatalError(`unexpected type ${ta}`);
    const tb = this.type(b);
    if (tb !== 'object') throw new JPLFatalError(`unexpected type ${tb}`);

    return compareObjects.call(
      this,
      this.unwrapValue(a) as Record<string, unknown>,
      this.unwrapValue(b) as Record<string, unknown>,
    );
  };

  /** Determine if the specified normalized values can be considered to be equal */
  equals = (a: unknown, b: unknown) => this.compare(a, b) === 0;

  /** Deep merge the specified normalized values */
  merge = async (a: unknown, b: unknown): Promise<unknown> => {
    // Call stack decoupling - This is necessary as some browsers (i.e. Safari) have very limited call stack sizes which result in stack overflow exceptions in certain situations.
    await undefined;

    if (this.type(a) !== 'object' || this.type(b) !== 'object') return b;

    return this.alterValue(a, async (value) => {
      const changes = await Promise.all(
        Object.entries(this.unwrapValue(b) as Record<string, unknown>).map(
          async ([k, v]): Promise<[string, unknown]> => [
            k,
            await this.merge((value as Record<string, unknown>)[k] ?? null, v),
          ],
        ),
      );

      return applyObject(value, changes);
    });
  };

  /** Stringify the specified normalized value for usage in program outputs */
  stringifyJSON = (value: unknown, unescapeString: boolean): string =>
    stringify(value, unescapeString);

  /** Strip the specified normalized value for usage in program outputs */
  stripJSON = (value: unknown): unknown => strip(value);

  /**
   * Multiplex the specified array of arguments by calling cb for all possible combinations of arguments.
   *
   * `mux([[1,2], [3,4]], cb)` for example yields:
   * - `cb(1, 3)`
   * - `cb(1, 4)`
   * - `cb(2, 3)`
   * - `cb(2, 4)`
   */
  mux = mux;

  /** Multiplex the specified array of arguments and return the results produced by the callbacks */
  muxOne = muxOne;

  /** Multiplex the specified array of arguments asynchronously and return the results produced by the callbacks */
  muxAsync = muxAsync;

  /** Multiplex the specified array of arguments asynchronously and return a single array of all merged result arrays produced by the callbacks */
  muxAll = muxAll;
}

/** Compare the specified normalized values */
function compare(this: JPLRuntime, a: unknown, b: unknown): number {
  const ta = this.type(a);
  const tb = this.type(b);

  if (ta !== tb) return typeOrder.indexOf(ta) - typeOrder.indexOf(tb);

  const ua = this.unwrapValue(a);
  const ub = this.unwrapValue(b);

  switch (ta) {
    case 'null':
    case 'function':
      return 0;

    case 'boolean':
    case 'number':
      return +(ua as number) - +(ub as number);

    case 'string':
      return compareStrings(ua as string, ub as string);

    case 'array':
      return compareArrays.call(this, ua as unknown[], ub as unknown[]);

    case 'object':
      return compareObjects.call(
        this,
        ua as Record<string, unknown>,
        ub as Record<string, unknown>,
      );

    default:
      throw new JPLFatalError(`unexpected type ${ta}`);
  }
}

/** Compare the specified normalized strings based on their unicode code points */
function compareStrings(a: string, b: string): number {
  const min = Math.min(a.length, b.length);
  let i = 0;

  for (const _ of a) {
    if (i >= min) {
      break;
    }
    const cp1 = a.codePointAt(i)!;
    const cp2 = b.codePointAt(i)!;
    const order = cp1 - cp2;
    if (order !== 0) {
      return order;
    }
    i += 1;
    if (cp1 > 0xffff) {
      i += 1;
    }
  }
  return a.length - b.length;
}

/** Compare the specified normalized arrays based on their lexical order */
function compareArrays(this: JPLRuntime, a: unknown[], b: unknown[]): number {
  const min = Math.min(a.length, b.length);
  for (let i = 0; i < min; i += 1) {
    const c = compare.call(this, a[i], b[i]);
    if (c !== 0) return c;
  }
  return a.length - b.length;
}

/** Compare the specified normalized objects */
function compareObjects(
  this: JPLRuntime,
  a: Record<string, unknown>,
  b: Record<string, unknown>,
): number {
  const aKeys = Object.keys(a).sort(compareStrings);
  const bKeys = Object.keys(b).sort(compareStrings);
  let order = compareArrays.call(this, aKeys, bKeys);
  if (order !== 0) return order;
  for (const key of aKeys) {
    order = compare.call(this, a[key], b[key]);
    if (order !== 0) return order;
  }
  return 0;
}
