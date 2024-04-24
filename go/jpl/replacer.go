package jpl

// JPL replacer
type JPLReplacer interface {
	Replace(key string, value any) (result any, remove bool, err JPLError)
}

type JPLReplacerFunc func(key string, value any) (result any, remove bool, err JPLError)

// JPLReplacerFunc implements JPLReplacer
var _ JPLReplacer = JPLReplacerFunc(nil)

func (p JPLReplacerFunc) Replace(key string, value any) (result any, remove bool, err JPLError) {
	return p(key, value)
}
