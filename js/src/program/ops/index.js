import {
  OP_ACCESS,
  OP_AND,
  OP_ARRAY_CONSTRUCTOR,
  OP_ASSIGNMENT,
  OP_CALCULATION,
  OP_COMPARISON,
  OP_CONSTANT,
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
} from '../../library';
import opAccess from './opAccess';
import opAnd from './opAnd';
import opArrayConstructor from './opArrayConstructor';
import opAssignment from './opAssignment';
import opCalculation from './opCalculation';
import opComparison from './opComparison';
import opConstant from './opConstant';
import opConstantFalse from './opConstantFalse';
import opConstantNull from './opConstantNull';
import opConstantTrue from './opConstantTrue';
import opFunctionDefinition from './opFunctionDefinition';
import opIf from './opIf';
import opInterpolatedString from './opInterpolatedString';
import opNegation from './opNegation';
import opNot from './opNot';
import opNullCoalescence from './opNullCoalescence';
import opNumber from './opNumber';
import opObjectConstructor from './opObjectConstructor';
import opOr from './opOr';
import opOutputConcat from './opOutputConcat';
import opString from './opString';
import opTry from './opTry';
import opVariable from './opVariable';
import opVariableDefinition from './opVariableDefinition';
import opVoid from './opVoid';

const ops = {
  [OP_ACCESS]: opAccess,
  [OP_AND]: opAnd,
  [OP_ARRAY_CONSTRUCTOR]: opArrayConstructor,
  [OP_ASSIGNMENT]: opAssignment,
  [OP_CALCULATION]: opCalculation,
  [OP_COMPARISON]: opComparison,
  [OP_CONSTANT]: opConstant,
  [OP_CONSTANT_FALSE]: opConstantFalse,
  [OP_CONSTANT_NULL]: opConstantNull,
  [OP_CONSTANT_TRUE]: opConstantTrue,
  [OP_FUNCTION_DEFINITION]: opFunctionDefinition,
  [OP_IF]: opIf,
  [OP_INTERPOLATED_STRING]: opInterpolatedString,
  [OP_NEGATION]: opNegation,
  [OP_NOT]: opNot,
  [OP_NULL_COALESCENCE]: opNullCoalescence,
  [OP_NUMBER]: opNumber,
  [OP_OBJECT_CONSTRUCTOR]: opObjectConstructor,
  [OP_OR]: opOr,
  [OP_OUTPUT_CONCAT]: opOutputConcat,
  [OP_STRING]: opString,
  [OP_TRY]: opTry,
  [OP_VARIABLE]: opVariable,
  [OP_VARIABLE_DEFINITION]: opVariableDefinition,
  [OP_VOID]: opVoid,
};

export default ops;
