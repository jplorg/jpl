package library

type JPLRuntimeScopeConfig struct {
	Signal JPLRuntimeSignal
	Vars   map[string]any
}

type JPLRuntimeScope interface {
	Signal() JPLRuntimeSignal
	Vars() map[string]any

	// Inherit the next scope based on the specified modifications
	Next(modifications *JPLRuntimeScopeConfig) JPLRuntimeScope
}

func NewRuntimeScope(presets *JPLRuntimeScopeConfig) JPLRuntimeScope {
	if presets == nil {
		presets = new(JPLRuntimeScopeConfig)
	}

	signal := presets.Signal
	if signal == nil {
		signal = NewRuntimeSignal()
	}

	return &runtimeScope{
		signal: signal,
		vars:   presets.Vars,
	}
}

type runtimeScope struct {
	signal JPLRuntimeSignal
	vars   map[string]any
}

func (s *runtimeScope) Signal() JPLRuntimeSignal {
	return s.signal
}

func (s *runtimeScope) Vars() map[string]any {
	return s.vars
}

func (s *runtimeScope) Next(modifications *JPLRuntimeScopeConfig) JPLRuntimeScope {
	if modifications == nil {
		modifications = new(JPLRuntimeScopeConfig)
	}

	signal := modifications.Signal
	if signal == nil {
		signal = s.signal
	}

	return &runtimeScope{
		signal: signal,
		vars:   applyObject(s.vars, objectEntries(modifications.Vars)),
	}
}
