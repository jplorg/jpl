package library

type IterFunc = func(key *string, value any) (any, error)

// JPL stripper
type JPLStripper interface {
	Strip(key *string, value any, iter IterFunc) (any, error)
}

type JPLStripperFunc func(key *string, value any, iter IterFunc) (any, error)

// JPLStripperFunc implements JPLStripper
var _ JPLStripper = JPLStripperFunc(nil)

func (p JPLStripperFunc) Strip(key *string, value any, iter IterFunc) (any, error) {
	return p(key, value, iter)
}
