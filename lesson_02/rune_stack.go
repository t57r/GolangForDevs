package main

type RuneStack struct {
	data []rune
}

func NewRuneStack() *RuneStack {
	return &RuneStack{data: make([]rune, 0)}
}

func (s *RuneStack) IsEmpty() bool {
	return len(s.data) == 0
}

func (s *RuneStack) Push(r rune) {
	s.data = append(s.data, r)
}

func (s *RuneStack) Pop() (rune, bool) {
	n := len(s.data)
	if n == 0 {
		return 0, false
	}
	top := s.data[n-1]
	s.data = s.data[:n-1]
	return top, true
}
