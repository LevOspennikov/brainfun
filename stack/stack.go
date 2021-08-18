package stack

type ArrayStack []int

func NewStack() *ArrayStack {
	stack := ArrayStack([]int{})
	return &stack
}

func (s ArrayStack) IsEmpty() bool {
	return len(s) == 0
}

func (s ArrayStack) Size() int {
	return len(s)
}

func (s *ArrayStack) Push(str int) {
	*s = append(*s, str)
}

func (s *ArrayStack) Pop() (int, bool) {
	if s.IsEmpty() {
		return 0, false
	} else {
		index := len(*s) - 1
		element := (*s)[index]
		*s = (*s)[:index]
		return element, true
	}
}

type Stack interface {
	IsEmpty() bool
	Size() int
	Push(str int)
	Pop() (int, bool)
}
