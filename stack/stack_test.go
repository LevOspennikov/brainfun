package stack

import (
	"testing"
)

func Test_Stack(t *testing.T) {
	s := ArrayStack{}

	if s.Size() != 0 {
		t.Errorf("Length of an empty stack should be 0")
	}

	s.Push(1)

	if s.Size() != 1 {
		t.Errorf("Length should be 0")
	}

	if val, _ := s.Pop(); val != 1 {
		t.Errorf("Top item should have been 1")
	}

	if s.Size() != 0 {
		t.Errorf("Stack should be empty")
	}

	s.Push(1)
	s.Push(2)

	if s.Size() != 2 {
		t.Errorf("Length should be 2")
	}

	if val, _ := s.Pop(); val != 2 {
		t.Errorf("Top item should have been 2")
	}

	if val, _ := s.Pop(); val != 1 {
		t.Errorf("Top item should have been 1")
	}

	_, ok := s.Pop()
	if ok {
		t.Errorf("should be not ok")
	}
}
