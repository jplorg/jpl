package interpreter

import "github.com/2manyvcos/jpl/go/definition"

type ParserContext struct {
	Interpreter JPLInterpreter
}

// Parse a single program at i.
// Throws an error if src contains additional content.
func parseEntrypoint(src string, i int, c *ParserContext) (n int, result definition.Pipe, err error) {
	n = i

	iResult, opsResult, err := parseProgram(src, n, c)
	if err != nil {
		return 0, nil, err
	}
	n = iResult

	if _, isEnd := eot(src, n, c); !isEnd {
		return 0, nil, errorUnexpectedToken(src, n, c, errorOptions{Operator: "program", Message: "expected EOT"})
	}

	return iResult, opsResult, nil
}

// Parse program at i
func parseProgram(src string, i int, c *ParserContext) (n int, result definition.Pipe, err error) {
	n = i

	if n, _, err = walkWhitespace(src, n, c); err != nil {
		return 0, nil, err
	}

	return opPipe(src, n, c)
}

// Parse pipe at i
func opPipe(src string, i int, c *ParserContext) (n int, result definition.Pipe, err error) {
	n = i

	var pipe definition.Pipe
	for {
		var ops definition.Pipe
		if n, ops, err = opOutputConcat(src, n, c); err != nil {
			return 0, nil, err
		}
		pipe = append(pipe, ops...)

		iM, isM, err := matchWord(src, n, c, matchOptions{Phrase: "|", NotBeforeSet: "="})
		if err != nil {
			return 0, nil, err
		}
		if !isM {
			break
		}
		n = iM
	}

	return n, pipe, nil
}

// Parse subpipe at i
func opSubPipe(src string, i int, c *ParserContext) (n int, result definition.Pipe, err error) {
	n = i

	var pipe definition.Pipe
	for {
		var ops definition.Pipe
		if n, ops, err = opTry(src, n, c); err != nil {
			return 0, nil, err
		}
		pipe = append(pipe, ops...)

		iM, isM, err := matchWord(src, n, c, matchOptions{Phrase: "|", NotBeforeSet: "="})
		if err != nil {
			return 0, nil, err
		}
		if !isM {
			break
		}
		n = iM
	}

	return n, pipe, nil
}

// Parse subroute at i
func opSubRoute(src string, i int, c *ParserContext) (n int, result definition.Pipe, err error) {
	return opTry(src, i, c)
}

// Parse output concat at i
func opOutputConcat(src string, i int, c *ParserContext) (n int, result definition.Pipe, err error) {
	n = i

	var pipes []definition.Pipe
	for {
		var ops definition.Pipe
		if n, ops, err = opTry(src, n, c); err != nil {
			return 0, nil, err
		}
		pipes = append(pipes, ops)

		iM, isM, err := matchWord(src, n, c, matchOptions{Phrase: ","})
		if err != nil {
			return 0, nil, err
		}
		if !isM {
			break
		}
		n = iM
	}

	if len(pipes) == 1 {
		return n, pipes[0], nil
	}

	return n, definition.Pipe{{OP: definition.OP_OUTPUT_CONCAT, Params: definition.JPLInstructionParams{Pipes: pipes}}}, nil
}

// Parse try at i
func opTry(src string, i int, c *ParserContext) (n int, result definition.Pipe, err error) {
	n = i

	iM, isM, err := matchWord(src, n, c, matchOptions{Phrase: "try", SpaceAfter: true})
	if err != nil {
		return 0, nil, err
	}
	if !isM {
		return opOr(src, n, c)
	}
	n = iM

	var opsTry definition.Pipe
	if n, opsTry, err = opOr(src, n, c); err != nil {
		return 0, nil, err
	}

	var opsCatch definition.Pipe
	iM, isM, err = matchWord(src, n, c, matchOptions{SpaceBefore: true, Phrase: "catch", SpaceAfter: true})
	if err != nil {
		return 0, nil, err
	}
	if isM {
		n = iM

		if n, opsCatch, err = opOr(src, n, c); err != nil {
			return 0, nil, err
		}
	} else {
		opsCatch = definition.Pipe{{OP: definition.OP_VOID}}
	}

	return n, definition.Pipe{{OP: definition.OP_TRY, Params: definition.JPLInstructionParams{Try: opsTry, Catch: opsCatch}}}, nil
}

// Parse or at i
func opOr(src string, i int, c *ParserContext) (n int, result definition.Pipe, err error) {
	n = i

	var pipes []definition.Pipe
	for {
		var ops definition.Pipe
		if n, ops, err = opAnd(src, n, c); err != nil {
			return 0, nil, err
		}
		pipes = append(pipes, ops)

		iM, isM, err := matchWord(src, n, c, matchOptions{SpaceBefore: true, Phrase: "or", SpaceAfter: true})
		if err != nil {
			return 0, nil, err
		}
		if !isM {
			break
		}
		n = iM
	}

	if len(pipes) == 1 {
		return n, pipes[0], nil
	}

	return n, definition.Pipe{{OP: definition.OP_OR, Params: definition.JPLInstructionParams{Pipes: pipes}}}, nil
}

// Parse and at i
func opAnd(src string, i int, c *ParserContext) (n int, result definition.Pipe, err error) {
	n = i

	var pipes []definition.Pipe
	for {
		var ops definition.Pipe
		if n, ops, err = opEquality(src, n, c); err != nil {
			return 0, nil, err
		}
		pipes = append(pipes, ops)

		iM, isM, err := matchWord(src, n, c, matchOptions{SpaceBefore: true, Phrase: "and", SpaceAfter: true})
		if err != nil {
			return 0, nil, err
		}
		if !isM {
			break
		}
		n = iM
	}

	if len(pipes) == 1 {
		return n, pipes[0], nil
	}

	return n, definition.Pipe{{OP: definition.OP_AND, Params: definition.JPLInstructionParams{Pipes: pipes}}}, nil
}

// Parse equality at i
func opEquality(src string, i int, c *ParserContext) (n int, result definition.Pipe, err error) {
	n = i

	var ops definition.Pipe
	if n, ops, err = opComparison(src, n, c); err != nil {
		return 0, nil, err
	}

	var comparisons []definition.JPLComparison
	for {
		iM, isM, err := matchWord(src, n, c, matchOptions{Phrase: "=="})
		if err != nil {
			return 0, nil, err
		}
		if isM {
			n = iM

			var opsBy definition.Pipe
			if n, opsBy, err = opComparison(src, n, c); err != nil {
				return 0, nil, err
			}
			comparisons = append(comparisons, definition.JPLComparison{OP: definition.OPC_EQUAL, Params: definition.JPLComparisonParams{By: opsBy}})
			continue
		}

		iM, isM, err = matchWord(src, n, c, matchOptions{Phrase: "!="})
		if err != nil {
			return 0, nil, err
		}
		if isM {
			n = iM

			var opsBy definition.Pipe
			if n, opsBy, err = opComparison(src, n, c); err != nil {
				return 0, nil, err
			}
			comparisons = append(comparisons, definition.JPLComparison{OP: definition.OPC_UNEQUAL, Params: definition.JPLComparisonParams{By: opsBy}})
			continue
		}

		break
	}

	if len(comparisons) == 0 {
		return n, ops, nil
	}

	return n, definition.Pipe{{OP: definition.OP_COMPARISON, Params: definition.JPLInstructionParams{Pipe: ops, Comparisons: comparisons}}}, nil
}

// Parse comparison at i
func opComparison(src string, i int, c *ParserContext) (n int, result definition.Pipe, err error) {
	n = i

	var ops definition.Pipe
	if n, ops, err = opNot(src, n, c); err != nil {
		return 0, nil, err
	}

	var comparisons []definition.JPLComparison
	for {
		iM, isM, err := matchWord(src, n, c, matchOptions{Phrase: "<="})
		if err != nil {
			return 0, nil, err
		}
		if isM {
			n = iM

			var opsBy definition.Pipe
			if n, opsBy, err = opNot(src, n, c); err != nil {
				return 0, nil, err
			}
			comparisons = append(comparisons, definition.JPLComparison{OP: definition.OPC_LESSEQUAL, Params: definition.JPLComparisonParams{By: opsBy}})
			continue
		}

		iM, isM, err = matchWord(src, n, c, matchOptions{Phrase: "<"})
		if err != nil {
			return 0, nil, err
		}
		if isM {
			n = iM

			var opsBy definition.Pipe
			if n, opsBy, err = opNot(src, n, c); err != nil {
				return 0, nil, err
			}
			comparisons = append(comparisons, definition.JPLComparison{OP: definition.OPC_LESS, Params: definition.JPLComparisonParams{By: opsBy}})
			continue
		}

		iM, isM, err = matchWord(src, n, c, matchOptions{Phrase: ">="})
		if err != nil {
			return 0, nil, err
		}
		if isM {
			n = iM

			var opsBy definition.Pipe
			if n, opsBy, err = opNot(src, n, c); err != nil {
				return 0, nil, err
			}
			comparisons = append(comparisons, definition.JPLComparison{OP: definition.OPC_GREATEREQUAL, Params: definition.JPLComparisonParams{By: opsBy}})
			continue
		}

		iM, isM, err = matchWord(src, n, c, matchOptions{Phrase: ">"})
		if err != nil {
			return 0, nil, err
		}
		if isM {
			n = iM

			var opsBy definition.Pipe
			if n, opsBy, err = opNot(src, n, c); err != nil {
				return 0, nil, err
			}
			comparisons = append(comparisons, definition.JPLComparison{OP: definition.OPC_GREATER, Params: definition.JPLComparisonParams{By: opsBy}})
			continue
		}

		break
	}

	if len(comparisons) == 0 {
		return n, ops, nil
	}

	return n, definition.Pipe{{OP: definition.OP_COMPARISON, Params: definition.JPLInstructionParams{Pipe: ops, Comparisons: comparisons}}}, nil
}

// Parse not at i
func opNot(src string, i int, c *ParserContext) (n int, result definition.Pipe, err error) {
	n = i

	iM, isM, err := matchWord(src, n, c, matchOptions{Phrase: "not", SpaceAfter: true})
	if err != nil {
		return 0, nil, err
	}
	if !isM {
		return opErrorSuppression(src, n, c)
	}
	n = iM

	var ops definition.Pipe
	if n, ops, err = opErrorSuppression(src, n, c); err != nil {
		return 0, nil, err
	}

	result = append(result, ops...)
	result = append(result, definition.JPLInstruction{OP: definition.OP_NOT})
	return n, result, nil
}

// Parse error suppression at i
func opErrorSuppression(src string, i int, c *ParserContext) (n int, result definition.Pipe, err error) {
	n = i

	iResult, opsResult, err := opDifference(src, n, c)
	if err != nil {
		return 0, nil, err
	}
	n = iResult

	iM, isM, err := matchWord(src, n, c, matchOptions{Phrase: "?", NotBeforeSet: "?="})
	if err != nil {
		return 0, nil, err
	}
	if !isM {
		return iResult, opsResult, nil
	}
	n = iM

	return n, definition.Pipe{{OP: definition.OP_TRY, Params: definition.JPLInstructionParams{Try: opsResult, Catch: definition.Pipe{{OP: definition.OP_VOID}}}}}, nil
}

// Parse difference at i
func opDifference(src string, i int, c *ParserContext) (n int, result definition.Pipe, err error) {
	n = i

	var ops definition.Pipe
	if n, ops, err = opMultiplication(src, n, c); err != nil {
		return 0, nil, err
	}

	var operations []definition.JPLOperation
	for {
		iM, isM, err := matchWord(src, n, c, matchOptions{Phrase: "+", NotBeforeSet: "="})
		if err != nil {
			return 0, nil, err
		}
		if isM {
			n = iM

			var opsBy definition.Pipe
			if n, opsBy, err = opMultiplication(src, n, c); err != nil {
				return 0, nil, err
			}
			operations = append(operations, definition.JPLOperation{OP: definition.OPM_ADDITION, Params: definition.JPLOperationParams{By: opsBy}})
			continue
		}

		iM, isM, err = matchWord(src, n, c, matchOptions{Phrase: "-", NotBeforeSet: "=>"})
		if err != nil {
			return 0, nil, err
		}
		if isM {
			n = iM

			var opsBy definition.Pipe
			if n, opsBy, err = opMultiplication(src, n, c); err != nil {
				return 0, nil, err
			}
			operations = append(operations, definition.JPLOperation{OP: definition.OPM_SUBTRACTION, Params: definition.JPLOperationParams{By: opsBy}})
			continue
		}

		break
	}

	if len(operations) == 0 {
		return n, ops, nil
	}

	return n, definition.Pipe{{OP: definition.OP_CALCULATION, Params: definition.JPLInstructionParams{Pipe: ops, Operations: operations}}}, nil
}

// Parse multiplication at i
func opMultiplication(src string, i int, c *ParserContext) (n int, result definition.Pipe, err error) {
	n = i

	var ops definition.Pipe
	if n, ops, err = opNullCoalescence(src, n, c); err != nil {
		return 0, nil, err
	}

	var operations []definition.JPLOperation
	for {
		iM, isM, err := matchWord(src, n, c, matchOptions{Phrase: "*", NotBeforeSet: "="})
		if err != nil {
			return 0, nil, err
		}
		if isM {
			n = iM

			var opsBy definition.Pipe
			if n, opsBy, err = opNullCoalescence(src, n, c); err != nil {
				return 0, nil, err
			}
			operations = append(operations, definition.JPLOperation{OP: definition.OPM_MULTIPLICATION, Params: definition.JPLOperationParams{By: opsBy}})
			continue
		}

		iM, isM, err = matchWord(src, n, c, matchOptions{Phrase: "/", NotBeforeSet: "="})
		if err != nil {
			return 0, nil, err
		}
		if isM {
			n = iM

			var opsBy definition.Pipe
			if n, opsBy, err = opNullCoalescence(src, n, c); err != nil {
				return 0, nil, err
			}
			operations = append(operations, definition.JPLOperation{OP: definition.OPM_DIVISION, Params: definition.JPLOperationParams{By: opsBy}})
			continue
		}

		iM, isM, err = matchWord(src, n, c, matchOptions{Phrase: "%", NotBeforeSet: "="})
		if err != nil {
			return 0, nil, err
		}
		if isM {
			n = iM

			var opsBy definition.Pipe
			if n, opsBy, err = opNullCoalescence(src, n, c); err != nil {
				return 0, nil, err
			}
			operations = append(operations, definition.JPLOperation{OP: definition.OPM_REMAINDER, Params: definition.JPLOperationParams{By: opsBy}})
			continue
		}

		break
	}

	if len(operations) == 0 {
		return n, ops, nil
	}

	return n, definition.Pipe{{OP: definition.OP_CALCULATION, Params: definition.JPLInstructionParams{Pipe: ops, Operations: operations}}}, nil
}

// Parse null coalescence at i
func opNullCoalescence(src string, i int, c *ParserContext) (n int, result definition.Pipe, err error) {
	n = i

	var pipes []definition.Pipe
	for {
		var ops definition.Pipe
		if n, ops, err = opNegation(src, n, c); err != nil {
			return 0, nil, err
		}
		pipes = append(pipes, ops)

		iM, isM, err := matchWord(src, n, c, matchOptions{Phrase: "??"})
		if err != nil {
			return 0, nil, err
		}
		if !isM {
			break
		}
		n = iM
	}

	if len(pipes) == 1 {
		return n, pipes[0], nil
	}

	return n, definition.Pipe{{OP: definition.OP_NULL_COALESCENCE, Params: definition.JPLInstructionParams{Pipes: pipes}}}, nil
}

// Parse negation at i
func opNegation(src string, i int, c *ParserContext) (n int, result definition.Pipe, err error) {
	n = i

	iM, isM, err := matchWord(src, n, c, matchOptions{Phrase: "-", NotBeforeSet: "=>"})
	if err != nil {
		return 0, nil, err
	}
	if !isM {
		return opIf(src, n, c)
	}
	n = iM

	var ops definition.Pipe
	if n, ops, err = opIf(src, n, c); err != nil {
		return 0, nil, err
	}

	result = append(result, ops...)
	result = append(result, definition.JPLInstruction{OP: definition.OP_NEGATION})
	return n, result, nil
}

// Parse if at i
func opIf(src string, i int, c *ParserContext) (n int, result definition.Pipe, err error) {
	n = i

	iM, isM, err := matchWord(src, n, c, matchOptions{Phrase: "if", SpaceAfter: true})
	if err != nil {
		return 0, nil, err
	}
	if !isM {
		return opConstant(src, n, c)
	}
	n = iM

	var ifs []definition.JPLIfThen
	for {
		var opsIf definition.Pipe
		if n, opsIf, err = opPipe(src, n, c); err != nil {
			return 0, nil, err
		}

		iM, isM, err = matchWord(src, n, c, matchOptions{SpaceBefore: true, Phrase: "then", SpaceAfter: true})
		if err != nil {
			return 0, nil, err
		}
		n = iM
		if !isM {
			return 0, nil, errorUnexpectedToken(src, n, c, errorOptions{Operator: "if statement", Message: "expected 'then'"})
		}

		var opsThen definition.Pipe
		if n, opsThen, err = opPipe(src, n, c); err != nil {
			return 0, nil, err
		}
		ifs = append(ifs, definition.JPLIfThen{If: opsIf, Then: opsThen})

		iM, isM, err = matchWord(src, n, c, matchOptions{SpaceBefore: true, Phrase: "elif", SpaceAfter: true})
		if err != nil {
			return 0, nil, err
		}
		if !isM {
			break
		}
		n = iM
	}

	var opsElse definition.Pipe
	iM, isM, err = matchWord(src, n, c, matchOptions{SpaceBefore: true, Phrase: "else", SpaceAfter: true})
	if err != nil {
		return 0, nil, err
	}
	if isM {
		n = iM

		if n, opsElse, err = opPipe(src, n, c); err != nil {
			return 0, nil, err
		}
	} else {
		opsElse = nil
	}

	iM, isM, err = matchWord(src, n, c, matchOptions{SpaceBefore: true, Phrase: "end"})
	if err != nil {
		return 0, nil, err
	}
	n = iM
	if !isM {
		return 0, nil, errorUnexpectedToken(src, n, c, errorOptions{Operator: "if statement", Message: "expected 'end'"})
	}

	return n, definition.Pipe{{OP: definition.OP_IF, Params: definition.JPLInstructionParams{Ifs: ifs, Else: opsElse}}}, nil
}

// Parse constant at i
func opConstant(src string, i int, c *ParserContext) (n int, result definition.Pipe, err error) {
	n = i

	iM, isM, err := matchWord(src, n, c, matchOptions{Phrase: "true", SpaceAfter: true})
	if err != nil {
		return 0, nil, err
	}
	if isM {
		n = iM
		return n, definition.Pipe{{OP: definition.OP_CONSTANT_TRUE}}, nil
	}

	iM, isM, err = matchWord(src, n, c, matchOptions{Phrase: "false", SpaceAfter: true})
	if err != nil {
		return 0, nil, err
	}
	if isM {
		n = iM
		return n, definition.Pipe{{OP: definition.OP_CONSTANT_FALSE}}, nil
	}

	iM, isM, err = matchWord(src, n, c, matchOptions{Phrase: "null", SpaceAfter: true})
	if err != nil {
		return 0, nil, err
	}
	if isM {
		n = iM
		return n, definition.Pipe{{OP: definition.OP_CONSTANT_NULL}}, nil
	}

	return opNumber(src, n, c)
}

// Parse number at i
func opNumber(src string, i int, c *ParserContext) (n int, result definition.Pipe, err error) {
	n = i

	iResult, isResult, opsResult, err := number(src, n, c)
	if err != nil {
		return 0, nil, err
	}
	if !isResult {
		return opNamedFunctionDefinition(src, n, c)
	}
	n = iResult

	return n, opsResult, nil
}

// Parse named function definition at i
func opNamedFunctionDefinition(src string, i int, c *ParserContext) (n int, result definition.Pipe, err error) {
	n = i

	iM, isM, err := matchWord(src, n, c, matchOptions{Phrase: "func", SpaceAfter: true})
	if err != nil {
		return 0, nil, err
	}
	if !isM {
		return opFunctionDefinition(src, n, c)
	}
	n = iM

	iV, isV, name, _, err := safeVariable(src, n, c)
	if err != nil {
		return 0, nil, err
	}
	if !isV {
		return opFunctionDefinition(src, i, c)
	}
	n = iV

	var argNames []string
	if n, argNames, err = functionHeader(src, n, c); err != nil {
		return 0, nil, err
	}

	var ops definition.Pipe
	if n, ops, err = opSubRoute(src, n, c); err != nil {
		return 0, nil, err
	}

	return n, definition.Pipe{{OP: definition.OP_VARIABLE_DEFINITION, Params: definition.JPLInstructionParams{Name: name, Pipe: definition.Pipe{{OP: definition.OP_FUNCTION_DEFINITION, Params: definition.JPLInstructionParams{ArgNames: argNames, Pipe: ops}}}}}}, nil
}

// Parse function definition at i
func opFunctionDefinition(src string, i int, c *ParserContext) (n int, result definition.Pipe, err error) {
	n = i

	iM, isM, err := matchWord(src, n, c, matchOptions{Phrase: "func", SpaceAfter: true})
	if err != nil {
		return 0, nil, err
	}
	if !isM {
		return opVariableAccess(src, n, c)
	}
	n = iM

	var argNames []string
	if n, argNames, err = functionHeader(src, n, c); err != nil {
		return 0, nil, err
	}

	var ops definition.Pipe
	if n, ops, err = opSubRoute(src, n, c); err != nil {
		return 0, nil, err
	}

	return n, definition.Pipe{{OP: definition.OP_FUNCTION_DEFINITION, Params: definition.JPLInstructionParams{ArgNames: argNames, Pipe: ops}}}, nil
}

// Parse variable definition at i
func opVariableAccess(src string, i int, c *ParserContext) (n int, result definition.Pipe, err error) {
	panic("TODO:")
}

// Parse function header at i
func functionHeader(src string, i int, c *ParserContext) (n int, argNames []string, err error) {
	n = i

	iM, isM, err := matchWord(src, n, c, matchOptions{Phrase: "("})
	if err != nil {
		return 0, nil, err
	}
	n = iM
	if !isM {
		return 0, nil, errorUnexpectedToken(src, n, c, errorOptions{Operator: "function definition", Message: "expected '('"})
	}

	iM, isM, err = matchWord(src, n, c, matchOptions{Phrase: ")"})
	if err != nil {
		return 0, nil, err
	}
	if isM {
		n = iM
	} else {
		for {
			iV, isV, name, _, err := safeVariable(src, n, c)
			if err != nil {
				return 0, nil, err
			}
			if !isV {
				return 0, nil, errorUnexpectedToken(src, n, c, errorOptions{Operator: "function definition", Message: "expected argument name"})
			}
			n = iV
			argNames = append(argNames, name)

			iM, isM, err = matchWord(src, n, c, matchOptions{Phrase: ")"})
			if err != nil {
				return 0, nil, err
			}
			if isM {
				n = iM
				break
			}

			iM, isM, err = matchWord(src, n, c, matchOptions{Phrase: ","})
			if err != nil {
				return 0, nil, err
			}
			n = iM
			if !isM {
				return 0, nil, errorUnexpectedToken(src, n, c, errorOptions{Operator: "function definition", Message: "expected ',' or ')'"})
			}
		}
	}

	iM, isM, err = matchWord(src, n, c, matchOptions{Phrase: ":"})
	if err != nil {
		return 0, nil, err
	}
	n = iM
	if !isM {
		return 0, nil, errorUnexpectedToken(src, n, c, errorOptions{Operator: "function definition", Message: "expected ':'"})
	}

	return n, argNames, nil
}

// Parse number at i
func number(src string, i int, c *ParserContext) (n int, is bool, result definition.Pipe, err error) {
	n = i
	var value string

	for {
		iSet, isSet, valueSet := matchSet(src, n, c, matchSetOptions{Set: "0123456789"})
		if !isSet {
			break
		}
		n = iSet
		is = true
		value += valueSet
	}

	if !is {
		return n, false, nil, nil
	}

	iM, isM := match(src, n, c, matchOptions{Phrase: "."})
	if isM {
		n = iM
		value += "."

		for {
			iSet, isSet, valueSet := matchSet(src, n, c, matchSetOptions{Set: "0123456789"})
			if !isSet {
				break
			}
			n = iSet
			value += valueSet
		}
	}

	iSet, isSet, valueSet := matchSet(src, n, c, matchSetOptions{Set: "eE"})
	if isSet {
		n = iSet
		value += valueSet

		iSet, isSet, valueSet = matchSet(src, n, c, matchSetOptions{Set: "+-"})
		if isSet {
			n = iSet
			value += valueSet
		}

		var isE bool
		for {
			iSet, isSet, valueSet = matchSet(src, n, c, matchSetOptions{Set: "0123456789"})
			if !isSet {
				break
			}
			n = iSet
			isE = true
			value += valueSet
		}
		if !isE {
			return 0, false, nil, errorUnexpectedToken(src, n, c, errorOptions{Operator: "number", Message: "expected digit"})
		}
	}

	if n, _, err = walkWhitespace(src, n, c); err != nil {
		return 0, false, nil, err
	}

	return n, true, definition.Pipe{{OP: definition.OP_NUMBER, Params: definition.JPLInstructionParams{Value: value}}}, nil
}
