package brainfun

import (
	"testing"

	"github.com/LevOspennikov/brainfun/stack"
	"github.com/stretchr/testify/assert"
)

func getNewState() *State {
	data := make([]byte, DataSize)
	history := make([]rune, 0)
	return &State{Data: data, history: history, Stack: stack.NewStack()}
}

func getNewStateWithSize(size int) *State {
	data := make([]byte, size)
	history := make([]rune, 0)
	return &State{Data: data, history: history, Stack: stack.NewStack()}
}

func TestInc(t *testing.T) {
	state := getNewState()
	// add one
	err := Inc(state)
	assert.NoError(t, err)
	assert.Equal(t, 1, int(state.Data[0]))

	// add to the limit
	for i := 1; i < 255; i++ {
		err := Inc(state)
		assert.NoError(t, err)
	}
	assert.Equal(t, 255, int(state.Data[0]))

	// add to overflow
	err = Inc(state)
	assert.NoError(t, err)
	assert.Equal(t, 0, int(state.Data[0]))

	// test skip
	state.Skip = true
	err = Inc(state)
	assert.NoError(t, err)
	assert.Equal(t, 0, int(state.Data[0]))
}

func TestDec(t *testing.T) {
	state := getNewState()
	state.Data[0] = 255
	// sub one
	err := Dec(state)
	assert.NoError(t, err)
	assert.Equal(t, 254, int(state.Data[0]))

	// sub to the zero
	for i := 1; i < 255; i++ {
		err := Dec(state)
		assert.NoError(t, err)
	}
	assert.Equal(t, 0, int(state.Data[0]))

	// sub to overflow
	err = Dec(state)
	assert.NoError(t, err)
	assert.Equal(t, 255, int(state.Data[0]))

	// test skip
	state.Skip = true
	err = Dec(state)
	assert.NoError(t, err)
	assert.Equal(t, 255, int(state.Data[0]))
}

func TestMoveRight(t *testing.T) {
	state := getNewStateWithSize(4)

	state.DataPointer = 0
	// moveRight one
	err := MoveRight(state)
	assert.NoError(t, err)
	assert.Equal(t, 1, int(state.DataPointer))

	// moveRight to the edge
	for i := 2; i < 4; i++ {
		err := MoveRight(state)
		assert.NoError(t, err)
	}
	assert.Equal(t, 3, int(state.DataPointer))
	assert.NoError(t, Inc(state))

	// moveRight to overflow
	err = MoveRight(state)
	assert.NoError(t, err)
	assert.Equal(t, 0, int(state.DataPointer))

	// test skip
	state.Skip = true
	err = MoveRight(state)
	assert.NoError(t, err)
	assert.Equal(t, 0, int(state.DataPointer))
}

func TestMoveLeft(t *testing.T) {
	state := getNewStateWithSize(4)

	state.DataPointer = 3
	// MoveLeft one
	err := MoveLeft(state)
	assert.NoError(t, err)
	assert.Equal(t, 2, int(state.DataPointer))

	// MoveLeft to the edge
	for i := 2; i < 4; i++ {
		err := MoveLeft(state)
		assert.NoError(t, err)
	}
	assert.Equal(t, 0, int(state.DataPointer))
	assert.NoError(t, Inc(state))

	// MoveLeft to overflow
	err = MoveLeft(state)
	assert.NoError(t, err)
	assert.Equal(t, 3, int(state.DataPointer))

	// test skip
	state.Skip = true
	err = MoveLeft(state)
	assert.NoError(t, err)
	assert.Equal(t, 0, int(state.DataPointer))
}

func TestLoopStart(t *testing.T) {
	state := getNewState()

	// loop skip
	err := LoopStart(state)
	assert.NoError(t, err)
	assert.True(t, state.Skip)
	assert.Equal(t, 1, state.Stack.Size())

	state = getNewState()
	state.HistoryPointer = 1
	state.Data[0] = 1
	err = LoopStart(state)
	assert.NoError(t, err)
	assert.False(t, state.Skip)
	assert.Equal(t, state.HistoryPointer, 1)

}

func TestLoopEnd(t *testing.T) {
	state := getNewState()

	// loop skip
	err := LoopEnd(state)
	assert.Error(t, err)

	// loop back
	state = getNewState()
	state.HistoryPointer = 1
	state.Data[0] = 1
	state.Stack.Push(0)
	err = LoopEnd(state)
	assert.NoError(t, err)
	assert.False(t, state.Skip)
	assert.Equal(t, state.HistoryPointer, 0)

	// stop skipping
	state = getNewState()
	state.HistoryPointer = 1
	state.Data[0] = 1
	state.Skip = true
	state.Stack.Push(0)
	err = LoopEnd(state)
	assert.NoError(t, err)
	assert.False(t, state.Skip)
	assert.Equal(t, state.HistoryPointer, 1)
}
