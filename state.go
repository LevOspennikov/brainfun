package brainfun

import (
	"errors"

	"github.com/LevOspennikov/brainfun/stack"
)

type State struct {
	Stack          stack.Stack
	Data           []byte
	Skip           bool
	DataPointer    uint64
	HistoryPointer int
	history        []rune
}

func (state *State) next() (rune, error) {
	if state.HistoryPointer >= len(state.history) {
		return 0, errors.New("no next rune")
	}
	state.HistoryPointer++
	return state.history[state.HistoryPointer-1], nil
}

func (state *State) put(char rune) {
	state.history = append(state.history, char)
}

func (state State) hasNext() bool {
	return state.HistoryPointer < len(state.history)
}

func (state State) isTerminal() bool {
	stackEmpty := state.Stack == nil || state.Stack.IsEmpty()
	historyEnded := !state.hasNext()
	return stackEmpty && historyEnded
}
