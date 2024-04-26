package interpreter

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/2manyvcos/jpl/go/jpl"
	"github.com/2manyvcos/jpl/go/library"
)

const setAZ = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const setDigit = "0123456789"
const setVarFirst = setAZ + "_$"
const setVarRest = setVarFirst + setDigit
const setVarAll = setVarRest
const setWhitespace = " \r\n\t"
const setHex = setDigit + "abcdefABCDEF"

// Walk whitespace at i
func walkWhitespace(src string, i int, c *ParserContext) (n int, is bool, err jpl.JPLSyntaxError) {
	n = i

	for {
		if iM, isM := match(src, n, c, matchOptions{Phrase: "#"}); isM {
			n = iM
			for {
				if iSet, isSet, _ := matchSet(src, n, c, matchSetOptions{Set: "\r\n"}); isSet {
					n = iSet
					break
				}

				if iEnd, isEnd := eot(src, n, c); isEnd {
					n = iEnd
					break
				}

				n += 1
			}
			is = true
			continue
		}

		if iM, isM := match(src, n, c, matchOptions{Phrase: "/*"}); isM {
			n = iM
			for {
				if iM, isM := match(src, n, c, matchOptions{Phrase: "*/"}); isM {
					n = iM
					break
				}

				if iEnd, isEnd := eot(src, n, c); isEnd {
					n = iEnd
					return 0, false, errorUnexpectedToken(src, n, c, errorOptions{Operator: "comment", Message: "incomplete comment"})
				}

				n += 1
			}
			is = true
			continue
		}

		if iM, isM := match(src, n, c, matchOptions{Phrase: "\u00a0"}); isM {
			n = iM
			is = true
			continue
		}

		iSet, isSet, _ := matchSet(src, n, c, matchSetOptions{Set: setWhitespace})
		n = iSet
		if !isSet {
			break
		}
		is = true
	}

	return n, is, nil
}

type matchOptions struct {
	Phrase       string
	NotBeforeSet string
	SpaceBefore  bool
	SpaceAfter   bool
}

// Check if src matches phrase at i.
// If notBeforeSet is provided, the phrase is only considered to be matching if it is not immediately followed by one of the symbols of the set.
// If spaceBefore is true, the phrase is only considered to be matching if there is space (not to be confused with whitespace) before i (or i is 0).
// If spaceAfter is true, the phrase is only considered to be matching if there is space (not to be confused with whitespace) after the phrase (or the end is reached).
// The returned index is positioned directly after the phrase or at the first unmatched character, even if spaceAfter is true.
// Space is defined as anything that cannot occur in a variable name.
func match(src string, i int, c *ParserContext, options matchOptions) (n int, is bool) {
	n = i

	if options.SpaceBefore && i > 0 {
		if _, isSet, _ := matchSet(src, n-1, c, matchSetOptions{Set: setVarAll, Exclusive: true}); !isSet {
			return n, false
		}
	}

	for j := 0; j < len(options.Phrase); j, n = j+1, n+1 {
		if iEnd, isEnd := eot(src, n, c); isEnd || src[n] != options.Phrase[j] {
			return iEnd, false
		}
	}

	if _, isEnd := eot(src, n, c); !isEnd {
		if options.NotBeforeSet != "" {
			if _, isSet, _ := matchSet(src, n, c, matchSetOptions{Set: options.NotBeforeSet}); isSet {
				return n, false
			}
		}

		if options.SpaceAfter {
			if _, isSet, _ := matchSet(src, n, c, matchSetOptions{Set: setVarAll, Exclusive: true}); !isSet {
				return n, false
			}
		}
	}

	return n, true
}

// The same as match, but also walk whitespace after the match if present
func matchWord(src string, i int, c *ParserContext, options matchOptions) (n int, is bool, err jpl.JPLSyntaxError) {
	n = i

	iM, isM := match(src, n, c, options)
	if !isM {
		return iM, false, nil
	}
	n = iM

	if n, _, err = walkWhitespace(src, n, c); err != nil {
		return 0, false, err
	}

	return n, true, nil
}

// Check if i is a hexadecimal value
func hex(src string, i int, c *ParserContext) (n int, is bool, value string) {
	return matchSet(src, i, c, matchSetOptions{Set: setHex})
}

type matchSetOptions struct {
	Set       string
	Exclusive bool
}

// Check if i is one of the chars in set.
// If exclusive is true, i must not be one of the chars in set for the check to succeed.
func matchSet(src string, i int, c *ParserContext, options matchSetOptions) (n int, is bool, value string) {
	if iEnd, isEnd := eot(src, i, c); isEnd {
		return iEnd, false, ""
	}

	char := string(src[i])
	if strings.Contains(options.Set, char) == !options.Exclusive {
		return i + 1, true, char
	}

	return i, false, ""
}

// Parse variable selector at i
func variable(src string, i int, c *ParserContext) (n int, is bool, value string, err jpl.JPLSyntaxError) {
	n = i

	for {
		set := setVarRest
		if !is {
			set = setVarFirst
		}
		iSet, isSet, valueSet := matchSet(src, n, c, matchSetOptions{Set: set})
		n = iSet
		if !isSet {
			break
		}
		is = true
		value += valueSet
	}

	if is {
		if n, _, err = walkWhitespace(src, n, c); err != nil {
			return 0, false, "", err
		}
	}

	return n, is, value, nil
}

// The same as variable, but also check that the result is not a reserved term
func safeVariable(src string, i int, c *ParserContext) (n int, is bool, value string, reserved bool, err jpl.JPLSyntaxError) {
	n, is, value, err = variable(src, i, c)

	if !is {
		return
	}

	switch value {
	case "and", "catch", "func", "elif", "else", "end", "false", "if", "not", "null", "or", "then", "true", "try":
		return i, false, value, true, nil

	default:
		return
	}
}

// Check if i is at the end of src
func eot(src string, i int, _ *ParserContext) (n int, is bool) {
	if l := len(src); i >= l {
		return l, true
	}
	return i, false
}

// Get (zero based) line and column for i
func whereIs(src string, i int, _ *ParserContext) (n int, line int, column int) {
	lines := regexp.MustCompile(`\r?\n|\r`).Split(src[0:i], -1)
	line = len(lines) - 1
	currentLine := lines[line]
	return i, line, len(currentLine)
}

type highlightOptions struct {
	Area int
}

// Get a descriptive text highlighting i
func highlightLocation(src string, i int, _ *ParserContext, options highlightOptions) (n int, value string) {
	area := options.Area
	if area <= 0 {
		area = 25
	}

	l := len(src)
	s := max(min(i, l-1-area), area) - area
	e := min(s+area+1+area, l)
	view := strings.ReplaceAll(regexp.MustCompile(`\r?\n|\r`).ReplaceAllString(src[s:e], "⏎"), "\t", "→")

	prefix := " > "
	if s > 0 {
		prefix += "…"
	}

	suffix := ""
	if e < l {
		suffix = "…"
	}

	description := prefix + view + suffix + "\n" + strings.Repeat(" ", len(prefix)+(i-s)) + "^ here"

	return i, description
}

type errorOptions struct {
	Operator string
	Message  string
}

// Throw an error caused by an unexpected token at i
func errorUnexpectedToken(src string, i int, c *ParserContext, options errorOptions) jpl.JPLSyntaxError {
	var errorMessage string
	if _, isEnd := eot(src, i, c); isEnd {
		errorMessage = "unexpected EOT"
	} else {
		_, line, column := whereIs(src, i, c)
		errorMessage = fmt.Sprintf("unexpected token '%s' at line %v, column %v", string(src[i]), line+1, column+1)
	}
	if options.Operator != "" {
		errorMessage += " while parsing " + options.Operator
	}
	if options.Message != "" {
		errorMessage += ": " + options.Message
	}
	_, description := highlightLocation(src, i, c, highlightOptions{})
	errorMessage += "\n" + description
	return library.NewSyntaxError(errorMessage)
}

// Throw an error caused by a generic parser error at i
func errorGeneric(src string, i int, c *ParserContext, options errorOptions) jpl.JPLSyntaxError {
	errorMessage := "error"
	if options.Operator != "" {
		errorMessage += " while parsing " + options.Operator
	}
	if options.Message != "" {
		errorMessage += ": " + options.Message
	}
	_, description := highlightLocation(src, i, c, highlightOptions{})
	errorMessage += "\n" + description
	return library.NewSyntaxError(errorMessage)
}
