package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/DTrader-store/formula-go/errors"
	"github.com/DTrader-store/formula-go/lexer"
	"github.com/DTrader-store/formula-go/parser/ast"
)

// Parser performs syntactic analysis on tokens
type Parser struct {
	tokens  []*lexer.Token
	pos     int
	current *lexer.Token
}

// NewParser creates a new Parser instance
func NewParser(tokens []*lexer.Token) *Parser {
	p := &Parser{
		tokens: tokens,
		pos:    0,
	}
	if len(tokens) > 0 {
		p.current = tokens[0]
	}
	return p
}

// Parse parses the tokens and returns an AST
func (p *Parser) Parse() (*ast.Program, error) {
	statements := make([]ast.Statement, 0)

	for !p.isAtEnd() {
		// Skip newlines
		if p.current.Type == lexer.NEWLINE {
			p.advance()
			continue
		}

		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}

		if stmt != nil {
			statements = append(statements, stmt)
		}
	}

	return &ast.Program{Body: statements}, nil
}

// parseStatement parses a single statement
func (p *Parser) parseStatement() (ast.Statement, error) {
	// Check for assignment (identifier := expression)
	if p.current.Type == lexer.IDENTIFIER && p.peek() != nil && p.peek().Type == lexer.ASSIGN {
		return p.parseVariableDeclaration()
	}

	// TDX output declaration: name: expression[, style...]
	if p.current.Type == lexer.IDENTIFIER && p.peek() != nil && p.peek().Type == lexer.COLON {
		return p.parseOutputDeclaration()
	}

	// Otherwise, parse as expression statement
	expr, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	// Skip optional semicolon or newline
	if !p.isAtEnd() && (p.current.Type == lexer.SEMICOLON || p.current.Type == lexer.NEWLINE) {
		p.advance()
	}

	return &ast.ExpressionStatement{Expr: expr}, nil
}

// parseOutputDeclaration parses a TDX output declaration: name: expression[, style...]
func (p *Parser) parseOutputDeclaration() (*ast.OutputDeclaration, error) {
	name := p.current.Value
	p.advance() // consume identifier

	if p.current.Type != lexer.COLON {
		return nil, p.error("expected :")
	}
	p.advance() // consume :

	value, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	style, err := p.parseStyleSuffixes()
	if err != nil {
		return nil, err
	}

	// Skip optional semicolon or newline
	if !p.isAtEnd() && (p.current.Type == lexer.SEMICOLON || p.current.Type == lexer.NEWLINE) {
		p.advance()
	}

	return &ast.OutputDeclaration{Name: name, Value: value, Style: style}, nil
}

// parseVariableDeclaration parses a variable declaration: name := expression
func (p *Parser) parseVariableDeclaration() (*ast.VariableDeclaration, error) {
	name := p.current.Value
	p.advance() // consume identifier

	if p.current.Type != lexer.ASSIGN {
		return nil, p.error("expected :=")
	}
	p.advance() // consume :=

	value, err := p.parseExpression()
	if err != nil {
		return nil, err
	}
	if _, err := p.parseStyleSuffixes(); err != nil {
		return nil, err
	}

	// Skip optional semicolon or newline
	if !p.isAtEnd() && (p.current.Type == lexer.SEMICOLON || p.current.Type == lexer.NEWLINE) {
		p.advance()
	}

	return &ast.VariableDeclaration{Name: name, Value: value}, nil
}

// parseExpression parses an expression (handles operator precedence)
func (p *Parser) parseExpression() (ast.Expression, error) {
	return p.parseLogicalOr()
}

// parseLogicalOr parses logical OR expressions
func (p *Parser) parseLogicalOr() (ast.Expression, error) {
	left, err := p.parseLogicalAnd()
	if err != nil {
		return nil, err
	}

	for !p.isAtEnd() && p.current.Type == lexer.OR {
		p.advance()
		right, err := p.parseLogicalAnd()
		if err != nil {
			return nil, err
		}
		left = &ast.BinaryExpression{
			Left:     left,
			Operator: ast.OpOr,
			Right:    right,
		}
	}

	return left, nil
}

// parseLogicalAnd parses logical AND expressions
func (p *Parser) parseLogicalAnd() (ast.Expression, error) {
	left, err := p.parseComparison()
	if err != nil {
		return nil, err
	}

	for !p.isAtEnd() && p.current.Type == lexer.AND {
		p.advance()
		right, err := p.parseComparison()
		if err != nil {
			return nil, err
		}
		left = &ast.BinaryExpression{
			Left:     left,
			Operator: ast.OpAnd,
			Right:    right,
		}
	}

	return left, nil
}

// parseComparison parses comparison expressions
func (p *Parser) parseComparison() (ast.Expression, error) {
	left, err := p.parseAdditive()
	if err != nil {
		return nil, err
	}

	for !p.isAtEnd() {
		var op ast.BinaryOperator
		switch p.current.Type {
		case lexer.GT:
			op = ast.OpGreaterThan
		case lexer.LT:
			op = ast.OpLessThan
		case lexer.GTE:
			op = ast.OpGreaterThanOrEqual
		case lexer.LTE:
			op = ast.OpLessThanOrEqual
		case lexer.EQ:
			op = ast.OpEqual
		case lexer.NEQ:
			op = ast.OpNotEqual
		default:
			return left, nil
		}

		p.advance()
		right, err := p.parseAdditive()
		if err != nil {
			return nil, err
		}
		left = &ast.BinaryExpression{
			Left:     left,
			Operator: op,
			Right:    right,
		}
	}

	return left, nil
}

// parseAdditive parses addition and subtraction expressions
func (p *Parser) parseAdditive() (ast.Expression, error) {
	left, err := p.parseMultiplicative()
	if err != nil {
		return nil, err
	}

	for !p.isAtEnd() && (p.current.Type == lexer.PLUS || p.current.Type == lexer.MINUS) {
		op := ast.OpPlus
		if p.current.Type == lexer.MINUS {
			op = ast.OpMinus
		}
		p.advance()

		right, err := p.parseMultiplicative()
		if err != nil {
			return nil, err
		}
		left = &ast.BinaryExpression{
			Left:     left,
			Operator: op,
			Right:    right,
		}
	}

	return left, nil
}

// parseMultiplicative parses multiplication and division expressions
func (p *Parser) parseMultiplicative() (ast.Expression, error) {
	left, err := p.parseUnary()
	if err != nil {
		return nil, err
	}

	for !p.isAtEnd() && (p.current.Type == lexer.MULTIPLY || p.current.Type == lexer.DIVIDE) {
		op := ast.OpMultiply
		if p.current.Type == lexer.DIVIDE {
			op = ast.OpDivide
		}
		p.advance()

		right, err := p.parseUnary()
		if err != nil {
			return nil, err
		}
		left = &ast.BinaryExpression{
			Left:     left,
			Operator: op,
			Right:    right,
		}
	}

	return left, nil
}

// parseUnary parses unary expressions
func (p *Parser) parseUnary() (ast.Expression, error) {
	if !p.isAtEnd() && p.current.Type == lexer.MINUS {
		p.advance()
		operand, err := p.parseUnary()
		if err != nil {
			return nil, err
		}
		return &ast.UnaryExpression{
			Operator: ast.OpUnaryMinus,
			Operand:  operand,
		}, nil
	}

	return p.parsePrimary()
}

// parsePrimary parses primary expressions (literals, identifiers, function calls, parentheses)
func (p *Parser) parsePrimary() (ast.Expression, error) {
	if p.isAtEnd() {
		return nil, p.error("unexpected end of input")
	}

	switch p.current.Type {
	case lexer.NUMBER:
		return p.parseNumber()
	case lexer.STRING, lexer.EXTERNAL_REFERENCE:
		return p.parseString()
	case lexer.IDENTIFIER, lexer.IF: // IF can be used as function name
		return p.parseIdentifierOrCall()
	case lexer.LPAREN:
		return p.parseGroupedExpression()
	default:
		return nil, p.error(fmt.Sprintf("unexpected token: %s", p.current.Type))
	}
}

// parseNumber parses a number literal
func (p *Parser) parseNumber() (*ast.NumberLiteral, error) {
	value, err := strconv.ParseFloat(p.current.Value, 64)
	if err != nil {
		return nil, p.error(fmt.Sprintf("invalid number: %s", p.current.Value))
	}
	p.advance()
	return &ast.NumberLiteral{Value: value}, nil
}

// parseString parses a string literal.
func (p *Parser) parseString() (*ast.StringLiteral, error) {
	value := p.current.Value
	external := p.current.Type == lexer.EXTERNAL_REFERENCE
	p.advance()
	return &ast.StringLiteral{Value: value, External: external}, nil
}

// parseIdentifierOrCall parses an identifier or function call
func (p *Parser) parseIdentifierOrCall() (ast.Expression, error) {
	name := p.current.Value
	p.advance()

	// Check if this is a function call
	if !p.isAtEnd() && p.current.Type == lexer.LPAREN {
		return p.parseFunctionCall(name)
	}

	// Just an identifier
	return &ast.Identifier{Name: name}, nil
}

// parseFunctionCall parses a function call
func (p *Parser) parseFunctionCall(name string) (*ast.FunctionCall, error) {
	p.advance() // consume '('

	args := make([]ast.Expression, 0)

	// Parse arguments
	for !p.isAtEnd() && p.current.Type != lexer.RPAREN {
		arg, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		args = append(args, arg)

		if p.current.Type == lexer.COMMA {
			p.advance()
		} else if p.current.Type != lexer.RPAREN {
			return nil, p.error("expected ',' or ')' in function call")
		}
	}

	if p.current.Type != lexer.RPAREN {
		return nil, p.error("expected ')' after function arguments")
	}
	p.advance() // consume ')'

	return &ast.FunctionCall{Name: name, Arguments: args}, nil
}

// parseGroupedExpression parses a parenthesized expression
func (p *Parser) parseGroupedExpression() (ast.Expression, error) {
	p.advance() // consume '('

	expr, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	if p.current.Type != lexer.RPAREN {
		return nil, p.error("expected ')' after expression")
	}
	p.advance() // consume ')'

	return expr, nil
}

func (p *Parser) parseStyleSuffixes() (*ast.DrawingStyle, error) {
	var style *ast.DrawingStyle

	for !p.isAtEnd() && p.current.Type == lexer.COMMA {
		next := p.peek()
		if next == nil || !p.isStyleToken(next) {
			return style, nil
		}

		p.advance() // consume comma
		if style == nil {
			style = &ast.DrawingStyle{}
		}

		if err := p.applyStyleToken(style, p.current); err != nil {
			return nil, err
		}
		p.advance()
	}

	return style, nil
}

func (p *Parser) isStyleToken(token *lexer.Token) bool {
	if token == nil {
		return false
	}

	switch token.Type {
	case lexer.COLOR, lexer.LINETHICK, lexer.DOTLINE, lexer.STICK:
		return true
	case lexer.IDENTIFIER:
		upper := strings.ToUpper(token.Value)
		return upper == "NODRAW" || upper == "COLORSTICK" || upper == "VOLSTICK"
	default:
		return false
	}
}

func (p *Parser) applyStyleToken(style *ast.DrawingStyle, token *lexer.Token) error {
	upper := strings.ToUpper(token.Value)

	switch token.Type {
	case lexer.COLOR:
		if upper == "COLORSTICK" {
			drawMethod := "colorstick"
			style.DrawMethod = &drawMethod
			return nil
		}
		color := upper
		if strings.HasPrefix(upper, "COLOR") && len(upper) > len("COLOR") {
			color = upper[len("COLOR"):]
		}
		style.Color = &color
	case lexer.LINETHICK:
		widthText := strings.TrimPrefix(upper, "LINETHICK")
		if widthText == "" {
			width := 1
			style.LineWidth = &width
			return nil
		}
		width, err := strconv.Atoi(widthText)
		if err != nil {
			return p.error(fmt.Sprintf("invalid line width: %s", token.Value))
		}
		style.LineWidth = &width
	case lexer.DOTLINE:
		lineStyle := "dotted"
		style.LineStyle = &lineStyle
	case lexer.STICK:
		drawMethod := strings.ToLower(upper)
		style.DrawMethod = &drawMethod
	case lexer.IDENTIFIER:
		switch upper {
		case "NODRAW":
			style.Hidden = true
		case "COLORSTICK", "VOLSTICK":
			drawMethod := strings.ToLower(upper)
			style.DrawMethod = &drawMethod
		default:
			return p.error(fmt.Sprintf("unexpected style token: %s", token.Value))
		}
	}

	return nil
}

// advance moves to the next token
func (p *Parser) advance() {
	if !p.isAtEnd() {
		p.pos++
		if p.pos < len(p.tokens) {
			p.current = p.tokens[p.pos]
		}
	}
}

// peek returns the next token without advancing
func (p *Parser) peek() *lexer.Token {
	if p.pos+1 < len(p.tokens) {
		return p.tokens[p.pos+1]
	}
	return nil
}

// isAtEnd checks if we've reached the end of tokens
func (p *Parser) isAtEnd() bool {
	return p.current == nil || p.current.Type == lexer.EOF
}

// error creates a parser error with current position
func (p *Parser) error(message string) error {
	if p.current != nil {
		return errors.NewParserError(message, p.current.Line, p.current.Column)
	}
	return errors.NewParserError(message, 0, 0)
}
