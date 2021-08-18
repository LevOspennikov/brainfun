package brainfun

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/LevOspennikov/brainfun/stack"
)

const DataSize int = 32767

type Transform func(*State) error

type Interpreter struct {
	state      *State
	operations map[rune]Transform
}

// NewInterpreter returns new instance with the default set of functions
func NewInterpreter() Interpreter {
	data := make([]byte, DataSize)
	history := make([]rune, 0)
	return Interpreter{
		state: &State{Data: data, history: history, Stack: stack.NewStack()},
		operations: map[rune]Transform{
			'+': Inc,
			'-': Dec,
			'<': MoveLeft,
			'>': MoveRight,
			'[': LoopStart,
			']': LoopEnd,
			'.': PrintWith(os.Stdout),
			',': ReadWith(os.Stdin),
		},
	}
}

// Execute executes functions linked with runes from reader. Read is happened one-by-one, not using look-ahead.
func (int *Interpreter) Execute(reader io.RuneReader) error {
	char, _, err := reader.ReadRune()
	for err == nil {
		// put new char to history
		int.state.put(char)
		// while history is not ended, execute
		for int.state.hasNext() {
			charFromState, err := int.state.next()
			if err != nil {
				return err
			}
			err = int.executeChar(charFromState)
			if err != nil {
				return err
			}
		}
		char, _, err = reader.ReadRune()
	}
	if !int.state.isTerminal() {
		return errors.New("execution error: state is not terminal")
	}
	return nil
}

func (int *Interpreter) executeChar(char rune) error {
	transformer, ok := int.operations[char]
	if !ok {
		return fmt.Errorf("no function for rune %v", char)
	}
	err := transformer(int.state)
	return err
}

// AddFunc creates link between rune and function. Function will be executed later, when rune is occurred
func (int *Interpreter) AddFunc(char rune, function Transform) {
	int.operations[char] = function
}

// DeleteFunc deletes link between rune and function. Function will not be executed later, when rune is occurred
func (int *Interpreter) DeleteFunc(char rune) {
	delete(int.operations, char)
}
