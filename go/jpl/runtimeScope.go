package jpl

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
