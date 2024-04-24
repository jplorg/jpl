package interpreter

import (
	"strconv"
	"unicode/utf16"
	"unicode/utf8"

	"github.com/2manyvcos/jpl/go/definition"
)

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

// Parse function header at i
func parseFunctionHeader(src string, i int, c *ParserContext) (n int, argNames []string, err error) {
	n = i

	iM, isM, err := matchWord(src, n, c, matchOptions{Phrase: "("})
	if err != nil {
		return 0, nil, err
	}
	n = iM
	if !isM {
		return 0, nil, errorUnexpectedToken(src, n, c, errorOptions{
			Operator: "function definition",
			Message:  "expected '('",
		})
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
				return 0, nil, errorUnexpectedToken(src, n, c, errorOptions{
					Operator: "function definition",
					Message:  "expected argument name",
				})
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
				return 0, nil, errorUnexpectedToken(src, n, c, errorOptions{
					Operator: "function definition",
					Message:  "expected ',' or ')'",
				})
			}
		}
	}

	iM, isM, err = matchWord(src, n, c, matchOptions{Phrase: ":"})
	if err != nil {
		return 0, nil, err
	}
	n = iM
	if !isM {
		return 0, nil, errorUnexpectedToken(src, n, c, errorOptions{
			Operator: "function definition",
			Message:  "expected ':'",
		})
	}

	return n, argNames, nil
}

type accessOptions struct {
	Identity bool
}

// Parse access at i
func parseAccess(src string, i int, c *ParserContext, options accessOptions) (n int, is bool, operations []definition.JPLOperation, canAssign bool, err error) {
	n = i

	canAssign = true
	for {
		iM, isM, err := matchWord(src, n, c, matchOptions{Phrase: "."})
		if err != nil {
			return 0, false, nil, false, err
		}
		isIdentity := options.Identity && len(operations) == 0
		if !isIdentity && isM {
			n = iM

			iV, isV, name, _, err := safeVariable(src, n, c)
			if err != nil {
				return 0, false, nil, false, err
			}
			if !isV {
				return 0, false, nil, false, errorUnexpectedToken(src, n, c, errorOptions{
					Operator: "field access operator",
					Message:  "expected field name",
				})
			}
			n = iV

			var optional bool
			iM, isM, err = matchWord(src, n, c, matchOptions{Phrase: "?", NotBeforeSet: "?="})
			if err != nil {
				return 0, false, nil, false, err
			}
			if isM {
				n = iM
				optional = true
			}
			operations = append(operations, definition.JPLOperation{OP: definition.OPA_FIELD, Params: definition.JPLOperationParams{Pipe: definition.Pipe{{OP: definition.OP_STRING, Params: definition.JPLInstructionParams{Value: name}}}, Optional: optional}})
			continue
		}

		iM, isM, err = matchWord(src, n, c, matchOptions{Phrase: "["})
		if err != nil {
			return 0, false, nil, false, err
		}
		if isM {
			n = iM

			iM, isM, err = matchWord(src, n, c, matchOptions{Phrase: "]"})
			if err != nil {
				return 0, false, nil, false, err
			}
			if isM {
				n = iM

				var optional bool
				iM, isM, err = matchWord(src, n, c, matchOptions{Phrase: "?", NotBeforeSet: "?="})
				if err != nil {
					return 0, false, nil, false, err
				}
				if isM {
					n = iM
					optional = true
				}
				operations = append(operations, definition.JPLOperation{OP: definition.OPA_ITER, Params: definition.JPLOperationParams{Optional: optional}})
				continue
			}

			iM, isM, err = matchWord(src, n, c, matchOptions{Phrase: ":"})
			if err != nil {
				return 0, false, nil, false, err
			}
			if isM {
				n = iM

				iM, isM, err = matchWord(src, n, c, matchOptions{Phrase: "]"})
				if err != nil {
					return 0, false, nil, false, err
				}
				if isM {
					n = iM

					var optional bool
					iM, isM, err = matchWord(src, n, c, matchOptions{Phrase: "?", NotBeforeSet: "?="})
					if err != nil {
						return 0, false, nil, false, err
					}
					if isM {
						n = iM
						optional = true
					}
					operations = append(operations, definition.JPLOperation{OP: definition.OPA_SLICE, Params: definition.JPLOperationParams{From: definition.Pipe{{OP: definition.OP_CONSTANT_NULL}}, To: definition.Pipe{{OP: definition.OP_CONSTANT_NULL}}, Optional: optional}})
					continue
				}

				var opsRight definition.Pipe
				if n, opsRight, err = opPipe(src, n, c); err != nil {
					return 0, false, nil, false, err
				}

				iM, isM, err = matchWord(src, n, c, matchOptions{Phrase: "]"})
				if err != nil {
					return 0, false, nil, false, err
				}
				n = iM
				if !isM {
					return 0, false, nil, false, errorUnexpectedToken(src, n, c, errorOptions{
						Operator: "array slice operator",
						Message:  "expected ']'",
					})
				}

				var optional bool
				iM, isM, err = matchWord(src, n, c, matchOptions{Phrase: "?", NotBeforeSet: "?="})
				if err != nil {
					return 0, false, nil, false, err
				}
				if isM {
					n = iM
					optional = true
				}
				operations = append(operations, definition.JPLOperation{OP: definition.OPA_SLICE, Params: definition.JPLOperationParams{From: definition.Pipe{{OP: definition.OP_CONSTANT_NULL}}, To: opsRight, Optional: optional}})
				continue
			}

			var opsLeft definition.Pipe
			if n, opsLeft, err = opPipe(src, n, c); err != nil {
				return 0, false, nil, false, err
			}

			iM, isM, err = matchWord(src, n, c, matchOptions{Phrase: "]"})
			if err != nil {
				return 0, false, nil, false, err
			}
			if isM {
				n = iM

				var optional bool
				iM, isM, err = matchWord(src, n, c, matchOptions{Phrase: "?", NotBeforeSet: "?="})
				if err != nil {
					return 0, false, nil, false, err
				}
				if isM {
					n = iM
					optional = true
				}
				operations = append(operations, definition.JPLOperation{OP: definition.OPA_FIELD, Params: definition.JPLOperationParams{Pipe: opsLeft, Optional: optional}})
				continue
			}

			iM, isM, err = matchWord(src, n, c, matchOptions{Phrase: ":"})
			if err != nil {
				return 0, false, nil, false, err
			}
			if !isM {
				return 0, false, nil, false, errorUnexpectedToken(src, n, c, errorOptions{
					Operator: "variable access operator",
					Message:  "expected ':' or ']'",
				})
			}
			n = iM

			iM, isM, err = matchWord(src, n, c, matchOptions{Phrase: "]"})
			if err != nil {
				return 0, false, nil, false, err
			}
			if isM {
				n = iM

				var optional bool
				iM, isM, err = matchWord(src, n, c, matchOptions{Phrase: "?", NotBeforeSet: "?="})
				if err != nil {
					return 0, false, nil, false, err
				}
				if isM {
					n = iM
					optional = true
				}
				operations = append(operations, definition.JPLOperation{OP: definition.OPA_SLICE, Params: definition.JPLOperationParams{From: opsLeft, To: definition.Pipe{{OP: definition.OP_CONSTANT_NULL}}, Optional: optional}})
				continue
			}

			var opsRight definition.Pipe
			if n, opsRight, err = opPipe(src, n, c); err != nil {
				return 0, false, nil, false, err
			}

			iM, isM, err = matchWord(src, n, c, matchOptions{Phrase: "]"})
			if err != nil {
				return 0, false, nil, false, err
			}
			n = iM
			if !isM {
				return 0, false, nil, false, errorUnexpectedToken(src, n, c, errorOptions{
					Operator: "array slice operator",
					Message:  "expected ']'",
				})
			}

			var optional bool
			iM, isM, err = matchWord(src, n, c, matchOptions{Phrase: "?", NotBeforeSet: "?="})
			if err != nil {
				return 0, false, nil, false, err
			}
			if isM {
				n = iM
				optional = true
			}
			operations = append(operations, definition.JPLOperation{OP: definition.OPA_SLICE, Params: definition.JPLOperationParams{From: opsLeft, To: opsRight, Optional: optional}})
			continue
		}

		var bound bool
		iM, isM, err = matchWord(src, n, c, matchOptions{Phrase: "->"})
		if err != nil {
			return 0, false, nil, false, err
		}
		if isM {
			n = iM
			bound = true
		}

		iM, isM, err = matchWord(src, n, c, matchOptions{Phrase: "("})
		if err != nil {
			return 0, false, nil, false, err
		}
		if !isM && bound {
			return 0, false, nil, false, errorUnexpectedToken(src, n, c, errorOptions{
				Operator: "bound function call",
				Message:  "expected '('",
			})
		}
		if isM {
			n = iM

			var args []definition.Pipe
			iM, isM, err = matchWord(src, n, c, matchOptions{Phrase: ")"})
			if err != nil {
				return 0, false, nil, false, err
			}
			if isM {
				n = iM
			} else {
				for {
					var opsArg definition.Pipe
					if n, opsArg, err = opSubPipe(src, n, c); err != nil {
						return 0, false, nil, false, err
					}
					args = append(args, opsArg)

					iM, isM, err = matchWord(src, n, c, matchOptions{Phrase: ")"})
					if err != nil {
						return 0, false, nil, false, err
					}
					if isM {
						n = iM
						break
					}

					iM, isM, err = matchWord(src, n, c, matchOptions{Phrase: ","})
					if err != nil {
						return 0, false, nil, false, err
					}
					n = iM
					if !isM {
						return 0, false, nil, false, errorUnexpectedToken(src, n, c, errorOptions{
							Operator: "function call",
							Message:  "expected ',' or ')'",
						})
					}
				}
			}

			var optional bool
			iM, isM, err = matchWord(src, n, c, matchOptions{Phrase: "?", NotBeforeSet: "?="})
			if err != nil {
				return 0, false, nil, false, err
			}
			if isM {
				n = iM
				optional = true
			}
			operations = append(operations, definition.JPLOperation{OP: definition.OPA_FUNCTION, Params: definition.JPLOperationParams{Args: args, Bound: bound, Optional: optional}})
			canAssign = false
			continue
		}

		break
	}

	return n, len(operations) > 0, operations, canAssign, nil
}

// Parse assignment at i
func parseAssignment(src string, i int, c *ParserContext) (n int, is bool, assignment *definition.JPLAssignment, err error) {
	n = i

	iM, isM, err := matchWord(src, n, c, matchOptions{Phrase: "=", NotBeforeSet: "="})
	if err != nil {
		return 0, false, nil, err
	}
	if isM {
		n = iM

		var ops definition.Pipe
		if n, ops, err = opSubRoute(src, n, c); err != nil {
			return 0, false, nil, err
		}
		return n, true, &definition.JPLAssignment{OP: definition.OPU_SET, Params: definition.JPLAssignmentParams{Pipe: ops}}, nil
	}

	iM, isM, err = matchWord(src, n, c, matchOptions{Phrase: "|="})
	if err != nil {
		return 0, false, nil, err
	}
	if isM {
		n = iM

		var ops definition.Pipe
		if n, ops, err = opSubRoute(src, n, c); err != nil {
			return 0, false, nil, err
		}
		return n, true, &definition.JPLAssignment{OP: definition.OPU_UPDATE, Params: definition.JPLAssignmentParams{Pipe: ops}}, nil
	}

	iM, isM, err = matchWord(src, n, c, matchOptions{Phrase: "+="})
	if err != nil {
		return 0, false, nil, err
	}
	if isM {
		n = iM

		var ops definition.Pipe
		if n, ops, err = opSubRoute(src, n, c); err != nil {
			return 0, false, nil, err
		}
		return n, true, &definition.JPLAssignment{OP: definition.OPU_ADDITION, Params: definition.JPLAssignmentParams{Pipe: ops}}, nil
	}

	iM, isM, err = matchWord(src, n, c, matchOptions{Phrase: "-="})
	if err != nil {
		return 0, false, nil, err
	}
	if isM {
		n = iM

		var ops definition.Pipe
		if n, ops, err = opSubRoute(src, n, c); err != nil {
			return 0, false, nil, err
		}
		return n, true, &definition.JPLAssignment{OP: definition.OPU_SUBTRACTION, Params: definition.JPLAssignmentParams{Pipe: ops}}, nil
	}

	iM, isM, err = matchWord(src, n, c, matchOptions{Phrase: "*="})
	if err != nil {
		return 0, false, nil, err
	}
	if isM {
		n = iM

		var ops definition.Pipe
		if n, ops, err = opSubRoute(src, n, c); err != nil {
			return 0, false, nil, err
		}
		return n, true, &definition.JPLAssignment{OP: definition.OPU_MULTIPLICATION, Params: definition.JPLAssignmentParams{Pipe: ops}}, nil
	}

	iM, isM, err = matchWord(src, n, c, matchOptions{Phrase: "/="})
	if err != nil {
		return 0, false, nil, err
	}
	if isM {
		n = iM

		var ops definition.Pipe
		if n, ops, err = opSubRoute(src, n, c); err != nil {
			return 0, false, nil, err
		}
		return n, true, &definition.JPLAssignment{OP: definition.OPU_DIVISION, Params: definition.JPLAssignmentParams{Pipe: ops}}, nil
	}

	iM, isM, err = matchWord(src, n, c, matchOptions{Phrase: "%="})
	if err != nil {
		return 0, false, nil, err
	}
	if isM {
		n = iM

		var ops definition.Pipe
		if n, ops, err = opSubRoute(src, n, c); err != nil {
			return 0, false, nil, err
		}
		return n, true, &definition.JPLAssignment{OP: definition.OPU_REMAINDER, Params: definition.JPLAssignmentParams{Pipe: ops}}, nil
	}

	iM, isM, err = matchWord(src, n, c, matchOptions{Phrase: "?="})
	if err != nil {
		return 0, false, nil, err
	}
	if isM {
		n = iM

		var ops definition.Pipe
		if n, ops, err = opSubRoute(src, n, c); err != nil {
			return 0, false, nil, err
		}
		return n, true, &definition.JPLAssignment{OP: definition.OPU_NULL_COALESCENCE, Params: definition.JPLAssignmentParams{Pipe: ops}}, nil
	}

	return n, false, nil, nil
}

// Parse number at i
func parseNumber(src string, i int, c *ParserContext) (n int, is bool, result definition.Pipe, err error) {
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

// Parse string at i
func parseString(src string, i int, c *ParserContext) (n int, is bool, result definition.Pipe, err error) {
	n = i
	var value []byte

	iSet, isSet, boundary := matchSet(src, n, c, matchSetOptions{Set: "\"'`"})
	if !isSet {
		return n, false, nil, nil
	}
	n = iSet

	multilineString := boundary == "`"

	var interpolations []definition.JPLInterpolation
	for {
		if iM, isM := match(src, n, c, matchOptions{Phrase: boundary}); isM {
			n = iM
			break
		}

		if iEnd, isEnd := eot(src, n, c); isEnd {
			n = iEnd
			return 0, false, nil, errorUnexpectedToken(src, n, c, errorOptions{
				Operator: "string",
				Message:  "incomplete string literal",
			})
		}

		if iM, isM := match(src, n, c, matchOptions{Phrase: "\\"}); isM {
			n = iM

			iM, isM, err := matchWord(src, n, c, matchOptions{Phrase: "("})
			if err != nil {
				return 0, false, nil, err
			}
			if isM {
				n = iM

				var ops definition.Pipe
				if n, ops, err = opPipe(src, n, c); err != nil {
					return 0, false, nil, err
				}

				iM, isM = match(src, n, c, matchOptions{Phrase: ")"})
				n = iM
				if !isM {
					return 0, false, nil, errorUnexpectedToken(src, n, c, errorOptions{
						Operator: "string interpolation",
						Message:  "expected ')'",
					})
				}
				interpolations = append(interpolations, definition.JPLInterpolation{Before: string(value), Pipe: ops})
				value = nil
				continue
			}

			if multilineString {
				switch src[n] {
				case '\n':
					n += 1
					continue

				case '\r':
					n += 1
					if _, isEnd := eot(src, n, c); !isEnd && src[n] == '\n' {
						n += 1
					}
					continue

				default:
				}
			}

			switch src[n] {
			case '"':
				value = append(value, '"')
				n += 1
				continue

			case '\'':
				value = append(value, '\'')
				n += 1
				continue

			case '`':
				value = append(value, '`')
				n += 1
				continue

			case '\\':
				value = append(value, '\\')
				n += 1
				continue

			case '/':
				value = append(value, '/')
				n += 1
				continue

			case 'b':
				value = append(value, '\b')
				n += 1
				continue

			case 'f':
				value = append(value, '\f')
				n += 1
				continue

			case 'n':
				value = append(value, '\n')
				n += 1
				continue

			case 'r':
				value = append(value, '\r')
				n += 1
				continue

			case 't':
				value = append(value, '\t')
				n += 1
				continue

			case 'u':
				n += 1
				var hexVal string
				for j := 0; j < 4; j += 1 {
					iM, isM, valueM := hex(src, n, c)
					n = iM
					if !isM {
						return 0, false, nil, errorUnexpectedToken(src, n, c, errorOptions{
							Operator: "string",
							Message:  "incomplete unicode escape sequence: expected hex digit",
						})
					}
					hexVal += valueM
				}
				parsed, err := strconv.ParseUint(hexVal, 16, 16)
				if err != nil {
					return 0, false, nil, err
				}
				r := rune(parsed)
				// parse utf16
				if utf16.IsSurrogate(r) {
					iM, isM = match(src, n, c, matchOptions{Phrase: "\\u"})
					if isM {
						n = iM

						hexVal = ""
						for j := 0; j < 4; j += 1 {
							iM, isM, valueM := hex(src, n, c)
							n = iM
							if !isM {
								return 0, false, nil, errorUnexpectedToken(src, n, c, errorOptions{
									Operator: "string",
									Message:  "incomplete unicode escape sequence: expected hex digit",
								})
							}
							hexVal += valueM
						}
						parsed, err = strconv.ParseUint(hexVal, 16, 16)
						if err != nil {
							return 0, false, nil, err
						}
						r2 := rune(parsed)
						// if the second rune is no utf16 surrogate, we append each rune individually to the result value
						if !utf16.IsSurrogate(r2) {
							value = utf8.AppendRune(value, r)
							value = utf8.AppendRune(value, r2)
							continue
						}
						// if we have a valid utf16 surrogate pair, we parse it to utf8 and append it to the result value
						r = utf16.DecodeRune(r, r2)
					}
				}
				// append the resulting rune
				value = utf8.AppendRune(value, r)
				continue

			default:
				return 0, false, nil, errorUnexpectedToken(src, n, c, errorOptions{
					Operator: "string",
					Message:  "invalid escape sequence",
				})
			}
		}

		if multilineString {
			switch src[n] {
			case '\n', '\r', '\t':
				value = append(value, src[n])
				n += 1
				continue

			default:
			}
		}

		if src[n] < 0x20 {
			return 0, false, nil, errorUnexpectedToken(src, n, c, errorOptions{Operator: "string"})
		}

		value = append(value, src[n])
		n += 1
	}

	if n, _, err = walkWhitespace(src, n, c); err != nil {
		return 0, false, nil, err
	}

	if len(interpolations) == 0 {
		return n, true, definition.Pipe{{OP: definition.OP_STRING, Params: definition.JPLInstructionParams{Value: string(value)}}}, nil
	}

	return n, true, definition.Pipe{{OP: definition.OP_INTERPOLATED_STRING, Params: definition.JPLInstructionParams{Interpolations: interpolations, After: string(value)}}}, nil
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
			return 0, nil, errorUnexpectedToken(src, n, c, errorOptions{
				Operator: "if statement",
				Message:  "expected 'then'",
			})
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
		return 0, nil, errorUnexpectedToken(src, n, c, errorOptions{
			Operator: "if statement",
			Message:  "expected 'end'",
		})
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

	iResult, isResult, opsResult, err := parseNumber(src, n, c)
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
	if n, argNames, err = parseFunctionHeader(src, n, c); err != nil {
		return 0, nil, err
	}

	var ops definition.Pipe
	if n, ops, err = opSubRoute(src, n, c); err != nil {
		return 0, nil, err
	}

	return n, definition.Pipe{{
		OP:     definition.OP_VARIABLE_DEFINITION,
		Params: definition.JPLInstructionParams{Name: name, Pipe: definition.Pipe{{OP: definition.OP_FUNCTION_DEFINITION, Params: definition.JPLInstructionParams{ArgNames: argNames, Pipe: ops}}}},
	}}, nil
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
	if n, argNames, err = parseFunctionHeader(src, n, c); err != nil {
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
	n = i

	iV, isV, name, _, err := safeVariable(src, n, c)
	if err != nil {
		return 0, nil, err
	}
	if !isV {
		return opValueAccess(src, n, c)
	}
	n = iV

	iAc, isAc, operations, canAssign, err := parseAccess(src, n, c, accessOptions{})
	if err != nil {
		return 0, nil, err
	}
	if isAc {
		n = iAc
	} else {
		operations = nil
		canAssign = true
	}

	if !canAssign {
		ops := definition.Pipe{{OP: definition.OP_VARIABLE, Params: definition.JPLInstructionParams{Name: name}}}

		if len(operations) == 0 {
			return n, ops, nil
		}
		return n, definition.Pipe{{OP: definition.OP_ACCESS, Params: definition.JPLInstructionParams{Pipe: ops, Operations: operations}}}, nil
	}

	var opAssignment *definition.JPLAssignment
	iAs, isAs, opAssignment, err := parseAssignment(src, n, c)
	if err != nil {
		return 0, nil, err
	}
	if !isAs {
		ops := definition.Pipe{{OP: definition.OP_VARIABLE, Params: definition.JPLInstructionParams{Name: name}}}

		if len(operations) == 0 {
			return n, ops, nil
		}
		return n, definition.Pipe{{OP: definition.OP_ACCESS, Params: definition.JPLInstructionParams{Pipe: ops, Operations: operations}}}, nil
	}
	n = iAs

	if len(operations) == 0 && opAssignment.OP == definition.OPU_SET {
		return n, definition.Pipe{{OP: definition.OP_VARIABLE_DEFINITION, Params: definition.JPLInstructionParams{Name: name, Pipe: opAssignment.Params.Pipe}}}, nil
	}

	return n, definition.Pipe{{
		OP: definition.OP_VARIABLE_DEFINITION,
		Params: definition.JPLInstructionParams{
			Name: name,
			Pipe: definition.Pipe{{
				OP: definition.OP_ASSIGNMENT,
				Params: definition.JPLInstructionParams{
					Pipe:       definition.Pipe{{OP: definition.OP_VARIABLE, Params: definition.JPLInstructionParams{Name: name}}},
					Operations: operations,
					Assignment: opAssignment,
				},
			}},
		},
	}}, nil
}

// Parse variable access at i
func opValueAccess(src string, i int, c *ParserContext) (n int, result definition.Pipe, err error) {
	n = i

	var operations []definition.JPLOperation

	var ops definition.Pipe
	iM, isM, err := matchWord(src, n, c, matchOptions{Phrase: "."})
	if err != nil {
		return 0, nil, err
	}
	if !isM {
		if n, ops, err = opObjectConstructor(src, n, c); err != nil {
			return 0, nil, err
		}
	} else {
		n = iM
		ops = nil

		iV, isV, name, _, err := safeVariable(src, n, c)
		if err != nil {
			return 0, nil, err
		}
		if isV {
			n = iV

			var optional bool
			iM, isM, err = matchWord(src, n, c, matchOptions{Phrase: "?", NotBeforeSet: "?="})
			if err != nil {
				return 0, nil, err
			}
			if isM {
				n = iM
				optional = true
			}

			operations = append(operations, definition.JPLOperation{
				OP:     definition.OPA_FIELD,
				Params: definition.JPLOperationParams{Pipe: definition.Pipe{{OP: definition.OP_STRING, Params: definition.JPLInstructionParams{Value: name}}}, Optional: optional},
			})
		}
	}

	iAc, isAc, operationsAc, canAssign, err := parseAccess(src, n, c, accessOptions{Identity: len(ops) == 0 && len(operations) == 0})
	if err != nil {
		return 0, nil, err
	}
	if isAc {
		n = iAc
		operations = append(operations, operationsAc...)
	} else {
		canAssign = len(operations) > 0
	}

	if len(operations) == 0 {
		return n, ops, nil
	}

	if !canAssign {
		return n, definition.Pipe{{OP: definition.OP_ACCESS, Params: definition.JPLInstructionParams{Pipe: ops, Operations: operations}}}, nil
	}

	var opAssignment *definition.JPLAssignment
	iAs, isAs, opAssignment, err := parseAssignment(src, n, c)
	if err != nil {
		return 0, nil, err
	}
	if !isAs {
		return n, definition.Pipe{{OP: definition.OP_ACCESS, Params: definition.JPLInstructionParams{Pipe: ops, Operations: operations}}}, nil
	}
	n = iAs

	return n, definition.Pipe{{OP: definition.OP_ASSIGNMENT, Params: definition.JPLInstructionParams{Pipe: ops, Operations: operations, Assignment: opAssignment}}}, nil
}

// Parse object constructor at i
func opObjectConstructor(src string, i int, c *ParserContext) (n int, result definition.Pipe, err error) {
	n = i

	iM, isM, err := matchWord(src, n, c, matchOptions{Phrase: "{"})
	if err != nil {
		return 0, nil, err
	}
	if !isM {
		return opArrayConstructor(src, n, c)
	}
	n = iM

	var fields []definition.JPLField

	iM, isM, err = matchWord(src, n, c, matchOptions{Phrase: "}"})
	if err != nil {
		return 0, nil, err
	}
	if isM {
		n = iM
	} else {
		for {
			iM, isM, err = matchWord(src, n, c, matchOptions{Phrase: "("})
			if err != nil {
				return 0, nil, err
			}
			if isM {
				n = iM

				var opsKey definition.Pipe
				if n, opsKey, err = opPipe(src, n, c); err != nil {
					return 0, nil, err
				}

				iM, isM, err := matchWord(src, n, c, matchOptions{Phrase: ")"})
				if err != nil {
					return 0, nil, err
				}
				n = iM
				if !isM {
					return 0, nil, errorUnexpectedToken(src, n, c, errorOptions{Operator: "object", Message: "expected ')'"})
				}

				var optional bool
				if n, optional, err = matchWord(src, n, c, matchOptions{Phrase: "?"}); err != nil {
					return 0, nil, err
				}

				iM, isM, err = matchWord(src, n, c, matchOptions{Phrase: ":"})
				if err != nil {
					return 0, nil, err
				}
				n = iM
				if !isM {
					return 0, nil, errorUnexpectedToken(src, n, c, errorOptions{Operator: "object", Message: "expected ':'"})
				}

				var opsValue definition.Pipe
				if n, opsValue, err = opSubPipe(src, n, c); err != nil {
					return 0, nil, err
				}

				fields = append(fields, definition.JPLField{Key: opsKey, Value: opsValue, Optional: optional})

				iM, isM, err = matchWord(src, n, c, matchOptions{Phrase: "}"})
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
					return 0, nil, errorUnexpectedToken(src, n, c, errorOptions{
						Operator: "object",
						Message:  "expected ',' or '}'",
					})
				}

				continue
			}

			iS, isS, opsKey, err := parseString(src, n, c)
			if err != nil {
				return 0, nil, err
			}
			if isS {
				n = iS

				iM, isM, err = matchWord(src, n, c, matchOptions{Phrase: ":"})
				if err != nil {
					return 0, nil, err
				}
				n = iM
				if !isM {
					return 0, nil, errorUnexpectedToken(src, n, c, errorOptions{Operator: "object", Message: "expected ':'"})
				}

				var opsValue definition.Pipe
				if n, opsValue, err = opSubPipe(src, n, c); err != nil {
					return 0, nil, err
				}

				fields = append(fields, definition.JPLField{Key: opsKey, Value: opsValue, Optional: false})

				iM, isM, err = matchWord(src, n, c, matchOptions{Phrase: "}"})
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
					return 0, nil, errorUnexpectedToken(src, n, c, errorOptions{
						Operator: "object",
						Message:  "expected ',' or '}'",
					})
				}

				continue
			}

			iV, isV, name, _, err := safeVariable(src, n, c)
			if err != nil {
				return 0, nil, err
			}
			if isV {
				n = iV

				var opsValue definition.Pipe
				iM, isM, err = matchWord(src, n, c, matchOptions{Phrase: ":"})
				if err != nil {
					return 0, nil, err
				}
				if isM {
					n = iM

					if n, opsValue, err = opSubPipe(src, n, c); err != nil {
						return 0, nil, err
					}
				} else {
					var optional bool
					if n, optional, err = matchWord(src, n, c, matchOptions{Phrase: "?"}); err != nil {
						return 0, nil, err
					}

					if optional {
						opsValue = definition.Pipe{{
							OP: definition.OP_TRY,
							Params: definition.JPLInstructionParams{
								Try:   definition.Pipe{{OP: definition.OP_VARIABLE, Params: definition.JPLInstructionParams{Name: name}}},
								Catch: definition.Pipe{{OP: definition.OP_VOID}},
							},
						}}
					} else {
						opsValue = definition.Pipe{{OP: definition.OP_VARIABLE, Params: definition.JPLInstructionParams{Name: name}}}
					}
				}

				fields = append(fields, definition.JPLField{
					Key:      definition.Pipe{{OP: definition.OP_STRING, Params: definition.JPLInstructionParams{Value: name}}},
					Value:    opsValue,
					Optional: false,
				})

				iM, isM, err = matchWord(src, n, c, matchOptions{Phrase: "}"})
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
					return 0, nil, errorUnexpectedToken(src, n, c, errorOptions{
						Operator: "object",
						Message:  "expected ',' or '}'",
					})
				}

				continue
			}

			return 0, nil, errorUnexpectedToken(src, n, c, errorOptions{
				Operator: "object",
				Message:  "expected field declaration",
			})
		}
	}

	return n, definition.Pipe{{OP: definition.OP_OBJECT_CONSTRUCTOR, Params: definition.JPLInstructionParams{Fields: fields}}}, nil
}

// Parse array constructor at i
func opArrayConstructor(src string, i int, c *ParserContext) (n int, result definition.Pipe, err error) {
	n = i

	iM, isM, err := matchWord(src, n, c, matchOptions{Phrase: "["})
	if err != nil {
		return 0, nil, err
	}
	n = iM
	if !isM {
		return opStringLiteral(src, n, c)
	}

	var ops definition.Pipe

	iM, isM, err = matchWord(src, n, c, matchOptions{Phrase: "]"})
	if err != nil {
		return 0, nil, err
	}
	if isM {
		n = iM

		ops = definition.Pipe{{OP: definition.OP_VOID}}
	} else {
		if n, ops, err = opPipe(src, n, c); err != nil {
			return 0, nil, err
		}

		iM, isM, err = matchWord(src, n, c, matchOptions{Phrase: "]"})
		if err != nil {
			return 0, nil, err
		}
		n = iM
		if !isM {
			return 0, nil, errorUnexpectedToken(src, n, c, errorOptions{Operator: "array", Message: "expected ']'"})
		}
	}

	return n, definition.Pipe{{OP: definition.OP_ARRAY_CONSTRUCTOR, Params: definition.JPLInstructionParams{Pipe: ops}}}, nil
}

// Parse string literal at i
func opStringLiteral(src string, i int, c *ParserContext) (n int, result definition.Pipe, err error) {
	n = i

	iS, isS, ops, err := parseString(src, n, c)
	if err != nil {
		return 0, nil, err
	}
	if !isS {
		return opGroup(src, n, c)
	}
	n = iS

	return n, ops, nil
}

// Parse group at i
func opGroup(src string, i int, c *ParserContext) (n int, result definition.Pipe, err error) {
	n = i

	iM, isM, err := matchWord(src, n, c, matchOptions{Phrase: "("})
	if err != nil {
		return 0, nil, err
	}
	n = iM
	if !isM {
		return 0, nil, errorUnexpectedToken(src, n, c, errorOptions{})
	}

	var ops definition.Pipe
	if n, ops, err = opPipe(src, n, c); err != nil {
		return 0, nil, err
	}

	iM, isM, err = matchWord(src, n, c, matchOptions{Phrase: ")"})
	if err != nil {
		return 0, nil, err
	}
	n = iM
	if !isM {
		return 0, nil, errorUnexpectedToken(src, n, c, errorOptions{Operator: "group", Message: "expected ')'"})
	}

	return n, ops, nil
}
