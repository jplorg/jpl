package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"slices"
	"strings"
	"time"
	"unicode"

	"os"
	"path/filepath"

	"github.com/chzyer/readline"
	gojpl "github.com/jplorg/jpl/go"
	"github.com/jplorg/jpl/go/jpl"
	"github.com/jplorg/jpl/go/library"
)

var replKeys = []rune{':', '!'}
var defaultReplKey = replKeys[0]

const defaultPrompt = "> "
const multilinePrompt = "… "

var multilineInput *string

var muted bool
var keep bool
var inputs []any
var measureTime bool

func main() {
	muted = !readline.IsTerminal(int(os.Stdin.Fd()))

	var historyFile string
	if homeDir, err := os.UserHomeDir(); err == nil {
		historyFile = filepath.Join(homeDir, ".jpl_repl_history")
	}

	rl, err := readline.NewEx(&readline.Config{
		Stdin:           os.Stdin,
		Stdout:          os.Stdout,
		Stderr:          os.Stderr,
		HistoryFile:     historyFile,
		HistoryLimit:    50,
		Prompt:          defaultPrompt,
		InterruptPrompt: "^C",
	})
	if err != nil {
		panic(err)
	}
	defer rl.Close()
	rl.CaptureExitSignal()

	if !muted {
		fmt.Println("Welcome to JPL.")
		fmt.Printf("Type \"%ch\" for more information.\n\n", defaultReplKey)
	}

	gojpl.Options.Runtime.Vars["exit"] = library.NativeFunction(func(runtime jpl.JPLRuntime, input any, args ...any) ([]any, error) {
		rl.Close()
		os.Exit(0)
		return nil, nil
	})
	gojpl.Options.Runtime.Vars["clear"] = library.NativeFunction(func(runtime jpl.JPLRuntime, input any, args ...any) ([]any, error) {
		readline.ClearScreen(rl)
		return nil, nil
	})

	for {
		line, err := rl.Readline()
		if err == readline.ErrInterrupt {
			if multilineInput != nil {
				multilineInput = nil
				rl.SetPrompt(defaultPrompt)
				continue
			}
			if len(line) == 0 {
				break
			} else {
				fmt.Fprintf(rl, "To exit, press Ctrl+C again or type %ce\n", defaultReplKey)
				continue
			}
		} else if err == io.EOF {
			break
		}

		handle(line, rl)
	}
}

func handle(input string, rl *readline.Instance) {
	if !keep || len(inputs) == 0 {
		inputs = []any{nil}
	}

	fullLine := input
	if multilineInput != nil {
		fullLine = *multilineInput + input
	}
	line := fullLine
	t := []rune(strings.TrimLeftFunc(line, unicode.IsSpace))

	if len(t) == 0 {
		return
	}

	if slices.Contains(replKeys, t[0]) {
		var command rune
		if len(t) < 2 {
			command = ' '
			line = ""
		} else {
			command = unicode.ToLower(t[1])
			line = string(t[2:])
		}

		switch command {
		case 'h':
			printHelp(rl)

		case 'e', 'q':
			rl.Close()
			os.Exit(0)

		case 'c':
			readline.ClearScreen(rl)

		case 'k':
			keep = parseBool(line, !keep, keep, "keep", rl)

		case 't':
			measureTime = parseBool(line, !measureTime, measureTime, "time", rl)

		case 'i':
			// reset prompt after potential previous multiline input
			multilineInput = nil
			rl.SetPrompt(defaultPrompt)
			program, err := gojpl.Parse(line, nil)
			if err != nil {
				if err, ok := err.(jpl.JPLSyntaxError); ok && err.At() >= len(err.Src()) {
					// program is incomplete -> request additional input
					ml := fullLine + "\n"
					multilineInput = &ml
					rl.SetPrompt(multilinePrompt)
					return
				}
				printError(rl, err)
				return
			}
			if json, err := marshalJSONIndent(program.Definition()); err != nil {
				fmt.Fprintf(rl, "Error: %s\n", err)
				return
			} else {
				fmt.Fprintln(rl, string(json))
			}

		case ' ':
			fmt.Fprintf(rl, "Error: missing REPL command\n\n")
			printHelp(rl)

		default:
			fmt.Fprintf(rl, "Error: unrecognized REPL command %c%c\n\n", defaultReplKey, command)
			printHelp(rl)
		}
	} else {
		// reset prompt after potential previous multiline input
		multilineInput = nil
		rl.SetPrompt(defaultPrompt)
		program, err := gojpl.Parse(line, nil)
		if err != nil {
			if err, ok := err.(jpl.JPLSyntaxError); ok && err.At() >= len(err.Src()) {
				// program is incomplete -> request additional input
				ml := fullLine + "\n"
				multilineInput = &ml
				rl.SetPrompt(multilinePrompt)
				return
			}
			printError(rl, err)
			return
		}
		var before, diff int64
		if measureTime {
			before = time.Now().UnixMilli()
		}
		nextInputs, err := program.Run(inputs, nil)
		if err != nil {
			printError(rl, err)
			return
		}
		inputs = nextInputs
		if measureTime {
			diff = time.Now().UnixMilli() - before
		}
		outputs := make([]string, len(inputs))
		for i, output := range inputs {
			json, err := marshalJSONIndent(output)
			if err != nil {
				fmt.Fprintf(rl, "Error: %s\n", err)
				return
			}
			outputs[i] = string(json)
		}
		fmt.Fprintln(rl, strings.Join(outputs, ", "))
		if measureTime {
			fmt.Fprintf(rl, " -> took %vs\n", float64(diff)/1000)
		}
	}
}

func parseBool(input string, defaultValue, fallbackValue bool, label string, rl *readline.Instance) bool {
	b := strings.ToLower(strings.TrimSpace(input))
	var v bool
	if b == "" {
		v = defaultValue
	} else if b == "on" || strings.HasPrefix("true", b) || strings.HasPrefix("yes", b) || strings.HasPrefix("enabled", b) {
		v = true
	} else if b == "off" || strings.HasPrefix("false", b) || strings.HasPrefix("no", b) || strings.HasPrefix("disabled", b) {
		v = false
	} else {
		fmt.Fprintf(rl, "Error: invalid boolean %s\n", b)
		return fallbackValue
	}
	if v {
		fmt.Fprintf(rl, " -> %s on\n", label)
	} else {
		fmt.Fprintf(rl, " -> %s off\n", label)
	}
	return v
}

func printBool(value bool) string {
	if value {
		return "boolean (on)"
	}
	return "boolean (off)"
}

type command struct {
	Command     rune
	Args        string
	Description string
}

func printHelp(rl *readline.Instance) {
	commands := []command{
		{'c', "", "Clear the console screen"},
		{'e', "", "Exit the REPL"},
		{'h', "", "Print this help message"},
		{'i', "program", "Interpret the specified program without executing it"},
		{'k', printBool(keep), "Set whether program output should be kept as input for the next program"},
		{'t', printBool(measureTime), "Set whether execution time should be measured"},
		{'q', "", "Exit the REPL"},
	}
	var aLen int
	for _, c := range commands {
		aLen = max(aLen, len(c.Args))
	}

	fmt.Fprintf(rl, "JPL REPL reference\n\n")
	fmt.Fprintf(rl, "The following synonymous tokens may be used to precede a command: %s\n\n", string(replKeys))

	for _, c := range commands {
		fmt.Fprintf(rl, "%c%c %s%s%s\n", defaultReplKey, c.Command, c.Args, strings.Repeat(" ", aLen-len(c.Args)+3), c.Description)
	}

	fmt.Fprintln(rl, "\nPress Ctrl+C to abort current expression, Ctrl+D to exit the REPL")
}

// Marshal the specified value to JSON with indentation
func marshalJSONIndent(value any) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	err := encoder.Encode(value)
	return bytes.TrimRight(buffer.Bytes(), "\n"), err
}

// Print an error to the console
func printError(rl *readline.Instance, err error) {
	if err, ok := err.(jpl.JPLError); ok {
		name := err.JPLErrorName()
		if name == "" {
			name = "JPLError"
		}
		fmt.Fprintf(rl, "%s: %s\n", name, err.JPLErrorMessage())
	} else {
		fmt.Fprintf(rl, "Error: %s\n", err)
	}
}
