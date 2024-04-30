package library

import "github.com/jplorg/jpl/go/jpl"

func NewRuntimeSignal(parent jpl.JPLRuntimeSignal) jpl.JPLRuntimeSignal {
	return &runtimeSignal{
		parent:        parent,
		subscriptions: make(map[int]func()),
	}
}

type runtimeSignal struct {
	parent              jpl.JPLRuntimeSignal
	exited              bool
	subscriptions       map[int]func()
	nextSubscriptionKey int
}

func (s *runtimeSignal) Parent() jpl.JPLRuntimeSignal {
	return s.parent
}

func (s *runtimeSignal) Exited() bool {
	if !s.exited && s.parent != nil {
		s.exited = s.parent.Exited()
	}
	return s.exited
}

func (s *runtimeSignal) CheckHealth() jpl.JPLFatalError {
	if s.exited {
		return NewFatalError("execution aborted")
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

func (s *runtimeSignal) Next() jpl.JPLRuntimeSignal {
	return NewRuntimeSignal(s)
}
