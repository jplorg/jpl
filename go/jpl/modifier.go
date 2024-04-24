package jpl

// JPL modifier
type JPLModifier interface {
	Modify(input any) (output any, err JPLError)
}

type JPLModifierFunc func(input any) (output any, err JPLError)

// JPLModifierFunc implements JPLModifier
var _ JPLModifier = JPLModifierFunc(nil)

func (p JPLModifierFunc) Modify(input any) (any, JPLError) {
	return p(input)
}
