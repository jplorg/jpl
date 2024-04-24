package library

type RuntimeSignal interface {
	// Return the signal's parent
	Parent() RuntimeSignal

	// Check whether the current runtime area has been requested to be exited
	Exited() bool

	// Check that the current runtime area has not been requested to be exited and throw a fatal error otherwise
	CheckHealth() error

	// Request the current runtime area to be exited.
	// This also involves all child areas (introduced using `RuntimeSignal.next`).
	Exit()

	// Subscribe for when the current runtime area is requested to be exited.
	// This also involves all parent areas.
	// The function returns an unsubscription hook which must be called when completed in order to prevent memory leaks.
	// If the area has already been exited when subscribing, the callback is called immediately after the current event cycle.
	Subscribe(cb func()) func()

	// Inherit the next child area for the current one
	Next() RuntimeSignal
}

func NewRuntimeSignal(parent RuntimeSignal) RuntimeSignal {
	return &runtimeSignal{
		parent:        parent,
		subscriptions: make(map[int]func()),
	}
}

type runtimeSignal struct {
	parent              RuntimeSignal
	exited              bool
	subscriptions       map[int]func()
	nextSubscriptionKey int
}

func (s *runtimeSignal) Parent() RuntimeSignal {
	return s.parent
}

func (s *runtimeSignal) Exited() bool {
	if !s.exited && s.parent != nil {
		s.exited = s.parent.Exited()
	}
	return s.exited
}

func (s *runtimeSignal) CheckHealth() error {
	if s.exited {
		return NewJPLFatalError("execution aborted")
	}
	return nil
}

func (s *runtimeSignal) Exit() {
	if s.exited {
		return
	}
	s.exited = true
	for _, subscription := range s.subscriptions {
		subscription()
	}
	s.subscriptions = nil
}

func noop() {}

func (s *runtimeSignal) Subscribe(cb func()) func() {
	if s.exited {
		cb()
		return noop
	}
	key := s.nextSubscriptionKey
	s.nextSubscriptionKey += 1
	s.subscriptions[key] = cb
	unsubscribeParent := noop
	if s.parent != nil {
		unsubscribeParent = s.parent.Subscribe(cb)
	}
	unsubscribe := func() {
		delete(s.subscriptions, key)
		unsubscribeParent()
	}
	return unsubscribe
}

func (s *runtimeSignal) Next() RuntimeSignal {
	return NewRuntimeSignal(s)
}
