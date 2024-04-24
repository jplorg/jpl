package options

type JPLRuntimeOptions struct {
	Vars map[string]any
}

func ApplyRuntimeDefaults(options JPLRuntimeOptions, defaults JPLRuntimeOptions) (result JPLRuntimeOptions) {
	result.Vars = make(map[string]any)
	for k, v := range defaults.Vars {
		result.Vars[k] = v
	}
	for k, v := range options.Vars {
		result.Vars[k] = v
	}

	return
}
