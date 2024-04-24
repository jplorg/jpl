package library

// JPL replacer
type JPLReplacer interface {
	Replace(key string, value any) (any, error)
}

type JPLReplacerFunc func(key string, value any) (any, error)

// JPLReplacerFunc implements JPLReplacer
var _ JPLReplacer = JPLReplacerFunc(nil)

func (p JPLReplacerFunc) Replace(key string, value any) (any, error) {
	return p(key, value)
}
