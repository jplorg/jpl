package runtime

import (
	"github.com/2manyvcos/jpl/go/config"
	"github.com/2manyvcos/jpl/go/definition"
	"github.com/2manyvcos/jpl/go/jpl"
	"github.com/2manyvcos/jpl/go/library"
)

var defaultOptions = config.JPLRuntimeOptions{}

func NewRuntime(program jpl.JPLProgram, options *config.JPLRuntimeConfig) jpl.JPLRuntime {
	if options == nil {
		options = new(config.JPLRuntimeConfig)
	}

	return &runtime{
		options: config.ApplyRuntimeDefaults(options.Runtime, defaultOptions),

		program: program,
	}
}

type runtime struct {
	options config.JPLRuntimeOptions

	program jpl.JPLProgram
}

func (r *runtime) Options() config.JPLRuntimeOptions {
	return r.options
}

func (r *runtime) Program() jpl.JPLProgram {
	return r.program
}

func (r *runtime) CreateScope(presets *library.RuntimeScopeConfig) library.RuntimeScope {
	return library.NewRuntimeScope(presets)
}

func (r *runtime) Execute(inputs []any) ([]any, error) {
	panic("TODO:")
}

func (r *runtime) ExecuteInstructions(instructions definition.Pipe, inputs []any, scope library.RuntimeScope, next library.JPLPiper) ([]any, error) {
	panic("TODO:")
}

func (r *runtime) OP(op definition.JPLOP, params map[string]any, inputs []any, scope library.RuntimeScope, next library.JPLPiper) ([]any, error) {
	panic("TODO:")
}

func (r *runtime) NormalizeValue(value any) (any, error) {
	panic("TODO:")
}

func (r *runtime) NormalizeValues(values []any, name string) ([]any, error) {
	panic("TODO:")
}

func (r *runtime) UnwrapValue(value any) (any, error) {
	panic("TODO:")
}

func (r *runtime) UnwrapValues(values []any, name string) ([]any, error) {
	panic("TODO:")
}

func (r *runtime) StripValue(value any) (any, error) {
	panic("TODO:")
}

func (r *runtime) StripValues(value []any) ([]any, error) {
	panic("TODO:")
}

func (r *runtime) AlterValue(value any, updater library.JPLModifier) (any, error) {
	panic("TODO:")
}

func (r *runtime) Type(value any) (library.JPLDataType, error) {
	panic("TODO:")
}

func (r *runtime) AssertType(value any, assertedType library.JPLDataType) (any, error) {
	panic("TODO:")
}

func (r *runtime) Truthy(value any) (bool, error) {
	panic("TODO:")
}

func (r *runtime) Compare(a, b any) (int, error) {
	panic("TODO:")
}

func (r *runtime) CompareStrings(a, b any) (int, error) {
	panic("TODO:")
}

func (r *runtime) CompareArrays(a, b any) (int, error) {
	panic("TODO:")
}

func (r *runtime) CompareObjects(a, b any) (int, error) {
	panic("TODO:")
}

func (r *runtime) Equals(a, b any) (bool, error) {
	panic("TODO:")
}

func (r *runtime) Merge(a, b any) (any, error) {
	panic("TODO:")
}

func (r *runtime) StringifyJSON(value any, unescapeString bool) (string, error) {
	panic("TODO:")
}

func (r *runtime) StripJSON(value any) (any, error) {
	panic("TODO:")
}

func (r *runtime) Mux(args [][]any, cb library.JPLMuxer) error {
	panic("TODO:")
}

func (r *runtime) MuxOne(args [][]any, cb library.JPLMuxerOne) ([]any, error) {
	panic("TODO:")
}

func (r *runtime) MuxAll(args [][]any, cb library.JPLMuxerAll) ([]any, error) {
	panic("TODO:")
}
