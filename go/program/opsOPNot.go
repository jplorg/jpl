package program

import (
	"github.com/jplorg/jpl/go/v2/definition"
	"github.com/jplorg/jpl/go/v2/jpl"
	"github.com/jplorg/jpl/go/v2/library"
)

type opNot struct{}

// {}
func (opNot) OP(runtime jpl.JPLRuntime, input any, params definition.JPLInstructionParams, scope jpl.JPLRuntimeScope, next jpl.JPLScopedPiper) ([]any, jpl.JPLError) {
	truthy, err := library.Truthy(input)
	if err != nil {
		return nil, err
	}
	return next.Pipe(!truthy, scope)
}

// {}
func (opNot) Map(runtime jpl.JPLRuntime, params jpl.JPLInstructionParams) (result definition.JPLInstructionParams, err jpl.JPLError) {
	return
}
