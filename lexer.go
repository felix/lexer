package lexer

import (
	"bytes"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

// StateFunc captures the movement from one state to the next.
type StateFunc func(*Lexer) StateFunc

// TokenType identifies the tokens emitted.
type TokenType int

const (
	// EOFRune is a convenience for EOF
	EOFRune rune = -1
	// ErrorToken is returned on error
	ErrorToken TokenType = -1
	// EOFToken is return on EOF
	EOFToken TokenType = 0
)

var lineSep = []byte{'\n'}

// Token is returned by the lexer.
type Token struct {
	Type     TokenType
	Value    string
	Position int
	Line     int
}

// String implements Stringer
func (t Token) String() string {
	return fmt.Sprintf("[%d] %s", t.Type, t.Value)
}

// Lexer represents the lexer machine.
type Lexer struct {
	source     string
	start      int
	line       int
	position   int
	lastWidth  int
	startState StateFunc
	tokens     chan Token
	history    stack
}

// New creates a returns a lexer ready to parse the given source code.
func New(src string, start StateFunc) *Lexer {
	return &Lexer{
		source:     src,
		startState: start,
		start:      0,
		line:       1,
		position:   0,
		history:    newStack(),
	}
}

// Start begins executing the Lexer in an asynchronous manner (using a goroutine).
func (l *Lexer) Start() {
	go l.StartSync()
}

// StartSync starts the lexer synchronously.
func (l *Lexer) StartSync() {
	// Take half the string length as a buffer size.
	buffSize := len(l.source) / 2
	if buffSize <= 0 {
		buffSize = 1
	}
	l.tokens = make(chan Token, buffSize)
	l.run()
}

func (l *Lexer) run() {
	state := l.startState
	for state != nil {
		state = state(l)
	}
	close(l.tokens)
}

// Current returns the value being being analyzed at this moment.
func (l *Lexer) Current() string {
	return l.source[l.start:l.position]
}

// Emit will receive a token type and push a new token with the current analyzed
// value into the tokens channel.
func (l *Lexer) Emit(t TokenType) {
	tok := Token{
		Type:     t,
		Value:    l.Current(),
		Position: l.position,
		Line:     l.line,
	}
	l.tokens <- tok
	l.checkLines()
	l.start = l.position
	l.history.clear()
}

func (l *Lexer) checkLines() {
	val := l.Current()
	l.line += bytes.Count([]byte(val), lineSep)
}

// Next pulls the next rune from the Lexer and returns it, moving the position
// forward in the source.
func (l *Lexer) Next() rune {
	var r rune
	var s int
	str := l.source[l.position:]
	if len(str) == 0 {
		r, s = EOFRune, 0
	} else {
		r, s = utf8.DecodeRuneInString(str)
	}
	l.position += s
	l.history.push(r)

	return r
}

// Ignore clears the history stack and then sets the current beginning position
// to the current position in the source which effectively ignores the section
// of the source being analyzed.
func (l *Lexer) Ignore() {
	l.history.clear()
	l.checkLines()
	l.start = l.position
}

// Peek performs a Next operation immediately followed by a Backup returning the
// peeked rune.
func (l *Lexer) Peek() rune {
	r := l.Next()
	l.Backup()

	return r
}

// Backup will take the last rune read (if any) and history back. Backups can
// occur more than once per call to Next but you can never history past the
// last point a token was emitted.
func (l *Lexer) Backup() {
	r := l.history.pop()
	if r > EOFRune {
		size := utf8.RuneLen(r)
		l.position -= size
		if l.position < l.start {
			l.position = l.start
		}
	}
}

// Accept receives a string containing all acceptable strings and will continue
// over each consecutive character in the source until a token not in the given
// string is encountered. This should be used to quickly pull token parts.
func (l *Lexer) Accept(valid string) bool {
	if strings.IndexRune(valid, l.Next()) >= 0 {
		return true
	}
	l.Backup() // last next wasn't a match
	return false
}

// AcceptRun consumes a run of runes from the valid set.
func (l *Lexer) AcceptRun(valid string) (n int) {
	for strings.IndexRune(valid, l.Next()) >= 0 {
		n++
	}
	l.Backup() // last next wasn't a match
	return n
}

// SkipWhitespace continues over all unicode whitespace.
func (l *Lexer) SkipWhitespace() {
	for {
		r := l.Next()

		if !unicode.IsSpace(r) {
			l.Backup()
			break
		}

		if r == EOFRune {
			l.Emit(EOFToken)
			break
		}
	}
}

// Tokens returns the a token channel.
func (l *Lexer) Tokens() <-chan Token {
	return l.tokens
}

// NextToken returns the next token from the lexer and done
func (l *Lexer) NextToken() (*Token, bool) {
	if tok, ok := <-l.tokens; ok {
		return &tok, false
	}
	return nil, true
}

func (l *Lexer) Error(format string, args ...interface{}) StateFunc {
	l.tokens <- Token{
		Type:     ErrorToken,
		Value:    fmt.Sprintf(format, args...),
		Position: l.position,
		Line:     l.line,
	}
	return nil
}
