// Package formula provides a formula parser and interpreter for technical indicators
package formula

import (
	"github.com/DTrader-store/formula-go/engine"
	"github.com/DTrader-store/formula-go/errors"
	"github.com/DTrader-store/formula-go/interpreter"
	"github.com/DTrader-store/formula-go/lexer"
	"github.com/DTrader-store/formula-go/parser"
	"github.com/DTrader-store/formula-go/parser/ast"
	"github.com/DTrader-store/formula-go/types"
)

// VERSION is the current version of formula-go
const VERSION = "1.0.0"

// Export error types for external use
type (
	FormulaError = errors.FormulaError
	LexerError   = errors.LexerError
	ParserError  = errors.ParserError
	RuntimeError = errors.RuntimeError
)

// Export lexer types for external use
type (
	Token     = lexer.Token
	TokenType = lexer.TokenType
	Lexer     = lexer.Lexer
)

// Export parser types for external use
type (
	Parser = parser.Parser
)

// Export AST types for external use
type (
	Node       = ast.Node
	Expression = ast.Expression
	Statement  = ast.Statement
	Program    = ast.Program
)

// Export types for external use
type (
	MarketData    = types.MarketData
	FormulaResult = types.FormulaResult
	OutputLine    = types.OutputLine
	LineStyle     = types.LineStyle
	DrawingEvent  = types.DrawingEvent
)

// Export engine types
type (
	FormulaEngine = engine.FormulaEngine
)

// Export interpreter types
type (
	Value            = interpreter.Value
	Interpreter      = interpreter.Interpreter
	FunctionRegistry = interpreter.FunctionRegistry
)

// Constructor functions
var (
	NewMarketData    = types.NewMarketData
	NewFormulaResult = types.NewFormulaResult
	NewLexerError    = errors.NewLexerError
	NewParserError   = errors.NewParserError
	NewRuntimeError  = errors.NewRuntimeError
	NewLexer         = lexer.NewLexer
	NewParser        = parser.NewParser
	NewFormulaEngine = engine.NewFormulaEngine
	NewInterpreter   = interpreter.NewInterpreter
)
