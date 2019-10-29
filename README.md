# Generic Go lexer

A very simple state-based lexer.

```go
import "src.userspace.com.au/lexer"

// Define the tokens for the lexer.
const (
	_ lexer.TokenType = iota
	tBREStart
	tBREEnd
	tRangeStart
	tRangeDash
	tRangeEnd
	tCharacter
	tClass
	tNot
)

// Define states returning a StateFunc.
func startState(l *lexer.Lexer) lexer.StateFunc {
	l.SkipWhitespace()
	r := l.Next()
	if r != '[' {
		return l.Error("expecting [")
	}
	l.Emit(tBREStart)
	return beFirstState
}
```

For a complete (but simple) example see the code in a [bracket expression
generator](https://src.userspace.com.au/bechars/tree/lexer.go).
