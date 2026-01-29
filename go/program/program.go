package program

import (
	"fmt"
	"regexp"
	"strconv"

	definition "github.com/jplorg/jpl/go/v2/definition"
	jpl "github.com/jplorg/jpl/go/v2/jpl"
	library "github.com/jplorg/jpl/go/v2/library"
	runtime "github.com/jplorg/jpl/go/v2/runtime"
)

var defaultOptions = jpl.JPLProgramOptions{}

var versionRegex = regexp.MustCompile(`^(?P<Major>\d+)\.(?P<Minor>\d+)$`)

func validateDefinition(programDefinition definition.JPLDefinition) jpl.JPLError {
	parts := versionRegex.FindStringSubmatch(programDefinition.Version)
	var major, minor string
	var majorV, minorV int
	if len(parts) > 2 {
		part := parts[1]
		if v, err := strconv.Atoi(part); err == nil {
			major = part
			majorV = v
		}
		part = parts[2]
		if v, err := strconv.Atoi(part); err == nil {
			minor = part
			minorV = v
		}
	}

	if major == "" || minor == "" {
		return library.NewFatalError("invalid program definition")
	}

	if majorV != definition.DEFINITION_VERSION_MAJOR || minorV > definition.DEFINITION_VERSION_MINOR {
		return library.NewFatalError(fmt.Sprintf("unsupported program definition v%s.%s - this version of JPL only supports v%v (up to v%s)", major, minor, definition.DEFINITION_VERSION_MAJOR, definition.DEFINITION_VERSION))
	}

	return nil
}

func NewProgram(programDefinition definition.JPLDefinition, options *jpl.JPLProgramConfig) (jpl.JPLProgram, jpl.JPLError) {
	if err := validateDefinition(programDefinition); err != nil {
		return nil, err
	}

	if options == nil {
		options = new(jpl.JPLProgramConfig)
	}

	return &program{
		options:        jpl.ApplyProgramDefaults(options.Program, defaultOptions),
		runtimeOptions: options.Runtime,

		definition: programDefinition,
		ops:        ops,
	}, nil
}

type program struct {
	options        jpl.JPLProgramOptions
	runtimeOptions jpl.JPLRuntimeOptions

	definition definition.JPLDefinition
	ops        map[definition.JPLOP]jpl.JPLOPHandler
}

func (p *program) Options() jpl.JPLProgramOptions {
	return p.options
}

func (p *program) Definition() definition.JPLDefinition {
	return p.definition
}

func (p *program) OPs() map[definition.JPLOP]jpl.JPLOPHandler {
	return p.ops
}

func (p *program) Run(inputs []any, options *jpl.JPLProgramConfig) ([]any, jpl.JPLError) {
	if options == nil {
		options = new(jpl.JPLProgramConfig)
	}

	r := runtime.NewRuntime(p, &jpl.JPLRuntimeConfig{
		Runtime: jpl.ApplyRuntimeDefaults(options.Runtime, p.runtimeOptions),
	})

	normalizedInputs, err := library.NormalizeValues(inputs, "program inputs")
	if err != nil {
		return nil, err
	}

	outputs, err := r.Execute(normalizedInputs)
	if err != nil {
		return nil, err
	}

	stripped, err := library.StripJSON(outputs)
	if err != nil {
		return nil, err
	}
	return stripped.([]any), nil
}

func (p *program) Prepend(programDefinition definition.JPLDefinition) (jpl.JPLProgram, jpl.JPLError) {
	if err := validateDefinition(programDefinition); err != nil {
		return nil, err
	}

	originalInstructions := p.Definition().Instructions
	additionalInstructions := programDefinition.Instructions

	instructions := make(definition.Pipe, 0, len(originalInstructions)+len(additionalInstructions))
	instructions = append(instructions, additionalInstructions...)
	instructions = append(instructions, originalInstructions...)

	mergedDefinition := definition.JPLDefinition{
		Version:      definition.DEFINITION_VERSION,
		Instructions: instructions,
	}

	return NewProgram(mergedDefinition, &jpl.JPLProgramConfig{
		Program: p.options,
		Runtime: p.runtimeOptions,
	})
}

func (p *program) Append(programDefinition definition.JPLDefinition) (jpl.JPLProgram, jpl.JPLError) {
	if err := validateDefinition(programDefinition); err != nil {
		return nil, err
	}

	originalInstructions := p.Definition().Instructions
	additionalInstructions := programDefinition.Instructions

	instructions := make(definition.Pipe, 0, len(originalInstructions)+len(additionalInstructions))
	instructions = append(instructions, originalInstructions...)
	instructions = append(instructions, additionalInstructions...)

	mergedDefinition := definition.JPLDefinition{
		Version:      definition.DEFINITION_VERSION,
		Instructions: instructions,
	}

	return NewProgram(mergedDefinition, &jpl.JPLProgramConfig{
		Program: p.options,
		Runtime: p.runtimeOptions,
	})
}
