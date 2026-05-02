package parser

import (
	"testing"

	"github.com/DTrader-store/formula-go/lexer"
	"github.com/DTrader-store/formula-go/parser/ast"
)

func TestParserSimpleExpressions(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"number", "42", "NumberLiteral"},
		{"identifier", "CLOSE", "Identifier"},
		{"addition", "1 + 2", "BinaryExpression"},
		{"subtraction", "5 - 3", "BinaryExpression"},
		{"multiplication", "4 * 5", "BinaryExpression"},
		{"division", "10 / 2", "BinaryExpression"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.NewLexer(tt.input)
			tokens, err := l.Tokenize()
			if err != nil {
				t.Fatalf("Lexer error: %v", err)
			}

			p := NewParser(tokens)
			program, err := p.Parse()
			if err != nil {
				t.Fatalf("Parser error: %v", err)
			}

			if len(program.Body) == 0 {
				t.Fatal("Expected at least one statement")
			}
		})
	}
}

func TestParserVariableDeclaration(t *testing.T) {
	input := "MA5 := MA(CLOSE, 5)"

	l := lexer.NewLexer(input)
	tokens, err := l.Tokenize()
	if err != nil {
		t.Fatalf("Lexer error: %v", err)
	}

	p := NewParser(tokens)
	program, err := p.Parse()
	if err != nil {
		t.Fatalf("Parser error: %v", err)
	}

	if len(program.Body) != 1 {
		t.Fatalf("Expected 1 statement, got %d", len(program.Body))
	}

	varDecl, ok := program.Body[0].(*ast.VariableDeclaration)
	if !ok {
		t.Fatalf("Expected VariableDeclaration, got %T", program.Body[0])
	}

	if varDecl.Name != "MA5" {
		t.Errorf("Expected name 'MA5', got '%s'", varDecl.Name)
	}

	funcCall, ok := varDecl.Value.(*ast.FunctionCall)
	if !ok {
		t.Fatalf("Expected FunctionCall, got %T", varDecl.Value)
	}

	if funcCall.Name != "MA" {
		t.Errorf("Expected function 'MA', got '%s'", funcCall.Name)
	}

	if len(funcCall.Arguments) != 2 {
		t.Errorf("Expected 2 arguments, got %d", len(funcCall.Arguments))
	}
}

func TestParserVariableDeclarationIgnoresTDXStyleSuffix(t *testing.T) {
	input := "阻力2 := MA(REF(HHV(H,15*Z1),1),M1), COLORCYAN;"

	l := lexer.NewLexer(input)
	tokens, err := l.Tokenize()
	if err != nil {
		t.Fatalf("Lexer error: %v", err)
	}

	p := NewParser(tokens)
	program, err := p.Parse()
	if err != nil {
		t.Fatalf("Parser error: %v", err)
	}

	if len(program.Body) != 1 {
		t.Fatalf("Expected 1 statement, got %d", len(program.Body))
	}
	varDecl, ok := program.Body[0].(*ast.VariableDeclaration)
	if !ok {
		t.Fatalf("Expected VariableDeclaration, got %T", program.Body[0])
	}
	if varDecl.Name != "阻力2" {
		t.Errorf("Expected Chinese identifier name, got %s", varDecl.Name)
	}
}

func TestParserOutputDeclarationWithStyle(t *testing.T) {
	input := "DIF: EMA(CLOSE, 12), COLORWHITE, LINETHICK2, DOTLINE"

	l := lexer.NewLexer(input)
	tokens, err := l.Tokenize()
	if err != nil {
		t.Fatalf("Lexer error: %v", err)
	}

	p := NewParser(tokens)
	program, err := p.Parse()
	if err != nil {
		t.Fatalf("Parser error: %v", err)
	}

	if len(program.Body) != 1 {
		t.Fatalf("Expected 1 statement, got %d", len(program.Body))
	}

	outputDecl, ok := program.Body[0].(*ast.OutputDeclaration)
	if !ok {
		t.Fatalf("Expected OutputDeclaration, got %T", program.Body[0])
	}
	if outputDecl.Name != "DIF" {
		t.Errorf("Expected name 'DIF', got '%s'", outputDecl.Name)
	}
	if outputDecl.Style == nil {
		t.Fatal("Expected style metadata")
	}
	if outputDecl.Style.Color == nil || *outputDecl.Style.Color != "WHITE" {
		t.Errorf("Expected color WHITE, got %#v", outputDecl.Style.Color)
	}
	if outputDecl.Style.LineWidth == nil || *outputDecl.Style.LineWidth != 2 {
		t.Errorf("Expected line width 2, got %#v", outputDecl.Style.LineWidth)
	}
	if outputDecl.Style.LineStyle == nil || *outputDecl.Style.LineStyle != "dotted" {
		t.Errorf("Expected dotted line style, got %#v", outputDecl.Style.LineStyle)
	}
}

func TestParserPrecedence(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"multiply before add", "2 + 3 * 4"},
		{"parentheses", "(2 + 3) * 4"},
		{"comparison", "a > b"},
		{"logical and", "a > 5 AND b < 10"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.NewLexer(tt.input)
			tokens, err := l.Tokenize()
			if err != nil {
				t.Fatalf("Lexer error: %v", err)
			}

			p := NewParser(tokens)
			_, err = p.Parse()
			if err != nil {
				t.Fatalf("Parser error: %v", err)
			}
		})
	}
}

func TestParserFunctionCall(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		funcName string
		argCount int
	}{
		{"no args", "COUNT()", "COUNT", 0},
		{"one arg", "SUM(prices)", "SUM", 1},
		{"two args", "MA(CLOSE, 5)", "MA", 2},
		{"three args", "IF(a > b, a, b)", "IF", 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.NewLexer(tt.input)
			tokens, err := l.Tokenize()
			if err != nil {
				t.Fatalf("Lexer error: %v", err)
			}

			p := NewParser(tokens)
			program, err := p.Parse()
			if err != nil {
				t.Fatalf("Parser error: %v", err)
			}

			if len(program.Body) == 0 {
				t.Fatal("Expected at least one statement")
			}

			exprStmt, ok := program.Body[0].(*ast.ExpressionStatement)
			if !ok {
				t.Fatalf("Expected ExpressionStatement, got %T", program.Body[0])
			}

			funcCall, ok := exprStmt.Expr.(*ast.FunctionCall)
			if !ok {
				t.Fatalf("Expected FunctionCall, got %T", exprStmt.Expr)
			}

			if funcCall.Name != tt.funcName {
				t.Errorf("Expected function '%s', got '%s'", tt.funcName, funcCall.Name)
			}

			if len(funcCall.Arguments) != tt.argCount {
				t.Errorf("Expected %d arguments, got %d", tt.argCount, len(funcCall.Arguments))
			}
		})
	}
}

func TestParserStringLiteral(t *testing.T) {
	input := `"MACD.DIF#WEEK"`

	l := lexer.NewLexer(input)
	tokens, err := l.Tokenize()
	if err != nil {
		t.Fatalf("Lexer error: %v", err)
	}

	p := NewParser(tokens)
	program, err := p.Parse()
	if err != nil {
		t.Fatalf("Parser error: %v", err)
	}

	exprStmt, ok := program.Body[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Expected ExpressionStatement, got %T", program.Body[0])
	}
	stringLit, ok := exprStmt.Expr.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("Expected StringLiteral, got %T", exprStmt.Expr)
	}
	if stringLit.Value != "MACD.DIF#WEEK" {
		t.Errorf("Expected external ref value, got %s", stringLit.Value)
	}
	if !stringLit.External {
		t.Error("Expected string literal to be marked as external")
	}
}

func TestParserUnaryExpression(t *testing.T) {
	input := "-5"

	l := lexer.NewLexer(input)
	tokens, err := l.Tokenize()
	if err != nil {
		t.Fatalf("Lexer error: %v", err)
	}

	p := NewParser(tokens)
	program, err := p.Parse()
	if err != nil {
		t.Fatalf("Parser error: %v", err)
	}

	if len(program.Body) == 0 {
		t.Fatal("Expected at least one statement")
	}

	exprStmt, ok := program.Body[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Expected ExpressionStatement, got %T", program.Body[0])
	}

	unaryExpr, ok := exprStmt.Expr.(*ast.UnaryExpression)
	if !ok {
		t.Fatalf("Expected UnaryExpression, got %T", exprStmt.Expr)
	}

	if unaryExpr.Operator != ast.OpUnaryMinus {
		t.Errorf("Expected unary minus operator")
	}
}

func TestParserMultipleStatements(t *testing.T) {
	input := `
		MA5 := MA(CLOSE, 5)
		MA10 := MA(CLOSE, 10)
		CROSS(MA5, MA10)
	`

	l := lexer.NewLexer(input)
	tokens, err := l.Tokenize()
	if err != nil {
		t.Fatalf("Lexer error: %v", err)
	}

	p := NewParser(tokens)
	program, err := p.Parse()
	if err != nil {
		t.Fatalf("Parser error: %v", err)
	}

	if len(program.Body) != 3 {
		t.Fatalf("Expected 3 statements, got %d", len(program.Body))
	}
}
