package builtins

import (
	"math"

	"github.com/jplorg/jpl/go/v2/jpl"
	"github.com/jplorg/jpl/go/v2/library"
)

func unwrapNumber(v any) (float64, jpl.JPLError) {
	t, err := library.Type(v)
	if err != nil {
		return 0, err
	}
	u, err := library.UnwrapValue(v)
	if err != nil {
		return 0, err
	}
	if t != jpl.JPLT_NUMBER {
		return 0, library.ThrowAny(library.NewTypeError("%s (%*<100v) cannot be used for mathematical operations", string(t), u))
	}
	return u.(float64), nil
}

type alterFunc = func(runtime jpl.JPLRuntime, value float64, args ...any) (any, jpl.JPLError)

func funcMath(alter alterFunc) jpl.JPLFunc {
	return enclose(func(runtime jpl.JPLRuntime, signal jpl.JPLRuntimeSignal, next jpl.JPLPiper, input any, args ...any) ([]any, error) {
		t, err := library.Type(input)
		if err != nil {
			return nil, err
		}
		if t != jpl.JPLT_NUMBER {
			u, err := library.UnwrapValue(input)
			if err != nil {
				return nil, err
			}
			return nil, library.ThrowAny(library.NewTypeError("%s (%*<100v) cannot be used for mathematical operations", string(t), u))
		}
		alteredValue, err := library.AlterValue(input, jpl.JPLModifierFunc(func(value any) (any, jpl.JPLError) {
			return alter(runtime, value.(float64), args...)
		}))
		if err != nil {
			return nil, err
		}
		return next.Pipe(alteredValue)
	})
}

var funcsMath = map[string]any{
	"pow": funcMath(func(runtime jpl.JPLRuntime, value float64, args ...any) (any, jpl.JPLError) {
		var arg0 any
		if len(args) > 0 {
			arg0 = args[0]
		}
		p, err := unwrapNumber(arg0)
		if err != nil {
			return nil, err
		}
		return math.Pow(value, p), nil
	}),
	"sqrt": funcMath(func(runtime jpl.JPLRuntime, value float64, args ...any) (any, jpl.JPLError) {
		return math.Sqrt(value), nil
	}),
	"exp": funcMath(func(runtime jpl.JPLRuntime, value float64, args ...any) (any, jpl.JPLError) {
		return math.Exp(value), nil
	}),
	"log": funcMath(func(runtime jpl.JPLRuntime, value float64, args ...any) (any, jpl.JPLError) {
		return math.Log(value), nil
	}),
	"log10": funcMath(func(runtime jpl.JPLRuntime, value float64, args ...any) (any, jpl.JPLError) {
		return math.Log10(value), nil
	}),

	"sin": funcMath(func(runtime jpl.JPLRuntime, value float64, args ...any) (any, jpl.JPLError) {
		return math.Sin(value), nil
	}),
	"cos": funcMath(func(runtime jpl.JPLRuntime, value float64, args ...any) (any, jpl.JPLError) {
		return math.Cos(value), nil
	}),
	"tan": funcMath(func(runtime jpl.JPLRuntime, value float64, args ...any) (any, jpl.JPLError) {
		return math.Tan(value), nil
	}),
	"asin": funcMath(func(runtime jpl.JPLRuntime, value float64, args ...any) (any, jpl.JPLError) {
		return math.Asin(value), nil
	}),
	"acos": funcMath(func(runtime jpl.JPLRuntime, value float64, args ...any) (any, jpl.JPLError) {
		return math.Acos(value), nil
	}),
	"atan": funcMath(func(runtime jpl.JPLRuntime, value float64, args ...any) (any, jpl.JPLError) {
		return math.Atan(value), nil
	}),

	"ceil": funcMath(func(runtime jpl.JPLRuntime, value float64, args ...any) (any, jpl.JPLError) {
		return math.Ceil(value), nil
	}),
	"floor": funcMath(func(runtime jpl.JPLRuntime, value float64, args ...any) (any, jpl.JPLError) {
		return math.Floor(value), nil
	}),
	"round": funcMath(func(runtime jpl.JPLRuntime, value float64, args ...any) (any, jpl.JPLError) {
		return math.RoundToEven(value), nil
	}),
	"trunc": funcMath(func(runtime jpl.JPLRuntime, value float64, args ...any) (any, jpl.JPLError) {
		return math.Trunc(value), nil
	}),

	"abs": funcMath(func(runtime jpl.JPLRuntime, value float64, args ...any) (any, jpl.JPLError) {
		return math.Abs(value), nil
	}),
}
