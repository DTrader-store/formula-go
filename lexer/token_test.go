package lexer

import "testing"

func TestNewToken(t *testing.T) {
	token := NewToken(NUMBER, "123", 1, 5)
	if token.Type != NUMBER {
		t.Errorf("Expected type NUMBER, got %s", token.Type)
	}
	if token.Value != "123" {
		t.Errorf("Expected value '123', got '%s'", token.Value)
	}
	if token.Line != 1 {
		t.Errorf("Expected line 1, got %d", token.Line)
	}
	if token.Column != 5 {
		t.Errorf("Expected column 5, got %d", token.Column)
	}
}

func TestTokenString(t *testing.T) {
	token := NewToken(IDENTIFIER, "test", 2, 10)
	expected := `Token(IDENTIFIER, "test", 2:10)`
	if token.String() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, token.String())
	}
}

func TestTokenEquals(t *testing.T) {
	tests := []struct {
		name   string
		token1 *Token
		token2 *Token
		equal  bool
	}{
		{
			name:   "equal tokens",
			token1: NewToken(NUMBER, "123", 1, 5),
			token2: NewToken(NUMBER, "123", 1, 5),
			equal:  true,
		},
		{
			name:   "different types",
			token1: NewToken(NUMBER, "123", 1, 5),
			token2: NewToken(IDENTIFIER, "123", 1, 5),
			equal:  false,
		},
		{
			name:   "different values",
			token1: NewToken(NUMBER, "123", 1, 5),
			token2: NewToken(NUMBER, "456", 1, 5),
			equal:  false,
		},
		{
			name:   "different lines",
			token1: NewToken(NUMBER, "123", 1, 5),
			token2: NewToken(NUMBER, "123", 2, 5),
			equal:  false,
		},
		{
			name:   "different columns",
			token1: NewToken(NUMBER, "123", 1, 5),
			token2: NewToken(NUMBER, "123", 1, 10),
			equal:  false,
		},
		{
			name:   "nil token",
			token1: NewToken(NUMBER, "123", 1, 5),
			token2: nil,
			equal:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.token1.Equals(tt.token2)
			if result != tt.equal {
				t.Errorf("Expected %v, got %v", tt.equal, result)
			}
		})
	}
}

func TestTokenTypes(t *testing.T) {
	// Test that all token type constants are defined
	tokenTypes := []TokenType{
		NUMBER, STRING, EXTERNAL_REFERENCE, IDENTIFIER,
		PLUS, MINUS, MULTIPLY, DIVIDE,
		GT, LT, GTE, LTE, EQ, NEQ,
		AND, OR,
		LPAREN, RPAREN, COMMA, SEMICOLON, COLON,
		ASSIGN, IF,
		COLOR, LINETHICK, DOTLINE, STICK,
		HASH, NEWLINE, EOF,
	}

	for _, tt := range tokenTypes {
		if string(tt) == "" {
			t.Errorf("Token type should not be empty")
		}
	}
}
