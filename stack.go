package lexer

type stackNode struct {
	r    rune
	next *stackNode
}

type stack struct {
	start *stackNode
}

func newStack() stack {
	return stack{}
}

func (s *stack) push(r rune) {
	node := &stackNode{r: r}
	if s.start == nil {
		s.start = node
	} else {
		node.next = s.start
		s.start = node
	}
}

func (s *stack) pop() rune {
	if s.start == nil {
		return EOFRune
	}

	n := s.start
	s.start = n.next
	return n.r
}

func (s *stack) clear() {
	s.start = nil
}
