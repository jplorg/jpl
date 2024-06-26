package program

import (
	"github.com/jplorg/jpl/go/definition"
	"github.com/jplorg/jpl/go/jpl"
	"github.com/jplorg/jpl/go/library"
)

type opVariable struct{}

// { name: string }
func (opVariable) OP(runtime jpl.JPLRuntime, input any, params definition.JPLInstructionParams, scope jpl.JPLRuntimeScope, next jpl.JPLScopedPiper) ([]any, jpl.JPLError) {
	value, ok := scope.Vars()[params.Name]
	if !ok {
		return nil, library.ThrowAny(library.NewReferenceError("%s is not defined", params.Name))
	}
	return next.Pipe(value, scope)
}

// { name: string }
func (opVariable) Map(runtime jpl.JPLRuntime, params jpl.JPLInstructionParams) (result definition.JPLInstructionParams, err jpl.JPLError) {
	return definition.JPLInstructionParams{
		Name: params.Name,
	}, nil
}
