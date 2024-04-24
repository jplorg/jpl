package library

// JPL modifier
type JPLModifier interface {
	Modify(input any) (output any, err error)
}

type JPLModifierFunc func(input any) (output any, err error)

// JPLModifierFunc implements JPLModifier
var _ JPLModifier = JPLModifierFunc(nil)

func (p JPLModifierFunc) Modify(input any) (any, error) {
	return p(input)
}
