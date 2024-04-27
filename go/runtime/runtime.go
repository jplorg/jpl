package runtime

import (
	"github.com/2manyvcos/jpl/go/definition"
	"github.com/2manyvcos/jpl/go/jpl"
	"github.com/2manyvcos/jpl/go/library"
)

var defaultOptions = jpl.JPLRuntimeOptions{}

func NewRuntime(program jpl.JPLProgram, options *jpl.JPLRuntimeConfig) jpl.JPLRuntime {
	if options == nil {
		options = new(jpl.JPLRuntimeConfig)
	}

	return &runtime{
		options: jpl.ApplyRuntimeDefaults(options.Runtime, defaultOptions),

		program: program,
	}
}

type runtime struct {
	options jpl.JPLRuntimeOptions

	program jpl.JPLProgram
}

func (r *runtime) Options() jpl.JPLRuntimeOptions {
	return r.options
}

func (r *runtime) Program() jpl.JPLProgram {
	return r.program
}

func (r *runtime) CreateScope(presets *jpl.JPLRuntimeScopeConfig) jpl.JPLRuntimeScope {
	return library.NewRuntimeScope(presets)
}

func (r *runtime) Execute(inputs []any) ([]any, jpl.JPLError) {
	varEntries, err := library.MuxOne([][]*library.ObjectEntry[any]{library.ObjectEntries(r.Options().Vars)}, jpl.IOMuxerFunc[*library.ObjectEntry[any], *library.ObjectEntry[any]](func(args ...*library.ObjectEntry[any]) (result *library.ObjectEntry[any], err jpl.JPLError) {
		result = args[0]
		result.Value, err = library.NormalizeValue(result.Value)
		return
	}))
	if err != nil {
		return nil, err
	}
	scope := r.CreateScope(&jpl.JPLRuntimeScopeConfig{
		Vars: library.ObjectFromEntries(varEntries),
	})

	defer scope.Signal().Exit()

	return r.ExecuteInstructions(
		r.Program().Definition().Instructions,
		inputs,
		scope,
		r.Options().AdjustResult,
	)
}

func (r *runtime) ExecuteInstructions(instructions definition.Pipe, inputs []any, scope jpl.JPLRuntimeScope, next jpl.JPLScopedPiper) ([]any, jpl.JPLError) {
	if next == nil {
		next = jpl.JPLScopedPiperFunc(func(output any, _ jpl.JPLRuntimeScope) ([]any, jpl.JPLError) {
			return []any{output}, nil
		})
	}

	var iter func(from int, input any, currentScope jpl.JPLRuntimeScope) ([]any, jpl.JPLError)
	iter = func(from int, input any, currentScope jpl.JPLRuntimeScope) ([]any, jpl.JPLError) {
		if err := currentScope.Signal().CheckHealth(); err != nil {
			return nil, err
		}

		if from >= len(instructions) {
			return next.Pipe(input, currentScope)
		}

		instruction := instructions[from]
		operator := r.Program().OPs()[instruction.OP]
		if operator == nil {
			return nil, library.NewFatalError("invalid OP '" + string(instruction.OP) + "'")
		}

		return operator.OP(r, input, instruction.Params, currentScope, jpl.JPLScopedPiperFunc(func(output any, nextScope jpl.JPLRuntimeScope) ([]any, jpl.JPLError) {
			return iter(from+1, output, nextScope)
		}))
	}

	return library.MuxAll([][]any{inputs}, jpl.IOMuxerFunc[any, []any](func(args ...any) ([]any, jpl.JPLError) {
		return iter(0, args[0], scope)
	}))
}

func (r *runtime) OP(op definition.JPLOP, params jpl.JPLInstructionParams, inputs []any, scope jpl.JPLRuntimeScope, next jpl.JPLScopedPiper) ([]any, jpl.JPLError) {
	operator := r.Program().OPs()[op]
	if operator == nil {
		return nil, library.NewFatalError("invalid OP '" + string(op) + "'")
	}

	opParams, err := operator.Map(r, params)
	if err != nil {
		return nil, err
	}
	return library.MuxAll([][]any{inputs}, jpl.IOMuxerFunc[any, []any](func(args ...any) ([]any, jpl.JPLError) {
		return operator.OP(r, args[0], opParams, scope, next)
	}))
}
