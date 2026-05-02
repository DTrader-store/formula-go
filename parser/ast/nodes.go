// Package ast defines the Abstract Syntax Tree node types and interfaces
package ast

// NodeType represents the type of an AST node
type NodeType string

// NodeType constants
const (
	// Program and Statements
	ProgramNode             NodeType = "Program"
	VariableDeclarationNode NodeType = "VariableDeclaration"
	OutputDeclarationNode   NodeType = "OutputDeclaration"
	ExpressionStatementNode NodeType = "ExpressionStatement"

	// Expressions
	BinaryExpressionNode      NodeType = "BinaryExpression"
	UnaryExpressionNode       NodeType = "UnaryExpression"
	FunctionCallNode          NodeType = "FunctionCall"
	ConditionalExpressionNode NodeType = "ConditionalExpression"

	// Literals and Identifiers
	IdentifierNode    NodeType = "Identifier"
	NumberLiteralNode NodeType = "NumberLiteral"
	StringLiteralNode NodeType = "StringLiteral"
)

// BinaryOperator represents binary operators
type BinaryOperator string

const (
	// Arithmetic
	OpPlus     BinaryOperator = "+"
	OpMinus    BinaryOperator = "-"
	OpMultiply BinaryOperator = "*"
	OpDivide   BinaryOperator = "/"
	OpModulo   BinaryOperator = "%"
	OpPower    BinaryOperator = "^"

	// Comparison
	OpEqual              BinaryOperator = "=="
	OpNotEqual           BinaryOperator = "!="
	OpLessThan           BinaryOperator = "<"
	OpLessThanOrEqual    BinaryOperator = "<="
	OpGreaterThan        BinaryOperator = ">"
	OpGreaterThanOrEqual BinaryOperator = ">="

	// Logical
	OpAnd BinaryOperator = "&&"
	OpOr  BinaryOperator = "||"
)

// UnaryOperator represents unary operators
type UnaryOperator string

const (
	OpUnaryMinus UnaryOperator = "-"
	OpNot        UnaryOperator = "!"
)

// DrawingStyle represents the drawing style configuration for output declarations
type DrawingStyle struct {
	Color      *string
	LineWidth  *int
	LineStyle  *string
	Hidden     bool
	DrawMethod *string
}

// Node is the base interface for all AST nodes
type Node interface {
	Type() NodeType
}

// Expression interface - all expressions are also statements
type Expression interface {
	Node
	exprNode()
}

// Statement interface
type Statement interface {
	Node
	stmtNode()
}

// Program node - root of the AST
type Program struct {
	Body []Statement
}

func (p *Program) Type() NodeType { return ProgramNode }
func (p *Program) stmtNode()      {}

// VariableDeclaration represents: var x = 10;
type VariableDeclaration struct {
	Name  string
	Value Expression
}

func (v *VariableDeclaration) Type() NodeType { return VariableDeclarationNode }
func (v *VariableDeclaration) stmtNode()      {}

// OutputDeclaration represents: output result = x + y; or output result = x + y; [color: red, size: 14];
type OutputDeclaration struct {
	Name  string
	Value Expression
	Style *DrawingStyle
}

func (o *OutputDeclaration) Type() NodeType { return OutputDeclarationNode }
func (o *OutputDeclaration) stmtNode()      {}

// ExpressionStatement represents: expression;
type ExpressionStatement struct {
	Expr Expression
}

func (e *ExpressionStatement) Type() NodeType { return ExpressionStatementNode }
func (e *ExpressionStatement) stmtNode()      {}

// BinaryExpression represents: left operator right
type BinaryExpression struct {
	Left     Expression
	Operator BinaryOperator
	Right    Expression
}

func (b *BinaryExpression) Type() NodeType { return BinaryExpressionNode }
func (b *BinaryExpression) exprNode()      {}

// UnaryExpression represents: operator operand
type UnaryExpression struct {
	Operator UnaryOperator
	Operand  Expression
}

func (u *UnaryExpression) Type() NodeType { return UnaryExpressionNode }
func (u *UnaryExpression) exprNode()      {}

// FunctionCall represents: functionName(arg1, arg2, ...)
type FunctionCall struct {
	Name      string
	Arguments []Expression
}

func (f *FunctionCall) Type() NodeType { return FunctionCallNode }
func (f *FunctionCall) exprNode()      {}

// ConditionalExpression represents: test ? consequent : alternate
type ConditionalExpression struct {
	Test       Expression
	Consequent Expression
	Alternate  Expression
}

func (c *ConditionalExpression) Type() NodeType { return ConditionalExpressionNode }
func (c *ConditionalExpression) exprNode()      {}

// Identifier represents: variable or function name reference
type Identifier struct {
	Name string
}

func (i *Identifier) Type() NodeType { return IdentifierNode }
func (i *Identifier) exprNode()      {}

// NumberLiteral represents: numeric constant
type NumberLiteral struct {
	Value float64
}

func (n *NumberLiteral) Type() NodeType { return NumberLiteralNode }
func (n *NumberLiteral) exprNode()      {}

// StringLiteral represents a quoted text or external indicator reference.
type StringLiteral struct {
	Value    string
	External bool
}

func (s *StringLiteral) Type() NodeType { return StringLiteralNode }
func (s *StringLiteral) exprNode()      {}
