package library

import "github.com/jplorg/jpl/go/v2/jpl"

func NewRuntimeScope(presets *jpl.JPLRuntimeScopeConfig) jpl.JPLRuntimeScope {
	if presets == nil {
		presets = new(jpl.JPLRuntimeScopeConfig)
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
	signal jpl.JPLRuntimeSignal
	vars   map[string]any
}

func (s *runtimeScope) Signal() jpl.JPLRuntimeSignal {
	return s.signal
}

func (s *runtimeScope) Vars() map[string]any {
	return s.vars
}

func (s *runtimeScope) Next(modifications *jpl.JPLRuntimeScopeConfig) jpl.JPLRuntimeScope {
	if modifications == nil {
		modifications = new(jpl.JPLRuntimeScopeConfig)
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
