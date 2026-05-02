// Package lexer provides the lexical analysis (tokenization) functionality
package lexer

// TokenType represents the type of a token
type TokenType string

// Token type constants
const (
	// Literals
	NUMBER             TokenType = "NUMBER"
	STRING             TokenType = "STRING"
	EXTERNAL_REFERENCE TokenType = "EXTERNAL_REFERENCE"
	IDENTIFIER         TokenType = "IDENTIFIER"

	// Operators
	PLUS     TokenType = "PLUS"
	MINUS    TokenType = "MINUS"
	MULTIPLY TokenType = "MULTIPLY"
	DIVIDE   TokenType = "DIVIDE"

	// Comparison operators
	GT  TokenType = "GT"
	LT  TokenType = "LT"
	GTE TokenType = "GTE"
	LTE TokenType = "LTE"
	EQ  TokenType = "EQ"
	NEQ TokenType = "NEQ"

	// Logical operators
	AND TokenType = "AND"
	OR  TokenType = "OR"

	// Punctuation
	LPAREN    TokenType = "LPAREN"
	RPAREN    TokenType = "RPAREN"
	COMMA     TokenType = "COMMA"
	SEMICOLON TokenType = "SEMICOLON"
	COLON     TokenType = "COLON"

	// Assignment
	ASSIGN TokenType = "ASSIGN"

	// Keywords
	IF TokenType = "IF"

	// Chart attributes
	COLOR     TokenType = "COLOR"
	LINETHICK TokenType = "LINETHICK"
	DOTLINE   TokenType = "DOTLINE"
	STICK     TokenType = "STICK"

	// Special
	HASH    TokenType = "HASH"
	NEWLINE TokenType = "NEWLINE"
	EOF     TokenType = "EOF"
)
