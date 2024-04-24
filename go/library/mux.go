package library

// JPL muxer
type IMuxer[Input any] interface {
	Mux(args ...Input) JPLError
}

type IMuxerFunc[Input any] func(args ...Input) JPLError

// IMuxerFunc implements IMuxer
var _ IMuxer[any] = IMuxerFunc[any](nil)

func (p IMuxerFunc[Input]) Mux(args ...Input) JPLError {
	return p(args...)
}

// JPL muxer
type IOMuxer[Input any, Output any] interface {
	Mux(args ...Input) (Output, JPLError)
}

type IOMuxerFunc[Input any, Output any] func(args ...Input) (Output, JPLError)

// IOMuxerFunc implements IOMuxer
var _ IOMuxer[any, any] = IOMuxerFunc[any, any](nil)

func (p IOMuxerFunc[Input, Output]) Mux(args ...Input) (Output, JPLError) {
	return p(args...)
}

// Multiplex the specified array of arguments by calling cb for all possible combinations of arguments.
//
// `mux([[1,2], [3,4]], cb)` for example yields:
// - `cb(1, 3)`
// - `cb(1, 4)`
// - `cb(2, 3)`
// - `cb(2, 4)`
func Mux[Input any](args [][]Input, cb IMuxer[Input]) JPLError {
	argCount := len(args)
	if argCount == 1 {
		for _, arg := range args[0] {
			if err := cb.Mux(arg); err != nil {
				return err
			}
		}
		return nil
	}
	execCount := 1
	indices := make([]*iteratorIndex[Input], argCount)
	for i, arg := range args {
		argLength := len(arg)
		execCount *= argLength
		indices[i] = &iteratorIndex[Input]{argLength, 0, arg}
	}
	if execCount == 0 {
		return nil
	}
	buffer := make([]Input, argCount)
	for i, arg := range args {
		buffer[i] = arg[0]
	}
	execIndex := 0
	for {
		if err := cb.Mux(buffer...); err != nil {
			return err
		}
		execIndex += 1
		if execIndex >= execCount {
			break
		}
		// determine next combination
		for argIndex := argCount - 1; argIndex >= 0; argIndex -= 1 {
			arg := indices[argIndex]
			next := arg.current + 1
			if next < arg.max {
				arg.current = next
				buffer[argIndex] = arg.values[next]
				break
			}
			arg.current = 0
			buffer[argIndex] = arg.values[0]
		}
	}
	return nil
}

// Multiplex the specified array of arguments and return the results produced by the callbacks
func MuxOne[Input any, Output any](args [][]Input, cb IOMuxer[Input, Output]) ([]Output, JPLError) {
	argCount := len(args)
	if argCount == 1 {
		inputs := args[0]
		result := make([]Output, len(inputs))
		for i, arg := range inputs {
			var err JPLError
			if result[i], err = cb.Mux(arg); err != nil {
				return nil, err
			}
		}
		return result, nil
	}
	execCount := 1
	indices := make([]*iteratorIndex[Input], argCount)
	for i, arg := range args {
		argLength := len(arg)
		execCount *= argLength
		indices[i] = &iteratorIndex[Input]{argLength, 0, arg}
	}
	if execCount == 0 {
		return nil, nil
	}
	outputs := make([]Output, execCount)
	buffer := make([]Input, argCount)
	for i, arg := range args {
		buffer[i] = arg[0]
	}
	execIndex := 0
	for {
		var err JPLError
		if outputs[execIndex], err = cb.Mux(buffer...); err != nil {
			return nil, err
		}
		execIndex += 1
		if execIndex >= execCount {
			break
		}
		// determine next combination
		for argIndex := argCount - 1; argIndex >= 0; argIndex -= 1 {
			arg := indices[argIndex]
			next := arg.current + 1
			if next < arg.max {
				arg.current = next
				buffer[argIndex] = arg.values[next]
				break
			}
			arg.current = 0
			buffer[argIndex] = arg.values[0]
		}
	}
	return outputs, nil
}

// Multiplex the specified array of arguments and return a single array of all merged result arrays produced by the callbacks
func MuxAll[Input any, Output any](args [][]Input, cb IOMuxer[Input, []Output]) ([]Output, JPLError) {
	segments, err := MuxOne(args, cb)
	if err != nil {
		return nil, err
	}
	return MergeSegments(segments), nil
}

// Create a single array from the specified array segments
func MergeSegments[Value any](segments [][]Value) []Value {
	if len(segments) == 0 {
		return nil
	}
	if len(segments) == 1 {
		return segments[0]
	}
	var count int
	for _, segment := range segments {
		count += len(segment)
	}
	result := make([]Value, 0, count)
	for _, segment := range segments {
		result = append(result, segment...)
	}
	return result
}
