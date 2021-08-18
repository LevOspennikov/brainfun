package brainfun

import (
	"errors"
	"io"
	"math"
)

func Inc(state *State) error {
	if state.Skip {
		return nil
	}
	if state.DataPointer >= uint64(len(state.Data)) {
		return errors.New("array out of bound")
	}
	val := state.Data[state.DataPointer]
	if val == math.MaxUint8 {
		val = 0
	} else {
		val++
	}
	state.Data[state.DataPointer] = val
	return nil
}

func Dec(state *State) error {
	if state.Skip {
		return nil
	}
	if state.DataPointer >= uint64(len(state.Data)) {
		return errors.New("array out of bound")
	}
	val := state.Data[state.DataPointer]
	if val == 0 {
		val = math.MaxUint8
	} else {
		val--
	}
	state.Data[state.DataPointer] = val
	return nil
}

func MoveLeft(state *State) error {
	if state.Skip {
		return nil
	}
	val := state.DataPointer
	if val == 0 {
		val = uint64(len(state.Data))
	}
	val--
	state.DataPointer = val
	return nil
}

func MoveRight(state *State) error {
	if state.Skip {
		return nil
	}
	val := state.DataPointer
	if val >= uint64(len(state.Data))-1 {
		val = 0
	} else {
		val++
	}
	state.DataPointer = val
	return nil
}

func LoopStart(state *State) error {
	if state.Skip {
		val, ok := state.Stack.Pop() // bracket counter, not the link to the start of the loop
		if !ok {
			return errors.New("wrong start of the loop: no beginning")
		}
		state.Stack.Push(val + 1)
		return nil
	}

	if state.DataPointer >= uint64(len(state.Data)) {
		return errors.New("array out of bound")
	}
	val := state.Data[state.DataPointer]
	if val == 0 {
		state.Skip = true
		state.Stack.Push(0) // bracket counter, not the link to the start
		return nil
	}
	state.Stack.Push(state.HistoryPointer - 1)

	return nil
}

func LoopEnd(state *State) error {
	if state.Stack.Size() == 0 {
		return errors.New("wrong end of the loop: no beginning")
	}
	stackVal, _ := state.Stack.Pop()
	if state.Skip {
		if stackVal == 0 {
			// bracket counter equals zero, can stop skipping
			state.Skip = false
			return nil
		}
		state.Stack.Push(stackVal - 1)
		return nil
	}

	val := state.Data[state.DataPointer]
	if val > 0 {
		state.HistoryPointer = stackVal
		return nil
	}
	return nil
}

func PrintWith(writer io.Writer) Transform {
	return func(state *State) error {
		if state.Skip {
			return nil
		}
		if state.DataPointer >= uint64(len(state.Data)) {
			return errors.New("array out of bound")
		}
		val := state.Data[state.DataPointer]
		_, err := writer.Write([]byte(string(val)))
		return err
	}
}

func ReadWith(reader io.Reader) Transform {
	return func(state *State) error {
		if state.Skip {
			return nil
		}
		if state.DataPointer >= uint64(len(state.Data)) {
			return errors.New("array out of bound")
		}
		buf := make([]byte, 1)
		_, err := reader.Read(buf)
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		state.Data[state.DataPointer] = buf[0]
		return nil
	}
}
