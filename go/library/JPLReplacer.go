package library

// JPL replacer
type JPLReplacer interface {
	Replace(key string, value any) (any, JPLError)
}

type JPLReplacerFunc func(key string, value any) (any, JPLError)

// JPLReplacerFunc implements JPLReplacer
var _ JPLReplacer = JPLReplacerFunc(nil)

func (p JPLReplacerFunc) Replace(key string, value any) (any, JPLError) {
	return p(key, value)
}
