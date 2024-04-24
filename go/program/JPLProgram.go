package program

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/2manyvcos/jpl/go/config"
	"github.com/2manyvcos/jpl/go/library/definition"
)

var defaultOptions = config.JPLProgramOptions{}

type Options struct {
	Program config.JPLProgramOptions
	Runtime config.JPLRuntimeOptions
}

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

// JPL program
type JPLProgram interface {
	// Run the program with the provided inputs and runtime options.
	// The program throws a JPLExecutionError for runtime failures.
	// Other errors may be thrown when execution fails.
	Run(inputs []any, options *Options) ([]any, error)

	// Return the program's definition.
	// The definition can be serialized as JSON to be reused in other JPL implementations.
	Definition() definition.JPLDefinition
}

func NewProgram(programDefinition definition.JPLDefinition, options *Options) (JPLProgram, error) {
	if err := validateDefinition(programDefinition); err != nil {
		return nil, err
	}

	if options == nil {
		options = new(Options)
	}

	return &program{
		Options: config.ApplyProgramDefaults(options.Program, defaultOptions),

		definition: programDefinition,
		OPs:        OPs,

		RuntimeOptions: options.Runtime,
	}, nil
}

type program struct {
	Options config.JPLProgramOptions

	definition definition.JPLDefinition
	OPs        map[string]JPLOP

	RuntimeOptions config.JPLRuntimeOptions
}

func (p *program) Run(inputs []any, options *Options) ([]any, error) {
	panic("TODO:")
}

func (p *program) Definition() definition.JPLDefinition {
	return p.definition
}
