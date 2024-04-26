package program

import (
	"github.com/2manyvcos/jpl/go/definition"
	"github.com/2manyvcos/jpl/go/jpl"
	"github.com/2manyvcos/jpl/go/library"
)

type opNegation struct{}

// {}
func (opNegation) OP(runtime jpl.JPLRuntime, input any, params definition.JPLInstructionParams, scope jpl.JPLRuntimeScope, next jpl.JPLScopedPiper) ([]any, jpl.JPLError) {
	alteredValue, err := library.AlterValue(input, jpl.JPLModifierFunc(func(value any) (any, jpl.JPLError) {
		t, err := library.Type(value)
		if err != nil {
			return nil, err
		}
		switch t {
		case jpl.JPLT_NUMBER:
			return -value.(float64), nil

		default:
		}

		return nil, library.ThrowAny(library.NewTypeError("%s (%*<100v) cannot be negated", string(t), value))
	}))
	if err != nil {
		return nil, err
	}
	return next.Pipe(alteredValue, scope)
}

// {}
func (opNegation) Map(runtime jpl.JPLRuntime, params jpl.JPLInstructionParams) (result definition.JPLInstructionParams, err jpl.JPLError) {
	return
}
