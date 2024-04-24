package jpl

type IterFunc = func(key *string, value any) (any, JPLError)

// JPL stripper
type JPLStripper interface {
	Strip(key *string, value any, iter IterFunc) (any, JPLError)
}

type JPLStripperFunc func(key *string, value any, iter IterFunc) (any, JPLError)

// JPLStripperFunc implements JPLStripper
var _ JPLStripper = JPLStripperFunc(nil)

func (p JPLStripperFunc) Strip(key *string, value any, iter IterFunc) (any, JPLError) {
	return p(key, value, iter)
}
