package jpl

type JPLRuntimeSignal interface {
	// Return the signal's parent
	Parent() JPLRuntimeSignal

	// Check whether the current runtime area has been requested to be exited
	Exited() bool

	// Check that the current runtime area has not been requested to be exited and throw a fatal error otherwise
	CheckHealth() JPLFatalError

	// Request the current runtime area to be exited.
	// This also involves all child areas (introduced using `RuntimeSignal.next`).
	Exit()

	// Subscribe for when the current runtime area is requested to be exited.
	// This also involves all parent areas.
	// The function returns an unsubscription hook which must be called when completed in order to prevent memory leaks.
	// If the area has already been exited when subscribing, the callback is called immediately after the current event cycle.
	Subscribe(cb func()) func()

	// Inherit the next child area for the current one
	Next() JPLRuntimeSignal
}
