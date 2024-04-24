package program

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/2manyvcos/jpl/go/config"
	"github.com/2manyvcos/jpl/go/definition"
	"github.com/2manyvcos/jpl/go/jpl"
	"github.com/2manyvcos/jpl/go/library"
	"github.com/2manyvcos/jpl/go/runtime"
)

var defaultOptions = config.JPLProgramOptions{}

var versionRegex = regexp.MustCompile(`^(?P<Major>\d+)\.(?P<Minor>\d+)$`)

func validateDefinition(programDefinition definition.JPLDefinition) error {
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
		return fmt.Errorf("invalid program definition")
	}

	if majorV != definition.DEFINITION_VERSION_MAJOR || minorV > definition.DEFINITION_VERSION_MINOR {
		return fmt.Errorf("unsupported program definition v%s.%s - this version of JPL only supports v%v (up to v%s)", major, minor, definition.DEFINITION_VERSION_MAJOR, definition.DEFINITION_VERSION)
	}

	return nil
}

func NewProgram(programDefinition definition.JPLDefinition, options *config.JPLProgramConfig) (jpl.JPLProgram, error) {
	if err := validateDefinition(programDefinition); err != nil {
		return nil, err
	}

	if options == nil {
		options = new(config.JPLProgramConfig)
	}

	return &program{
		options:        config.ApplyProgramDefaults(options.Program, defaultOptions),
		runtimeOptions: options.Runtime,

		definition: programDefinition,
		ops:        OPs,
	}, nil
}

type program struct {
	options        config.JPLProgramOptions
	runtimeOptions config.JPLRuntimeOptions

	definition definition.JPLDefinition
	ops        map[string]jpl.JPLOPHandler
}

func (p *program) Options() config.JPLProgramOptions {
	return p.options
}

func (p *program) Definition() definition.JPLDefinition {
	return p.definition
}

func (p *program) OPs() map[string]jpl.JPLOPHandler {
	return p.ops
}

func (p *program) Run(inputs []any, options *config.JPLProgramConfig) ([]any, error) {
	if options == nil {
		options = new(config.JPLProgramConfig)
	}

	r := runtime.NewRuntime(p, &config.JPLRuntimeConfig{
		Runtime: config.ApplyRuntimeDefaults(options.Runtime, p.runtimeOptions),
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
