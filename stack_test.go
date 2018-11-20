package lexer

import (
	"testing"
)

func TestStack(t *testing.T) {
	s := newStack()
	s.push('r')
	r := s.pop()
	if r != 'r' {
		t.Fatalf("Expected r but got %b", r)
	}
	r = s.pop()
	if r != EOFRune {
		t.Fatalf("Expected EOFRune but got %b", r)
	}
}
