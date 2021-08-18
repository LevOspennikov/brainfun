package brainfun

import (
	"bytes"
	"errors"
	"io"
	"strconv"
	"strings"
	"testing"
)

func TestInterpreter_Execute(t *testing.T) {
	var tests = []struct {
		name       string
		str        string
		wantErr    bool
		wantString string
	}{
		{
			name:       "hello world",
			str:        "++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++.+++++++++++++++++++++++++++++.+++++++..+++.-------------------------------------------------------------------------------.+++++++++++++++++++++++++++++++++++++++++++++++++++++++.++++++++++++++++++++++++.+++.------.--------.-------------------------------------------------------------------.-----------------------.",
			wantErr:    false,
			wantString: "Hello World!\n",
		},
		{
			name:       "hello world updated",
			str:        "++++++++++[>+++++++>++++++++++>+++>+<<<<-]>++.>+.+++++++..+++.>++.<<+++++++++++++++.>.+++.------.--------.>+.>.",
			wantErr:    false,
			wantString: "Hello World!\n",
		},
		{
			name:       "fibonacci",
			str:        "+++++++++++>+>>>>++++++++++++++++++++++++++++++++++++++++++++>++++++++++++++++++++++++++++++++<<<<<<[>[>>>>>>+>+<<<<<<<-]>>>>>>>[<<<<<<<+>>>>>>>-]<[>++++++++++[-<-[>>+>+<<<-]>>>[<<<+>>>-]+<[>[-]<[-]]>[<<[>>>+<<<-]>>[-]]<<]>>>[>>+>+<<<-]>>>[<<<+>>>-]+<[>[-]<[-]]>[<<+>>[-]]<<<<<<<]>>>>>[++++++++++++++++++++++++++++++++++++++++++++++++.[-]]++++++++++<[->-<]>++++++++++++++++++++++++++++++++++++++++++++++++.[-]<<<<<<<<<<<<[>>>+>+<<<<-]>>>>[<<<<+>>>>-]<-[>>.>.<<<[-]]<<[>>+>+<<<-]>>>[<<<+>>>-]<<[<+>-]>[<+>-]<<<-]",
			wantErr:    false,
			wantString: "1, 1, 2, 3, 5, 8, 13, 21, 34, 55, 89",
		},
		{
			name:    "correct brackets",
			str:     "[]",
			wantErr: false,
		},
		{
			name:    "wrong brackets",
			str:     "[]]",
			wantErr: true,
		},
		{
			name:    "brackets in brackets",
			str:     "+[[[-]]]",
			wantErr: false,
		},
		{
			name:    "not closing brackets",
			str:     "[[]][-",
			wantErr: true,
		},
		{
			name:       "bubblesort",
			str:        ">>,[>>,]<<[[<<]>>>>[<<[>+<<+>-]>>[>+<<<<[->]>[<]>>-]<<<[[-]>>[>+<-]>>[<<<+>>>-]]>>[[<+>-]>>]<]<<[>>+<<-]<<]>>>>[.>>]",
			wantErr:    false,
			wantString: "ABCD",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			interpreter := NewInterpreter()
			buf := &bytes.Buffer{}
			interpreter.AddFunc('.', PrintWith(buf))
			interpreter.AddFunc(',', ReadWith(strings.NewReader("DBAC")))
			if err := interpreter.Execute(strings.NewReader(tt.str)); (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
			} else if buf.String() != tt.wantString {
				t.Errorf("Execute() result = %s, want result %s", buf.String(), tt.wantString)
			}
		})
	}
}

func TestInterpreter_AddFunc(t *testing.T) {
	squareTransform := func(state *State) error {
		if state.Skip {
			return nil
		}
		if state.DataPointer >= uint64(len(state.Data)) {
			return errors.New("array out of bound")
		}
		val := state.Data[state.DataPointer]
		result := uint16(val) * uint16(val)

		state.Data[state.DataPointer] = uint8(result % 256)
		return nil
	}

	printNumTransform := func(writer io.Writer) Transform {
		return func(state *State) error {
			if state.Skip {
				return nil
			}
			if state.DataPointer >= uint64(len(state.Data)) {
				return errors.New("array out of bound")
			}
			val := state.Data[state.DataPointer]
			byteVal := strconv.Itoa(int(val))
			_, err := writer.Write([]byte(byteVal))
			return err
		}
	}
	var tests = []struct {
		name         string
		str          string
		newTransform Transform
		wantErr      bool
		wantString   string
	}{
		{
			name: "square",
			str:  "+'+'^'",
			newTransform: func(state *State) error {
				if state.Skip {
					return nil
				}
				if state.DataPointer >= uint64(len(state.Data)) {
					return errors.New("array out of bound")
				}
				val := state.Data[state.DataPointer]
				result := uint16(val) * uint16(val)

				state.Data[state.DataPointer] = uint8(result % 256)
				return nil

			},
			wantErr:    false,
			wantString: "124",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			interpreter := NewInterpreter()
			buf := &bytes.Buffer{}
			interpreter.AddFunc('.', PrintWith(buf))
			interpreter.AddFunc('\'', printNumTransform(buf))
			interpreter.AddFunc('^', squareTransform)
			if err := interpreter.Execute(strings.NewReader(tt.str)); (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
			} else if buf.String() != tt.wantString {
				t.Errorf("Execute() result = %s, want result %s", buf.String(), tt.wantString)
			}
		})
	}
}
