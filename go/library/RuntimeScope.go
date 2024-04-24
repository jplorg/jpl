package library

type RuntimeScopeConfig struct {
	Signal RuntimeSignal
	Vars   map[string]any
}

type RuntimeScope interface {
	Signal() RuntimeSignal

	Vars() map[string]any

	// Inherit the next scope based on the specified modifications
	Next(modifications *RuntimeScopeConfig) RuntimeScope
}

func NewRuntimeScope(presets *RuntimeScopeConfig) RuntimeScope {
	if presets == nil {
		presets = new(RuntimeScopeConfig)
	}

	signal := presets.Signal
	if signal == nil {
		signal = NewRuntimeSignal(nil)
	}

	return &runtimeScope{
		signal: signal,
		vars:   presets.Vars,
	}
}

type runtimeScope struct {
	signal RuntimeSignal
	vars   map[string]any
}

func (s *runtimeScope) Signal() RuntimeSignal {
	return s.signal
}

func (s *runtimeScope) Vars() map[string]any {
	return s.vars
}

func (s *runtimeScope) Next(modifications *RuntimeScopeConfig) RuntimeScope {
	if modifications == nil {
		modifications = new(RuntimeScopeConfig)
	}

	signal := modifications.Signal
	if signal == nil {
		signal = s.signal
	}

	return &runtimeScope{
		signal: signal,
		vars:   ApplyObject(s.vars, ObjectEntries(modifications.Vars)),
	}
}
