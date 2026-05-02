package lexer

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/DTrader-store/formula-go/errors"
)

// Lexer performs lexical analysis on formula source code
type Lexer struct {
	input  string   // the source code to analyze
	pos    int      // current position in input
	line   int      // current line number (1-indexed)
	column int      // current column number (1-indexed)
	tokens []*Token // collected tokens
}

// NewLexer creates a new Lexer instance
func NewLexer(input string) *Lexer {
	return &Lexer{
		input:  input,
		pos:    0,
		line:   1,
		column: 1,
		tokens: make([]*Token, 0),
	}
}

// Tokenize performs lexical analysis and returns all tokens
func (l *Lexer) Tokenize() ([]*Token, error) {
	for !l.isAtEnd() {
		if err := l.scanToken(); err != nil {
			return nil, err
		}
	}

	// Add EOF token
	l.addToken(EOF, "")
	return l.tokens, nil
}

// scanToken scans and adds a single token
func (l *Lexer) scanToken() error {
	// Skip whitespace (except newlines)
	for !l.isAtEnd() && l.isSpace(l.peek()) && l.peek() != '\n' {
		l.advance()
	}

	if l.isAtEnd() {
		return nil
	}

	ch := l.peek()

	// Handle newlines
	if ch == '\n' {
		l.addToken(NEWLINE, "\n")
		l.advance()
		l.line++
		l.column = 1
		return nil
	}

	// Handle numbers
	if unicode.IsDigit(ch) {
		return l.scanNumber()
	}

	// Handle quoted indicator references and text literals
	if ch == '\'' || ch == '"' {
		return l.scanString()
	}

	// Handle identifiers and keywords
	if unicode.IsLetter(ch) || ch == '_' {
		return l.scanIdentifier()
	}

	// Handle operators and punctuation
	return l.scanOperator()
}

// scanNumber scans a number token
func (l *Lexer) scanNumber() error {
	start := l.pos
	startCol := l.column

	// Scan integer part
	for !l.isAtEnd() && unicode.IsDigit(l.peek()) {
		l.advance()
	}

	// Scan decimal part
	if !l.isAtEnd() && l.peek() == '.' {
		l.advance() // consume '.'
		for !l.isAtEnd() && unicode.IsDigit(l.peek()) {
			l.advance()
		}
	}

	// Scan scientific notation
	if !l.isAtEnd() && (l.peek() == 'e' || l.peek() == 'E') {
		l.advance() // consume 'e' or 'E'
		if !l.isAtEnd() && (l.peek() == '+' || l.peek() == '-') {
			l.advance() // consume sign
		}
		for !l.isAtEnd() && unicode.IsDigit(l.peek()) {
			l.advance()
		}
	}

	value := l.input[start:l.pos]
	l.tokens = append(l.tokens, &Token{
		Type:   NUMBER,
		Value:  value,
		Line:   l.line,
		Column: startCol,
	})
	return nil
}

// scanString scans a quoted string token without the surrounding quotes.
func (l *Lexer) scanString() error {
	quote := l.advance()
	start := l.pos
	startCol := l.column - 1

	for !l.isAtEnd() && l.peek() != quote {
		if l.peek() == '\n' {
			return errors.NewLexerError("unterminated string", l.line, startCol, string(quote))
		}
		l.advance()
	}

	if l.isAtEnd() {
		return errors.NewLexerError("unterminated string", l.line, startCol, string(quote))
	}

	value := l.input[start:l.pos]
	l.advance()
	tokenType := STRING
	if quote == '"' {
		tokenType = EXTERNAL_REFERENCE
	}
	l.tokens = append(l.tokens, &Token{
		Type:   tokenType,
		Value:  value,
		Line:   l.line,
		Column: startCol,
	})
	return nil
}

// scanIdentifier scans an identifier or keyword
func (l *Lexer) scanIdentifier() error {
	start := l.pos
	startCol := l.column

	for !l.isAtEnd() && (unicode.IsLetter(l.peek()) || unicode.IsDigit(l.peek()) || l.peek() == '_') {
		l.advance()
	}

	value := l.input[start:l.pos]
	upper := strings.ToUpper(value)

	// Check for keywords
	tokenType := l.getKeywordType(upper)
	l.tokens = append(l.tokens, &Token{
		Type:   tokenType,
		Value:  value,
		Line:   l.line,
		Column: startCol,
	})
	return nil
}

// getKeywordType returns the token type for a keyword or IDENTIFIER
func (l *Lexer) getKeywordType(upper string) TokenType {
	keywords := map[string]TokenType{
		"IF":      IF,
		"AND":     AND,
		"OR":      OR,
		"DOTLINE": DOTLINE,
		"STICK":   STICK,
	}

	if tokenType, exists := keywords[upper]; exists {
		return tokenType
	}
	if strings.HasPrefix(upper, "COLOR") {
		return COLOR
	}
	if strings.HasPrefix(upper, "LINETHICK") {
		return LINETHICK
	}
	return IDENTIFIER
}

// scanOperator scans operators and punctuation
func (l *Lexer) scanOperator() error {
	startCol := l.column
	ch := l.advance()

	switch ch {
	case '+':
		l.addToken(PLUS, "+")
	case '-':
		l.addToken(MINUS, "-")
	case '*':
		l.addToken(MULTIPLY, "*")
	case '/':
		l.addToken(DIVIDE, "/")
	case '(':
		l.addToken(LPAREN, "(")
	case ')':
		l.addToken(RPAREN, ")")
	case ',':
		l.addToken(COMMA, ",")
	case ';':
		l.addToken(SEMICOLON, ";")
	case ':':
		// Check for := (assignment)
		if !l.isAtEnd() && l.peek() == '=' {
			l.advance()
			l.addToken(ASSIGN, ":=")
		} else {
			l.addToken(COLON, ":")
		}
	case '>':
		// Check for >=
		if !l.isAtEnd() && l.peek() == '=' {
			l.advance()
			l.addToken(GTE, ">=")
		} else {
			l.addToken(GT, ">")
		}
	case '<':
		// Check for <= or <>
		if !l.isAtEnd() && l.peek() == '=' {
			l.advance()
			l.addToken(LTE, "<=")
		} else if !l.isAtEnd() && l.peek() == '>' {
			l.advance()
			l.addToken(NEQ, "<>")
		} else {
			l.addToken(LT, "<")
		}
	case '=':
		// Check for ==
		if !l.isAtEnd() && l.peek() == '=' {
			l.advance()
			l.addToken(EQ, "==")
		} else {
			l.addToken(EQ, "=")
		}
	case '!':
		// Check for !=
		if !l.isAtEnd() && l.peek() == '=' {
			l.advance()
			l.addToken(NEQ, "!=")
		} else {
			return errors.NewLexerError("unexpected character", l.line, startCol, string(ch))
		}
	case '#':
		l.addToken(HASH, "#")
	default:
		return errors.NewLexerError(fmt.Sprintf("unexpected character: %c", ch), l.line, startCol, string(ch))
	}

	return nil
}

// peek returns the current character without advancing
func (l *Lexer) peek() rune {
	if l.isAtEnd() {
		return 0
	}
	ch, _ := utf8.DecodeRuneInString(l.input[l.pos:])
	return ch
}

// advance returns the current character and moves to the next
func (l *Lexer) advance() rune {
	if l.isAtEnd() {
		return 0
	}
	ch, size := utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += size
	l.column++
	return ch
}

// isAtEnd checks if we've reached the end of input
func (l *Lexer) isAtEnd() bool {
	return l.pos >= len(l.input)
}

// isSpace checks if a character is whitespace (excluding newline)
func (l *Lexer) isSpace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\r'
}

// addToken adds a token to the tokens list
func (l *Lexer) addToken(tokenType TokenType, value string) {
	l.tokens = append(l.tokens, &Token{
		Type:   tokenType,
		Value:  value,
		Line:   l.line,
		Column: l.column - len(value),
	})
}
