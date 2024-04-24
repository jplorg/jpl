import {
  OPA_FIELD,
  OPA_FUNCTION,
  OPA_ITER,
  OPA_SLICE,
  OPC_EQUAL,
  OPC_GREATER,
  OPC_GREATEREQUAL,
  OPC_LESS,
  OPC_LESSEQUAL,
  OPC_UNEQUAL,
  OPM_ADDITION,
  OPM_DIVISION,
  OPM_MULTIPLICATION,
  OPM_REMAINDER,
  OPM_SUBTRACTION,
  OPU_ADDITION,
  OPU_DIVISION,
  OPU_MULTIPLICATION,
  OPU_NULL_COALESCENCE,
  OPU_REMAINDER,
  OPU_SET,
  OPU_SUBTRACTION,
  OPU_UPDATE,
  OP_ACCESS,
  OP_AND,
  OP_ARRAY_CONSTRUCTOR,
  OP_ASSIGNMENT,
  OP_CALCULATION,
  OP_COMPARISON,
  OP_CONSTANT_FALSE,
  OP_CONSTANT_NULL,
  OP_CONSTANT_TRUE,
  OP_FUNCTION_DEFINITION,
  OP_IF,
  OP_INTERPOLATED_STRING,
  OP_NEGATION,
  OP_NOT,
  OP_NULL_COALESCENCE,
  OP_NUMBER,
  OP_OBJECT_CONSTRUCTOR,
  OP_OR,
  OP_OUTPUT_CONCAT,
  OP_STRING,
  OP_TRY,
  OP_VARIABLE,
  OP_VARIABLE_DEFINITION,
  OP_VOID,
} from '../library';
import {
  eot,
  errorUnexpectedToken,
  hex,
  match,
  matchSet,
  matchWord,
  safeVariable,
  walkWhitespace,
} from './util';

/**
 * Parse a single program at i.
 * Throws an error if src contains additional content.
 */
export async function parseEntrypoint(src, i, c) {
  let n = i;

  const result = await parseProgram(src, n, c);
  ({ i: n } = result);

  if (!eot(src, n, c).is)
    return errorUnexpectedToken(src, n, c, { operator: 'program', message: 'expected EOT' });

  return result;
}

/** Parse program at i */
export function parseProgram(src, i, c) {
  let n = i;

  ({ i: n } = walkWhitespace(src, n, c));

  return opPipe(src, n, c);
}

/** Parse pipe at i */
export async function opPipe(src, i, c) {
  // Call stack decoupling - This is necessary as some browsers (i.e. Safari) have very limited call stack sizes which result in stack overflow exceptions in certain situations.
  await undefined;

  let n = i;

  const pipe = [];
  for (;;) {
    let ops;
    ({ i: n, ops } = await opOutputConcat(src, n, c));
    pipe.push(...ops);

    const m = matchWord(src, n, c, { phrase: '|', notBeforeSet: '=' });
    if (!m.is) break;
    ({ i: n } = m);
  }

  return { i: n, ops: pipe };
}

/** Parse subpipe at i */
export async function opSubPipe(src, i, c) {
  // Call stack decoupling - This is necessary as some browsers (i.e. Safari) have very limited call stack sizes which result in stack overflow exceptions in certain situations.
  await undefined;

  let n = i;

  const pipe = [];
  for (;;) {
    let ops;
    ({ i: n, ops } = await opTry(src, n, c));
    pipe.push(...ops);

    const m = matchWord(src, n, c, { phrase: '|', notBeforeSet: '=' });
    if (!m.is) break;
    ({ i: n } = m);
  }

  return { i: n, ops: pipe };
}

/** Parse subroute at i */
export async function opSubRoute(src, i, c) {
  // Call stack decoupling - This is necessary as some browsers (i.e. Safari) have very limited call stack sizes which result in stack overflow exceptions in certain situations.
  await undefined;

  return opTry(src, i, c);
}

/** Parse output concat at i */
export async function opOutputConcat(src, i, c) {
  let n = i;

  const pipes = [];
  for (;;) {
    let ops;
    ({ i: n, ops } = await opTry(src, n, c));
    pipes.push(ops);

    const m = matchWord(src, n, c, { phrase: ',' });
    if (!m.is) break;
    ({ i: n } = m);
  }

  if (pipes.length === 1) return { i: n, ops: pipes[0] };

  return { i: n, ops: [{ op: OP_OUTPUT_CONCAT, params: { pipes } }] };
}

/** Parse try at i */
export async function opTry(src, i, c) {
  let n = i;

  let m = matchWord(src, n, c, { phrase: 'try', spaceAfter: true });
  if (!m.is) return opOr(src, n, c);
  ({ i: n } = m);

  let opsTry;
  ({ i: n, ops: opsTry } = await opOr(src, n, c));

  let opsCatch;
  m = matchWord(src, n, c, { spaceBefore: true, phrase: 'catch', spaceAfter: true });
  if (m.is) {
    ({ i: n } = m);

    ({ i: n, ops: opsCatch } = await opOr(src, n, c));
  } else opsCatch = [{ op: OP_VOID }];

  return { i: n, ops: [{ op: OP_TRY, params: { try: opsTry, catch: opsCatch } }] };
}

/** Parse or at i */
export async function opOr(src, i, c) {
  let n = i;

  const pipes = [];
  for (;;) {
    let ops;
    ({ i: n, ops } = await opAnd(src, n, c));
    pipes.push(ops);

    const m = matchWord(src, n, c, { spaceBefore: true, phrase: 'or', spaceAfter: true });
    if (!m.is) break;
    ({ i: n } = m);
  }

  if (pipes.length === 1) return { i: n, ops: pipes[0] };

  return { i: n, ops: [{ op: OP_OR, params: { pipes } }] };
}

/** Parse and at i */
export async function opAnd(src, i, c) {
  let n = i;

  const pipes = [];
  for (;;) {
    let ops;
    ({ i: n, ops } = await opEquality(src, n, c));
    pipes.push(ops);

    const m = matchWord(src, n, c, { spaceBefore: true, phrase: 'and', spaceAfter: true });
    if (!m.is) break;
    ({ i: n } = m);
  }

  if (pipes.length === 1) return { i: n, ops: pipes[0] };

  return { i: n, ops: [{ op: OP_AND, params: { pipes } }] };
}

/** Parse equality at i */
export async function opEquality(src, i, c) {
  let n = i;

  let ops;
  ({ i: n, ops } = await opComparison(src, n, c));

  const comparisons = [];

  for (;;) {
    let m = matchWord(src, n, c, { phrase: '==' });
    if (m.is) {
      ({ i: n } = m);

      let opsBy;
      ({ i: n, ops: opsBy } = await opComparison(src, n, c));

      comparisons.push({ op: OPC_EQUAL, params: { by: opsBy } });
      continue;
    }

    m = matchWord(src, n, c, { phrase: '!=' });
    if (m.is) {
      ({ i: n } = m);

      let opsBy;
      ({ i: n, ops: opsBy } = await opComparison(src, n, c));

      comparisons.push({ op: OPC_UNEQUAL, params: { by: opsBy } });
      continue;
    }

    break;
  }

  if (comparisons.length === 0) return { i: n, ops };

  return { i: n, ops: [{ op: OP_COMPARISON, params: { pipe: ops, comparisons } }] };
}

/** Parse comparison at i */
export async function opComparison(src, i, c) {
  let n = i;

  let ops;
  ({ i: n, ops } = await opNot(src, n, c));

  const comparisons = [];

  for (;;) {
    let m = matchWord(src, n, c, { phrase: '<=' });
    if (m.is) {
      ({ i: n } = m);

      let opsBy;
      ({ i: n, ops: opsBy } = await opNot(src, n, c));

      comparisons.push({ op: OPC_LESSEQUAL, params: { by: opsBy } });
      continue;
    }

    m = matchWord(src, n, c, { phrase: '<' });
    if (m.is) {
      ({ i: n } = m);

      let opsBy;
      ({ i: n, ops: opsBy } = await opNot(src, n, c));

      comparisons.push({ op: OPC_LESS, params: { by: opsBy } });
      continue;
    }

    m = matchWord(src, n, c, { phrase: '>=' });
    if (m.is) {
      ({ i: n } = m);

      let opsBy;
      ({ i: n, ops: opsBy } = await opNot(src, n, c));

      comparisons.push({ op: OPC_GREATEREQUAL, params: { by: opsBy } });
      continue;
    }

    m = matchWord(src, n, c, { phrase: '>' });
    if (m.is) {
      ({ i: n } = m);

      let opsBy;
      ({ i: n, ops: opsBy } = await opNot(src, n, c));

      comparisons.push({ op: OPC_GREATER, params: { by: opsBy } });
      continue;
    }

    break;
  }

  if (comparisons.length === 0) return { i: n, ops };

  return { i: n, ops: [{ op: OP_COMPARISON, params: { pipe: ops, comparisons } }] };
}

/** Parse not at i */
export async function opNot(src, i, c) {
  let n = i;

  const m = matchWord(src, n, c, { phrase: 'not', spaceAfter: true });
  if (!m.is) return opErrorSuppression(src, n, c);
  ({ i: n } = m);

  let ops;
  ({ i: n, ops } = await opErrorSuppression(src, n, c));

  return { i: n, ops: [...ops, { op: OP_NOT }] };
}

/** Parse error suppression at i */
export async function opErrorSuppression(src, i, c) {
  let n = i;

  const result = await opDifference(src, n, c);
  ({ i: n } = result);

  const m = matchWord(src, n, c, { phrase: '?', notBeforeSet: '?=' });
  if (!m.is) return result;
  ({ i: n } = m);

  return {
    i: n,
    ops: [{ op: OP_TRY, params: { try: result.ops, catch: [{ op: OP_VOID }] } }],
  };
}

/** Parse difference at i */
export async function opDifference(src, i, c) {
  let n = i;

  let ops;
  ({ i: n, ops } = await opMultiplication(src, n, c));

  const operations = [];

  for (;;) {
    let m = matchWord(src, n, c, { phrase: '+', notBeforeSet: '=' });
    if (m.is) {
      ({ i: n } = m);

      let opsBy;
      ({ i: n, ops: opsBy } = await opMultiplication(src, n, c));

      operations.push({ op: OPM_ADDITION, params: { by: opsBy } });
      continue;
    }

    m = matchWord(src, n, c, { phrase: '-', notBeforeSet: '=>' });
    if (m.is) {
      ({ i: n } = m);

      let opsBy;
      ({ i: n, ops: opsBy } = await opMultiplication(src, n, c));

      operations.push({ op: OPM_SUBTRACTION, params: { by: opsBy } });
      continue;
    }

    break;
  }

  if (operations.length === 0) return { i: n, ops };

  return { i: n, ops: [{ op: OP_CALCULATION, params: { pipe: ops, operations } }] };
}

/** Parse multiplication at i */
export async function opMultiplication(src, i, c) {
  let n = i;

  let ops;
  ({ i: n, ops } = await opNullCoalescence(src, n, c));

  const operations = [];

  for (;;) {
    let m = matchWord(src, n, c, { phrase: '*', notBeforeSet: '=' });
    if (m.is) {
      ({ i: n } = m);

      let opsBy;
      ({ i: n, ops: opsBy } = await opNullCoalescence(src, n, c));

      operations.push({ op: OPM_MULTIPLICATION, params: { by: opsBy } });
      continue;
    }

    m = matchWord(src, n, c, { phrase: '/', notBeforeSet: '=' });
    if (m.is) {
      ({ i: n } = m);

      let opsBy;
      ({ i: n, ops: opsBy } = await opNullCoalescence(src, n, c));

      operations.push({ op: OPM_DIVISION, params: { by: opsBy } });
      continue;
    }

    m = matchWord(src, n, c, { phrase: '%', notBeforeSet: '=' });
    if (m.is) {
      ({ i: n } = m);

      let opsBy;
      ({ i: n, ops: opsBy } = await opNullCoalescence(src, n, c));

      operations.push({ op: OPM_REMAINDER, params: { by: opsBy } });
      continue;
    }

    break;
  }

  if (operations.length === 0) return { i: n, ops };

  return { i: n, ops: [{ op: OP_CALCULATION, params: { pipe: ops, operations } }] };
}

/** Parse null coalescence at i */
export async function opNullCoalescence(src, i, c) {
  let n = i;

  const pipes = [];
  for (;;) {
    const result = await opNegation(src, n, c);
    let ops;
    ({ i: n, ops } = result);

    pipes.push(ops);

    const m = matchWord(src, n, c, { phrase: '??' });
    if (!m.is) break;
    ({ i: n } = m);
  }

  if (pipes.length === 1) return { i: n, ops: pipes[0] };

  return { i: n, ops: [{ op: OP_NULL_COALESCENCE, params: { pipes } }] };
}

/** Parse negation at i */
export async function opNegation(src, i, c) {
  let n = i;

  const m = matchWord(src, n, c, { phrase: '-', notBeforeSet: '=>' });
  if (!m.is) return opIf(src, n, c);
  ({ i: n } = m);

  let ops;
  ({ i: n, ops } = await opIf(src, n, c));

  return { i: n, ops: [...ops, { op: OP_NEGATION }] };
}

/** Parse if at i */
export async function opIf(src, i, c) {
  let n = i;

  let m = matchWord(src, n, c, { phrase: 'if', spaceAfter: true });
  if (!m.is) return opConstant(src, n, c);
  ({ i: n } = m);

  const ifs = [];
  for (;;) {
    let opsIf;
    ({ i: n, ops: opsIf } = await opPipe(src, n, c));

    m = matchWord(src, n, c, { spaceBefore: true, phrase: 'then', spaceAfter: true });
    ({ i: n } = m);
    if (!m.is)
      return errorUnexpectedToken(src, n, c, {
        operator: 'if statement',
        message: "expected 'then'",
      });

    let opsThen;
    ({ i: n, ops: opsThen } = await opPipe(src, n, c));

    ifs.push({ if: opsIf, then: opsThen });

    m = matchWord(src, n, c, { spaceBefore: true, phrase: 'elif', spaceAfter: true });
    if (!m.is) break;
    ({ i: n } = m);
  }

  let opsElse;
  m = matchWord(src, n, c, { spaceBefore: true, phrase: 'else', spaceAfter: true });
  if (m.is) {
    ({ i: n } = m);

    ({ i: n, ops: opsElse } = await opPipe(src, n, c));
  } else opsElse = [];

  m = matchWord(src, n, c, { spaceBefore: true, phrase: 'end' });
  ({ i: n } = m);
  if (!m.is)
    return errorUnexpectedToken(src, n, c, {
      operator: 'if statement',
      message: "expected 'end'",
    });

  return { i: n, ops: [{ op: OP_IF, params: { ifs, else: opsElse } }] };
}

/** Parse constant at i */
export function opConstant(src, i, c) {
  let n = i;

  let m = matchWord(src, n, c, { phrase: 'true', spaceAfter: true });
  if (m.is) {
    ({ i: n } = m);
    return { i: n, ops: [{ op: OP_CONSTANT_TRUE }] };
  }

  m = matchWord(src, n, c, { phrase: 'false', spaceAfter: true });
  if (m.is) {
    ({ i: n } = m);
    return { i: n, ops: [{ op: OP_CONSTANT_FALSE }] };
  }

  m = matchWord(src, n, c, { phrase: 'null', spaceAfter: true });
  if (m.is) {
    ({ i: n } = m);
    return { i: n, ops: [{ op: OP_CONSTANT_NULL }] };
  }

  return opNumber(src, n, c);
}

/** Parse number at i */
export async function opNumber(src, i, c) {
  let n = i;

  const result = await number(src, n, c);
  if (!result.is) return opNamedFunctionDefinition(src, n, c);
  ({ i: n } = result);

  return { i: n, ops: result.ops };
}

/** Parse named function definition at i */
export async function opNamedFunctionDefinition(src, i, c) {
  let n = i;

  const m = matchWord(src, n, c, { phrase: 'func', spaceAfter: true });
  if (!m.is) return opFunctionDefinition(src, n, c);
  ({ i: n } = m);

  const v = safeVariable(src, n, c);
  if (!v.is) return opFunctionDefinition(src, i, c);
  let name;
  ({ i: n, value: name } = v);

  let argNames;
  ({ i: n, argNames } = await functionHeader(src, n, c));

  let ops;
  ({ i: n, ops } = await opSubRoute(src, n, c));

  return {
    i: n,
    ops: [
      {
        op: OP_VARIABLE_DEFINITION,
        params: { name, pipe: [{ op: OP_FUNCTION_DEFINITION, params: { argNames, pipe: ops } }] },
      },
    ],
  };
}

/** Parse function definition at i */
export async function opFunctionDefinition(src, i, c) {
  let n = i;

  const m = matchWord(src, n, c, { phrase: 'func', spaceAfter: true });
  if (!m.is) return opVariableAccess(src, n, c);
  ({ i: n } = m);

  let argNames;
  ({ i: n, argNames } = await functionHeader(src, n, c));

  let ops;
  ({ i: n, ops } = await opSubRoute(src, n, c));

  return { i: n, ops: [{ op: OP_FUNCTION_DEFINITION, params: { argNames, pipe: ops } }] };
}

/** Parse variable definition at i */
export async function opVariableAccess(src, i, c) {
  let n = i;

  const v = safeVariable(src, n, c);
  if (!v.is) return opValueAccess(src, n, c);
  let name;
  ({ i: n, value: name } = v);

  let operations;
  let canAssign;
  const ac = await access(src, n, c);
  if (ac.is) ({ i: n, operations, canAssign } = ac);
  else {
    operations = [];
    canAssign = true;
  }

  if (!canAssign) {
    const ops = [{ op: OP_VARIABLE, params: { name } }];

    if (operations.length === 0) return { i: n, ops };
    return { i: n, ops: [{ op: OP_ACCESS, params: { pipe: ops, operations } }] };
  }

  const as = await assignment(src, n, c);
  if (!as.is) {
    const ops = [{ op: OP_VARIABLE, params: { name } }];

    if (operations.length === 0) return { i: n, ops };
    return { i: n, ops: [{ op: OP_ACCESS, params: { pipe: ops, operations } }] };
  }
  let opAssignment;
  ({ i: n, assignment: opAssignment } = as);

  if (operations.length === 0 && opAssignment.op === OPU_SET) {
    return {
      i: n,
      ops: [{ op: OP_VARIABLE_DEFINITION, params: { name, pipe: opAssignment.params.pipe } }],
    };
  }

  return {
    i: n,
    ops: [
      {
        op: OP_VARIABLE_DEFINITION,
        params: {
          name,
          pipe: [
            {
              op: OP_ASSIGNMENT,
              params: {
                pipe: [{ op: OP_VARIABLE, params: { name } }],
                operations,
                assignment: opAssignment,
              },
            },
          ],
        },
      },
    ],
  };
}

/** Parse variable access at i */
export async function opValueAccess(src, i, c) {
  let n = i;

  const operations = [];

  let ops;
  let m = matchWord(src, n, c, { phrase: '.' });
  if (!m.is) ({ i: n, ops } = await opObjectConstructor(src, n, c));
  else {
    ({ i: n } = m);
    ops = [];

    const v = safeVariable(src, n, c);
    if (v.is) {
      let name;
      ({ i: n, value: name } = v);

      let optional;
      m = matchWord(src, n, c, { phrase: '?', notBeforeSet: '?=' });
      if (m.is) ({ i: n, is: optional } = m);

      operations.push({
        op: OPA_FIELD,
        params: { pipe: [{ op: OP_STRING, params: { value: name } }], optional },
      });
    }
  }

  const ac = await access(src, n, c, { identity: ops.length === 0 && operations.length === 0 });
  let canAssign;
  if (ac.is) {
    ({ i: n, canAssign } = ac);
    operations.push(...ac.operations);
  } else canAssign = operations.length > 0;

  if (operations.length === 0) return { i: n, ops };

  if (!canAssign) return { i: n, ops: [{ op: OP_ACCESS, params: { pipe: ops, operations } }] };

  const as = await assignment(src, n, c);
  if (!as.is) return { i: n, ops: [{ op: OP_ACCESS, params: { pipe: ops, operations } }] };
  let opAssignment;
  ({ i: n, assignment: opAssignment } = as);

  return {
    i: n,
    ops: [{ op: OP_ASSIGNMENT, params: { pipe: ops, operations, assignment: opAssignment } }],
  };
}

/** Parse object constructor at i */
export async function opObjectConstructor(src, i, c) {
  let n = i;

  let m = matchWord(src, n, c, { phrase: '{' });
  if (!m.is) return opArrayConstructor(src, n, c);
  ({ i: n } = m);

  const fields = [];

  m = matchWord(src, n, c, { phrase: '}' });
  if (m.is) ({ i: n } = m);
  else
    for (;;) {
      m = matchWord(src, n, c, { phrase: '(' });
      if (m.is) {
        ({ i: n } = m);

        let opsKey;
        ({ i: n, ops: opsKey } = await opPipe(src, n, c));

        m = matchWord(src, n, c, { phrase: ')' });
        ({ i: n } = m);
        if (!m.is)
          return errorUnexpectedToken(src, n, c, { operator: 'object', message: "expected ')'" });

        let optional;
        ({ i: n, is: optional } = matchWord(src, n, c, { phrase: '?' }));

        m = matchWord(src, n, c, { phrase: ':' });
        ({ i: n } = m);
        if (!m.is)
          return errorUnexpectedToken(src, n, c, { operator: 'object', message: "expected ':'" });

        let opsValue;
        ({ i: n, ops: opsValue } = await opSubPipe(src, n, c));

        fields.push({ key: opsKey, value: opsValue, optional });

        m = matchWord(src, n, c, { phrase: '}' });
        if (m.is) {
          ({ i: n } = m);
          break;
        }

        m = matchWord(src, n, c, { phrase: ',' });
        ({ i: n } = m);
        if (!m.is)
          return errorUnexpectedToken(src, n, c, {
            operator: 'object',
            message: "expected ',' or '}'",
          });

        continue;
      }

      const s = await string(src, n, c);
      if (s.is) {
        let opsKey;
        ({ i: n, ops: opsKey } = s);

        m = matchWord(src, n, c, { phrase: ':' });
        ({ i: n } = m);
        if (!m.is)
          return errorUnexpectedToken(src, n, c, { operator: 'object', message: "expected ':'" });

        let opsValue;
        ({ i: n, ops: opsValue } = await opSubPipe(src, n, c));

        fields.push({ key: opsKey, value: opsValue, optional: false });

        m = matchWord(src, n, c, { phrase: '}' });
        if (m.is) {
          ({ i: n } = m);
          break;
        }

        m = matchWord(src, n, c, { phrase: ',' });
        ({ i: n } = m);
        if (!m.is)
          return errorUnexpectedToken(src, n, c, {
            operator: 'object',
            message: "expected ',' or '}'",
          });

        continue;
      }

      const v = safeVariable(src, n, c);
      if (v.is) {
        let name;
        ({ i: n, value: name } = v);

        let opsValue;
        m = matchWord(src, n, c, { phrase: ':' });
        if (m.is) {
          ({ i: n } = m);

          ({ i: n, ops: opsValue } = await opSubPipe(src, n, c));
        } else {
          let optional;
          ({ i: n, is: optional } = matchWord(src, n, c, { phrase: '?' }));

          opsValue = optional
            ? [
                {
                  op: OP_TRY,
                  params: {
                    try: [{ op: OP_VARIABLE, params: { name } }],
                    catch: [{ op: OP_VOID }],
                  },
                },
              ]
            : [{ op: OP_VARIABLE, params: { name } }];
        }

        fields.push({
          key: [{ op: OP_STRING, params: { value: name } }],
          value: opsValue,
          optional: false,
        });

        m = matchWord(src, n, c, { phrase: '}' });
        if (m.is) {
          ({ i: n } = m);
          break;
        }

        m = matchWord(src, n, c, { phrase: ',' });
        ({ i: n } = m);
        if (!m.is)
          return errorUnexpectedToken(src, n, c, {
            operator: 'object',
            message: "expected ',' or '}'",
          });

        continue;
      }

      return errorUnexpectedToken(src, n, c, {
        operator: 'object',
        message: 'expected field declaration',
      });
    }

  return { i: n, ops: [{ op: OP_OBJECT_CONSTRUCTOR, params: { fields } }] };
}

/** Parse array constructor at i */
export async function opArrayConstructor(src, i, c) {
  let n = i;

  let m = matchWord(src, n, c, { phrase: '[' });
  ({ i: n } = m);
  if (!m.is) return opStringLiteral(src, n, c);

  let ops;

  m = matchWord(src, n, c, { phrase: ']' });
  if (m.is) {
    ({ i: n } = m);

    ops = [{ op: OP_VOID }];
  } else {
    ({ i: n, ops } = await opPipe(src, n, c));

    m = matchWord(src, n, c, { phrase: ']' });
    ({ i: n } = m);
    if (!m.is)
      return errorUnexpectedToken(src, n, c, { operator: 'array', message: "expected ']'" });
  }

  return { i: n, ops: [{ op: OP_ARRAY_CONSTRUCTOR, params: { pipe: ops } }] };
}

/** Parse string literal at i */
export async function opStringLiteral(src, i, c) {
  let n = i;

  const result = await string(src, n, c);
  if (!result.is) return opGroup(src, n, c);
  ({ i: n } = result);

  return { i: n, ops: result.ops };
}

/** Parse group at i */
export async function opGroup(src, i, c) {
  let n = i;

  let m = matchWord(src, n, c, { phrase: '(' });
  ({ i: n } = m);
  if (!m.is) return errorUnexpectedToken(src, n, c);

  let ops;
  ({ i: n, ops } = await opPipe(src, n, c));

  m = matchWord(src, n, c, { phrase: ')' });
  ({ i: n } = m);
  if (!m.is) return errorUnexpectedToken(src, n, c, { operator: 'group', message: "expected ')'" });

  return { i: n, ops };
}

/** Parse function header at i */
export function functionHeader(src, i, c) {
  let n = i;

  let m = matchWord(src, n, c, { phrase: '(' });
  ({ i: n } = m);
  if (!m.is)
    return errorUnexpectedToken(src, n, c, {
      operator: 'function definition',
      message: "expected '('",
    });

  const argNames = [];

  m = matchWord(src, n, c, { phrase: ')' });
  if (m.is) ({ i: n } = m);
  else
    for (;;) {
      const v = safeVariable(src, n, c);
      if (!v.is)
        return errorUnexpectedToken(src, n, c, {
          operator: 'function definition',
          message: 'expected argument name',
        });
      let name;
      ({ i: n, value: name } = v);

      argNames.push(name);

      m = matchWord(src, n, c, { phrase: ')' });
      if (m.is) {
        ({ i: n } = m);
        break;
      }

      m = matchWord(src, n, c, { phrase: ',' });
      ({ i: n } = m);
      if (!m.is)
        return errorUnexpectedToken(src, n, c, {
          operator: 'function definition',
          message: "expected ',' or ')'",
        });
    }

  m = matchWord(src, n, c, { phrase: ':' });
  ({ i: n } = m);
  if (!m.is)
    return errorUnexpectedToken(src, n, c, {
      operator: 'function definition',
      message: "expected ':'",
    });

  return { i: n, argNames };
}

/** Parse access at i */
export async function access(src, i, c, { identity } = {}) {
  let n = i;

  const operations = [];
  let canAssign = true;

  for (;;) {
    let m = matchWord(src, n, c, { phrase: '.' });
    const isIdentity = identity && operations.length === 0;
    if (!isIdentity && m.is) {
      ({ i: n } = m);

      const v = safeVariable(src, n, c);
      if (!v.is)
        return errorUnexpectedToken(src, n, c, {
          operator: 'field access operator',
          message: 'expected field name',
        });
      let name;
      ({ i: n, value: name } = v);

      let optional;
      m = matchWord(src, n, c, { phrase: '?', notBeforeSet: '?=' });
      if (m.is) ({ i: n, is: optional } = m);

      operations.push({
        op: OPA_FIELD,
        params: { pipe: [{ op: OP_STRING, params: { value: name } }], optional },
      });
      continue;
    }

    m = matchWord(src, n, c, { phrase: '[' });
    if (m.is) {
      ({ i: n } = m);

      m = matchWord(src, n, c, { phrase: ']' });
      if (m.is) {
        ({ i: n } = m);

        let optional;
        m = matchWord(src, n, c, { phrase: '?', notBeforeSet: '?=' });
        if (m.is) ({ i: n, is: optional } = m);

        operations.push({ op: OPA_ITER, params: { optional } });
        continue;
      }

      m = matchWord(src, n, c, { phrase: ':' });
      if (m.is) {
        ({ i: n } = m);

        m = matchWord(src, n, c, { phrase: ']' });
        if (m.is) {
          ({ i: n } = m);

          let optional;
          m = matchWord(src, n, c, { phrase: '?', notBeforeSet: '?=' });
          if (m.is) ({ i: n, is: optional } = m);

          operations.push({
            op: OPA_SLICE,
            params: { from: [{ op: OP_CONSTANT_NULL }], to: [{ op: OP_CONSTANT_NULL }], optional },
          });
          continue;
        }

        let opsRight;
        ({ i: n, ops: opsRight } = await opPipe(src, n, c));

        m = matchWord(src, n, c, { phrase: ']' });
        ({ i: n } = m);
        if (!m.is)
          return errorUnexpectedToken(src, n, c, {
            operator: 'array slice operator',
            message: "expected ']'",
          });

        let optional;
        m = matchWord(src, n, c, { phrase: '?', notBeforeSet: '?=' });
        if (m.is) ({ i: n, is: optional } = m);

        operations.push({
          op: OPA_SLICE,
          params: { from: [{ op: OP_CONSTANT_NULL }], to: opsRight, optional },
        });
        continue;
      }

      let opsLeft;
      ({ i: n, ops: opsLeft } = await opPipe(src, n, c));

      m = matchWord(src, n, c, { phrase: ']' });
      if (m.is) {
        ({ i: n } = m);

        let optional;
        m = matchWord(src, n, c, { phrase: '?', notBeforeSet: '?=' });
        if (m.is) ({ i: n, is: optional } = m);

        operations.push({ op: OPA_FIELD, params: { pipe: opsLeft, optional } });
        continue;
      }

      m = matchWord(src, n, c, { phrase: ':' });
      if (!m.is)
        return errorUnexpectedToken(src, n, c, {
          operator: 'variable access operator',
          message: "expected ':' or ']'",
        });
      ({ i: n } = m);

      m = matchWord(src, n, c, { phrase: ']' });
      if (m.is) {
        ({ i: n } = m);

        let optional;
        m = matchWord(src, n, c, { phrase: '?', notBeforeSet: '?=' });
        if (m.is) ({ i: n, is: optional } = m);

        operations.push({
          op: OPA_SLICE,
          params: { from: opsLeft, to: [{ op: OP_CONSTANT_NULL }], optional },
        });
        continue;
      }

      let opsRight;
      ({ i: n, ops: opsRight } = await opPipe(src, n, c));

      m = matchWord(src, n, c, { phrase: ']' });
      ({ i: n } = m);
      if (!m.is)
        return errorUnexpectedToken(src, n, c, {
          operator: 'array slice operator',
          message: "expected ']'",
        });

      let optional;
      m = matchWord(src, n, c, { phrase: '?', notBeforeSet: '?=' });
      if (m.is) ({ i: n, is: optional } = m);

      operations.push({ op: OPA_SLICE, params: { from: opsLeft, to: opsRight, optional } });
      continue;
    }

    let bound;
    m = matchWord(src, n, c, { phrase: '->' });
    if (m.is) ({ i: n, is: bound } = m);

    m = matchWord(src, n, c, { phrase: '(' });
    if (!m.is && bound) {
      return errorUnexpectedToken(src, n, c, {
        operator: 'bound function call',
        message: "expected '('",
      });
    }
    if (m.is) {
      ({ i: n } = m);

      const args = [];

      m = matchWord(src, n, c, { phrase: ')' });
      if (m.is) ({ i: n } = m);
      else
        for (;;) {
          let opsArg;
          ({ i: n, ops: opsArg } = await opSubPipe(src, n, c));

          args.push(opsArg);

          m = matchWord(src, n, c, { phrase: ')' });
          if (m.is) {
            ({ i: n } = m);
            break;
          }

          m = matchWord(src, n, c, { phrase: ',' });
          ({ i: n } = m);
          if (!m.is)
            return errorUnexpectedToken(src, n, c, {
              operator: 'function call',
              message: "expected ',' or ')'",
            });
        }

      let optional;
      m = matchWord(src, n, c, { phrase: '?', notBeforeSet: '?=' });
      if (m.is) ({ i: n, is: optional } = m);

      operations.push({ op: OPA_FUNCTION, params: { args, bound, optional } });
      canAssign = false;
      continue;
    }

    break;
  }

  return { i: n, is: operations.length > 0, operations, canAssign };
}

/** Parse assignment at i */
export async function assignment(src, i, c) {
  let n = i;

  let m = matchWord(src, n, c, { phrase: '=', notBeforeSet: '=' });
  if (m.is) {
    ({ i: n } = m);

    let ops;
    ({ i: n, ops } = await opSubRoute(src, n, c));

    return { i: n, is: true, assignment: { op: OPU_SET, params: { pipe: ops } } };
  }

  m = matchWord(src, n, c, { phrase: '|=' });
  if (m.is) {
    ({ i: n } = m);

    let ops;
    ({ i: n, ops } = await opSubRoute(src, n, c));

    return { i: n, is: true, assignment: { op: OPU_UPDATE, params: { pipe: ops } } };
  }

  m = matchWord(src, n, c, { phrase: '+=' });
  if (m.is) {
    ({ i: n } = m);

    let ops;
    ({ i: n, ops } = await opSubRoute(src, n, c));

    return { i: n, is: true, assignment: { op: OPU_ADDITION, params: { pipe: ops } } };
  }

  m = matchWord(src, n, c, { phrase: '-=' });
  if (m.is) {
    ({ i: n } = m);

    let ops;
    ({ i: n, ops } = await opSubRoute(src, n, c));

    return { i: n, is: true, assignment: { op: OPU_SUBTRACTION, params: { pipe: ops } } };
  }

  m = matchWord(src, n, c, { phrase: '*=' });
  if (m.is) {
    ({ i: n } = m);

    let ops;
    ({ i: n, ops } = await opSubRoute(src, n, c));

    return { i: n, is: true, assignment: { op: OPU_MULTIPLICATION, params: { pipe: ops } } };
  }

  m = matchWord(src, n, c, { phrase: '/=' });
  if (m.is) {
    ({ i: n } = m);

    let ops;
    ({ i: n, ops } = await opSubRoute(src, n, c));

    return { i: n, is: true, assignment: { op: OPU_DIVISION, params: { pipe: ops } } };
  }

  m = matchWord(src, n, c, { phrase: '%=' });
  if (m.is) {
    ({ i: n } = m);

    let ops;
    ({ i: n, ops } = await opSubRoute(src, n, c));

    return { i: n, is: true, assignment: { op: OPU_REMAINDER, params: { pipe: ops } } };
  }

  m = matchWord(src, n, c, { phrase: '?=' });
  if (m.is) {
    ({ i: n } = m);

    let ops;
    ({ i: n, ops } = await opSubRoute(src, n, c));

    return { i: n, is: true, assignment: { op: OPU_NULL_COALESCENCE, params: { pipe: ops } } };
  }

  return { i: n, is: false };
}

/** Parse number at i */
export function number(src, i, c) {
  let n = i;
  let is = false;
  let value = '';

  for (;;) {
    const set = matchSet(src, n, c, { set: '0123456789' });
    if (!set.is) break;
    ({ i: n } = set);
    is = true;
    value += set.value;
  }

  if (!is) return { i: n, is: false };

  const m = match(src, n, c, { phrase: '.' });
  if (m.is) {
    ({ i: n } = m);
    value += '.';

    for (;;) {
      const set = matchSet(src, n, c, { set: '0123456789' });
      if (!set.is) break;
      ({ i: n } = set);
      value += set.value;
    }
  }

  let set = matchSet(src, n, c, { set: 'eE' });
  if (set.is) {
    ({ i: n } = set);
    value += set.value;

    set = matchSet(src, n, c, { set: '+-' });
    if (set.is) {
      ({ i: n } = set);
      value += set.value;
    }

    let isE;
    for (;;) {
      set = matchSet(src, n, c, { set: '0123456789' });
      if (!set.is) break;
      ({ i: n } = set);
      isE = true;
      value += set.value;
    }
    if (!isE)
      return errorUnexpectedToken(src, n, c, { operator: 'number', message: 'expected digit' });
  }

  ({ i: n } = walkWhitespace(src, n, c));

  return { i: n, is: true, ops: [{ op: OP_NUMBER, params: { value } }] };
}

/** Parse string at i */
export async function string(src, i, c) {
  let n = i;
  let value = '';

  const set = matchSet(src, n, c, { set: '"\'`' });
  if (!set.is) return { i: n, is: false };
  let boundary;
  ({ i: n, value: boundary } = set);

  const multilineString = boundary === '`';

  const interpolations = [];

  for (;;) {
    let m = match(src, n, c, { phrase: boundary });
    if (m.is) {
      ({ i: n } = m);
      break;
    }

    const end = eot(src, n, c);
    if (end.is) {
      ({ i: n } = end);
      return errorUnexpectedToken(src, n, c, {
        operator: 'string',
        message: 'incomplete string literal',
      });
    }

    m = match(src, n, c, { phrase: '\\' });
    if (m.is) {
      ({ i: n } = m);

      m = matchWord(src, n, c, { phrase: '(' });
      if (m.is) {
        ({ i: n } = m);

        let ops;
        ({ i: n, ops } = await opPipe(src, n, c));

        m = match(src, n, c, { phrase: ')' });
        ({ i: n } = m);
        if (!m.is)
          return errorUnexpectedToken(src, n, c, {
            operator: 'string interpolation',
            message: "expected ')'",
          });

        interpolations.push({ before: value, pipe: ops });
        value = '';
        continue;
      }

      if (multilineString) {
        switch (src[n]) {
          case '\n':
            n += 1;
            continue;

          case '\r':
            n += 1;
            if (!eot(src, n, c).is && src[n] === '\n') n += 1;
            continue;

          default:
        }
      }

      switch (src[n]) {
        case '"':
          value += '"';
          n += 1;
          continue;

        case "'":
          value += "'";
          n += 1;
          continue;

        case '`':
          value += '`';
          n += 1;
          continue;

        case '\\':
          value += '\\';
          n += 1;
          continue;

        case '/':
          value += '/';
          n += 1;
          continue;

        case 'b':
          value += '\b';
          n += 1;
          continue;

        case 'f':
          value += '\f';
          n += 1;
          continue;

        case 'n':
          value += '\n';
          n += 1;
          continue;

        case 'r':
          value += '\r';
          n += 1;
          continue;

        case 't':
          value += '\t';
          n += 1;
          continue;

        case 'u': {
          n += 1;
          let hexVal = '';
          for (let j = 0; j < 4; j += 1) {
            m = hex(src, n, c);
            ({ i: n } = m);
            if (!m.is)
              return errorUnexpectedToken(src, n, c, {
                operator: 'string',
                message: 'incomplete unicode escape sequence: expected hex digit',
              });
            hexVal += m.value;
          }
          value += String.fromCharCode(parseInt(hexVal, 16));
          continue;
        }

        default:
          return errorUnexpectedToken(src, n, c, {
            operator: 'string',
            message: 'invalid escape sequence',
          });
      }
    }

    if (multilineString) {
      switch (src[n]) {
        case '\n':
        case '\r':
        case '\t':
          value += src[n];
          n += 1;
          continue;

        default:
      }
    }

    if (src.charCodeAt(n) < 0x20) return errorUnexpectedToken(src, n, c, { operator: 'string' });

    value += src[n];
    n += 1;
  }

  ({ i: n } = walkWhitespace(src, n, c));

  if (interpolations.length === 0)
    return { i: n, is: true, ops: [{ op: OP_STRING, params: { value } }] };

  return {
    i: n,
    is: true,
    ops: [{ op: OP_INTERPOLATED_STRING, params: { interpolations, after: value } }],
  };
}
