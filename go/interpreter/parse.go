package interpreter

import "github.com/2manyvcos/jpl/go/definition"

type ParserContext struct {
	Interpreter JPLInterpreter
}

// Parse a single program at i.
// Throws an error if src contains additional content.
func parseEntrypoint(src string, i int, c *ParserContext) (n int, result definition.Pipe, err error) {
	n = i

	iResult, opsResult, errResult := parseProgram(src, n, c)
	if errResult != nil {
		return 0, nil, err
	}
	n = iResult

	if _, isEnd := eot(src, n, c); !isEnd {
		return 0, nil, errorUnexpectedToken(src, n, c, errorOptions{Operator: "program", Message: "expected EOT"})
	}

	return iResult, opsResult, errResult
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

		iM, isM, errM := matchWord(src, n, c, matchOptions{Phrase: "|", NotBeforeSet: "="})
		if errM != nil {
			return 0, nil, errM
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

		iM, isM, errM := matchWord(src, n, c, matchOptions{Phrase: "|", NotBeforeSet: "="})
		if errM != nil {
			return 0, nil, errM
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

		iM, isM, errM := matchWord(src, n, c, matchOptions{Phrase: ","})
		if errM != nil {
			return 0, nil, errM
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

	iM, isM, errM := matchWord(src, n, c, matchOptions{Phrase: "try", SpaceAfter: true})
	if errM != nil {
		return 0, nil, errM
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
	iM, isM, errM = matchWord(src, n, c, matchOptions{SpaceBefore: true, Phrase: "catch", SpaceAfter: true})
	if errM != nil {
		return 0, nil, errM
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

		iM, isM, errM := matchWord(src, n, c, matchOptions{SpaceBefore: true, Phrase: "or", SpaceAfter: true})
		if errM != nil {
			return 0, nil, errM
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

		iM, isM, errM := matchWord(src, n, c, matchOptions{SpaceBefore: true, Phrase: "and", SpaceAfter: true})
		if errM != nil {
			return 0, nil, errM
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
	panic("TODO:")
}
