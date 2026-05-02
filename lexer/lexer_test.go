package lexer

import (
	"testing"
)

func TestLexerSimpleTokens(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		expect []TokenType
	}{
		{
			name:   "arithmetic operators",
			input:  "+ - * /",
			expect: []TokenType{PLUS, MINUS, MULTIPLY, DIVIDE, EOF},
		},
		{
			name:   "comparison operators",
			input:  "> < >= <= = <>",
			expect: []TokenType{GT, LT, GTE, LTE, EQ, NEQ, EOF},
		},
		{
			name:   "parentheses and comma",
			input:  "(a, b)",
			expect: []TokenType{LPAREN, IDENTIFIER, COMMA, IDENTIFIER, RPAREN, EOF},
		},
		{
			name:   "assignment",
			input:  "x := 10",
			expect: []TokenType{IDENTIFIER, ASSIGN, NUMBER, EOF},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := NewLexer(tt.input)
			tokens, err := lexer.Tokenize()
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if len(tokens) != len(tt.expect) {
				t.Fatalf("Expected %d tokens, got %d", len(tt.expect), len(tokens))
			}

			for i, expectedType := range tt.expect {
				if tokens[i].Type != expectedType {
					t.Errorf("Token %d: expected type %s, got %s", i, expectedType, tokens[i].Type)
				}
			}
		})
	}
}

func TestLexerNumbers(t *testing.T) {
	tests := []struct {
		name  string
		input string
		value string
	}{
		{"integer", "123", "123"},
		{"decimal", "123.456", "123.456"},
		{"scientific", "1.23e10", "1.23e10"},
		{"scientific with sign", "1.23e-5", "1.23e-5"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := NewLexer(tt.input)
			tokens, err := lexer.Tokenize()
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if len(tokens) != 2 { // NUMBER + EOF
				t.Fatalf("Expected 2 tokens, got %d", len(tokens))
			}

			if tokens[0].Type != NUMBER {
				t.Errorf("Expected NUMBER token, got %s", tokens[0].Type)
			}

			if tokens[0].Value != tt.value {
				t.Errorf("Expected value %s, got %s", tt.value, tokens[0].Value)
			}
		})
	}
}

func TestLexerIdentifiers(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		expectType TokenType
	}{
		{"simple identifier", "price", IDENTIFIER},
		{"uppercase identifier", "CLOSE", IDENTIFIER},
		{"identifier with numbers", "MA5", IDENTIFIER},
		{"identifier with underscore", "_temp", IDENTIFIER},
		{"chinese identifier", "阻力1", IDENTIFIER},
		{"keyword IF", "IF", IF},
		{"keyword AND", "AND", AND},
		{"keyword OR", "OR", OR},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := NewLexer(tt.input)
			tokens, err := lexer.Tokenize()
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if len(tokens) < 1 {
				t.Fatalf("Expected at least 1 token, got %d", len(tokens))
			}

			if tokens[0].Type != tt.expectType {
				t.Errorf("Expected type %s, got %s", tt.expectType, tokens[0].Type)
			}
		})
	}
}

func TestLexerFormulas(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		expect []TokenType
	}{
		{
			name:  "simple assignment",
			input: "MA5 := MA(CLOSE, 5)",
			expect: []TokenType{
				IDENTIFIER, ASSIGN, IDENTIFIER, LPAREN, IDENTIFIER, COMMA, NUMBER, RPAREN, EOF,
			},
		},
		{
			name:  "arithmetic expression",
			input: "x := (a + b) * c",
			expect: []TokenType{
				IDENTIFIER, ASSIGN, LPAREN, IDENTIFIER, PLUS, IDENTIFIER, RPAREN, MULTIPLY, IDENTIFIER, EOF,
			},
		},
		{
			name:  "comparison",
			input: "CLOSE > MA5",
			expect: []TokenType{
				IDENTIFIER, GT, IDENTIFIER, EOF,
			},
		},
		{
			name:  "tdx output style suffix",
			input: "DIF:EMA(CLOSE,12),COLORWHITE,LINETHICK2",
			expect: []TokenType{
				IDENTIFIER, COLON, IDENTIFIER, LPAREN, IDENTIFIER, COMMA, NUMBER, RPAREN, COMMA, COLOR, COMMA, LINETHICK, EOF,
			},
		},
		{
			name:  "quoted external reference",
			input: `"MACD.DIF#WEEK"`,
			expect: []TokenType{
				EXTERNAL_REFERENCE, EOF,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := NewLexer(tt.input)
			tokens, err := lexer.Tokenize()
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if len(tokens) != len(tt.expect) {
				t.Fatalf("Expected %d tokens, got %d", len(tt.expect), len(tokens))
			}

			for i, expectedType := range tt.expect {
				if tokens[i].Type != expectedType {
					t.Errorf("Token %d: expected type %s, got %s (value: %s)",
						i, expectedType, tokens[i].Type, tokens[i].Value)
				}
			}
		})
	}
}

func TestLexerNewlines(t *testing.T) {
	input := "a := 1\nb := 2"
	lexer := NewLexer(input)
	tokens, err := lexer.Tokenize()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expected := []TokenType{IDENTIFIER, ASSIGN, NUMBER, NEWLINE, IDENTIFIER, ASSIGN, NUMBER, EOF}
	if len(tokens) != len(expected) {
		t.Fatalf("Expected %d tokens, got %d", len(expected), len(tokens))
	}

	for i, expectedType := range expected {
		if tokens[i].Type != expectedType {
			t.Errorf("Token %d: expected type %s, got %s", i, expectedType, tokens[i].Type)
		}
	}

	// Check line numbers
	if tokens[0].Line != 1 {
		t.Errorf("Expected token 0 on line 1, got line %d", tokens[0].Line)
	}
	if tokens[4].Line != 2 {
		t.Errorf("Expected token 4 on line 2, got line %d", tokens[4].Line)
	}
}

func TestLexerErrors(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"unexpected character", "@"},
		{"unterminated string", `"MACD.DIF`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := NewLexer(tt.input)
			_, err := lexer.Tokenize()
			if err == nil {
				t.Error("Expected error, got nil")
			}
		})
	}
}

func TestLexerWhitespace(t *testing.T) {
	input := "  a   +   b  "
	lexer := NewLexer(input)
	tokens, err := lexer.Tokenize()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expected := []TokenType{IDENTIFIER, PLUS, IDENTIFIER, EOF}
	if len(tokens) != len(expected) {
		t.Fatalf("Expected %d tokens, got %d", len(expected), len(tokens))
	}

	for i, expectedType := range expected {
		if tokens[i].Type != expectedType {
			t.Errorf("Token %d: expected type %s, got %s", i, expectedType, tokens[i].Type)
		}
	}
}
