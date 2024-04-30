package jpl

type IterFunc = func(key *string, value any) (result any, remove bool, err JPLError)

// JPL stripper
type JPLStripper interface {
	Strip(key *string, value any, iter IterFunc) (result any, remove bool, err JPLError)
}

type JPLStripperFunc func(key *string, value any, iter IterFunc) (result any, remove bool, err JPLError)

// JPLStripperFunc implements JPLStripper
var _ JPLStripper = JPLStripperFunc(nil)

func (p JPLStripperFunc) Strip(key *string, value any, iter IterFunc) (result any, remove bool, err JPLError) {
	return p(key, value, iter)
}
