package main

import (
	"encoding/json"
	"fmt"

	"github.com/2manyvcos/jpl/go/interpreter"
)

func main() {
	interpreter := interpreter.NewInterpreter(nil)

	instructions, err := interpreter.ParseInstructions(". + 1 | 'test\\uD83C\\uDDE6\\uD83C\\uDDE8\\uD83D\\uDE0A' + .")
	if err != nil {
		panic(err)
	}

	output, err := json.MarshalIndent(instructions, "", "  ")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(output))
	return

	panic("TODO: CLI REPL")
}
