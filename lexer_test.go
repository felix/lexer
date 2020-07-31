package lexer

import (
	"fmt"
	"testing"
)

const (
	NumberToken TokenType = iota
	OpToken
	IdentToken
)

func NumberState(l *Lexer) StateFunc {
	l.AcceptRun("0123456789")
	l.Emit(NumberToken)
	if l.Peek() == '.' {
		l.Next()
		l.Emit(OpToken)
		return IdentState
	}

	return nil
}

func IdentState(l *Lexer) StateFunc {
	r := l.Next()
	for (r >= 'a' && r <= 'z') || r == '_' {
		r = l.Next()
	}
	l.Backup()
	l.Emit(IdentToken)

	return WhitespaceState
}

func NewlineState(l *Lexer) StateFunc {
	l.AcceptRun("0123456789")
	l.Emit(NumberToken)
	l.SkipWhitespace()
	l.Ignore()
	l.AcceptRun("0123456789")
	l.Emit(NumberToken)
	l.SkipWhitespace()

	return nil
}

func WhitespaceState(l *Lexer) StateFunc {
	r := l.Next()
	if r == EOFRune {
		return nil
	}

	if r != ' ' && r != '\t' && r != '\n' && r != '\r' {
		l.Error(fmt.Sprintf("unexpected token %q", r))
		return nil
	}

	l.Accept(" \t\n\r")
	l.Ignore()

	return NumberState
}

func TestMovingThroughString(t *testing.T) {
	l := New("123", nil)
	run := []struct {
		s string
		r rune
	}{
		{"1", '1'},
		{"12", '2'},
		{"123", '3'},
		{"123", EOFRune},
	}

	for _, test := range run {
		r := l.Next()
		if r != test.r {
			t.Fatalf("Expected %q but got %q", test.r, r)
		}

		if l.Current() != test.s {
			t.Fatalf("Expected %q but got %q", test.s, l.Current())
		}
	}
}

func TestNumbers(t *testing.T) {
	l := New("123", NumberState)
	l.StartSync()
	tok, done := l.NextToken()
	if done {
		t.Fatal("Expected a token but lexer finished")
	}

	if tok.Type != NumberToken {
		t.Fatalf("Expected a %v but got %v", NumberToken, tok.Type)
	}

	if tok.Value != "123" {
		t.Fatalf("Expected %q but got %q", "123", tok.Value)
	}

	s := "[0] 123"
	if tok.String() != s {
		t.Fatalf("Expected %q but got %q", s, tok.String())
	}

	tok, done = l.NextToken()
	if !done {
		t.Fatal("Expected done but it wasn't.")
	}

	if tok != nil {
		t.Fatalf("Expected a nil token but got %v", *tok)
	}
}

func TestNewlines(t *testing.T) {
	src := `123
456
789`
	l := New(src, NewlineState)
	l.StartSync()
	tok, done := l.NextToken()
	if done {
		t.Fatal("Expected done, but it wasn't.")
	}

	if tok.Type != NumberToken {
		t.Fatalf("Expected a number token but got %v", *tok)
	}

	if tok.Value != "123" {
		t.Fatalf("Expected 123 but got %q", tok.Value)
	}

	if tok.Line != 1 {
		t.Fatalf("Expected line 1 but got %d", tok.Line)
	}

	tok, done = l.NextToken()
	if done {
		t.Fatal("Expected not done, but it was.")
	}

	if tok.Type != NumberToken {
		t.Fatalf("Expected a number token but got %v", *tok)
	}

	if tok.Value != "456" {
		t.Fatalf("Expected 456 but got %q", tok.Value)
	}

	if tok.Line != 2 {
		t.Fatalf("Expected line 2 but got %d", tok.Line)
	}
}

func TestBackup(t *testing.T) {
	l := New("1", nil)
	r := l.Next()
	if r != '1' {
		t.Fatalf("Expected %q but got %q", '1', r)
	}

	if l.Current() != "1" {
		t.Fatalf("Expected %q but got %q", "1", l.Current())
	}

	l.Backup()
	if l.Current() != "" {
		t.Fatalf("Expected empty string but got %q", l.Current())
	}
}

func TestWhitespace(t *testing.T) {
	l := New("    1", NumberState)
	l.StartSync()
	l.SkipWhitespace()

	tok, done := l.NextToken()
	if done {
		t.Fatal("Expected token to be !done, but it was.")
	}

	if tok.Type != NumberToken {
		t.Fatalf("Expected number token but got %v", *tok)
	}
}

func TestMultipleTokens(t *testing.T) {
	cases := []struct {
		tokType TokenType
		val     string
	}{
		{NumberToken, "123"},
		{OpToken, "."},
		{IdentToken, "hello"},
		{NumberToken, "675"},
		{OpToken, "."},
		{IdentToken, "world"},
	}

	l := New("123.hello  675.world", NumberState)
	l.StartSync()

	for _, c := range cases {
		tok, done := l.NextToken()
		if done {
			t.Fatal("Expected there to be more tokens but there weren't")
		}

		if c.tokType != tok.Type {
			t.Fatalf("Expected token type %v but got %v", c.tokType, tok.Type)
		}

		if c.val != tok.Value {
			t.Fatalf("Expected %q but got %q", c.val, tok.Value)
		}
	}

	tok, done := l.NextToken()
	if !done {
		t.Fatal("Expected the lexer to be done but it wasn't.")
	}

	if tok != nil {
		t.Fatalf("Did not expect a token, but got %v", *tok)
	}
}

func TestError(t *testing.T) {
	l := New("notaspace", WhitespaceState)
	l.StartSync()

	tok, done := l.NextToken()
	if done {
		t.Fatal("Expected token to be !done but it was.")
	}

	if tok.Type != ErrorToken {
		t.Fatalf("Expected error token but got %v", *tok)
	}
}
