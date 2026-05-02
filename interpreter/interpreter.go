// Package interpreter provides the execution engine for formula ASTs
package interpreter

import (
	"fmt"
	"math"

	"github.com/DTrader-store/formula-go/errors"
	"github.com/DTrader-store/formula-go/parser/ast"
	"github.com/DTrader-store/formula-go/types"
)

// Value represents a computed value (can be single value or array)
type Value struct {
	Single   float64               // Single value
	Array    []float64             // Array of values
	Text     string                // String value for text/drawing functions
	Drawings []*types.DrawingEvent // Drawing annotations
	IsArray  bool                  // Whether this is an array value
	IsString bool                  // Whether this is a string value
	IsDraw   bool                  // Whether this is a drawing result
}

// NewSingleValue creates a single value
func NewSingleValue(v float64) *Value {
	return &Value{Single: v, IsArray: false}
}

// NewArrayValue creates an array value
func NewArrayValue(arr []float64) *Value {
	return &Value{Array: arr, IsArray: true}
}

// NewStringValue creates a string value.
func NewStringValue(text string) *Value {
	return &Value{Text: text, IsString: true}
}

// NewDrawingValue creates a drawing result value.
func NewDrawingValue(drawings []*types.DrawingEvent) *Value {
	return &Value{Drawings: drawings, IsDraw: true}
}

// Interpreter executes formula ASTs
type Interpreter struct {
	marketData []*types.MarketData
	variables  map[string]*Value
	userVars   []string // Track user-defined variables in order
	outputs    []outputDeclaration
	functions  *FunctionRegistry
	exprCount  int
}

type outputDeclaration struct {
	Name  string
	Value *Value
	Style *types.LineStyle
}

// NewInterpreter creates a new Interpreter
func NewInterpreter(marketData []*types.MarketData) *Interpreter {
	return &Interpreter{
		marketData: marketData,
		variables:  make(map[string]*Value),
		userVars:   make([]string, 0),
		outputs:    make([]outputDeclaration, 0),
		functions:  NewFunctionRegistry(),
	}
}

// Execute executes a program and returns the result
func (interp *Interpreter) Execute(program *ast.Program) (*types.FormulaResult, error) {
	// Initialize market data variables
	interp.initMarketDataVariables()

	// Execute all statements
	for _, stmt := range program.Body {
		if err := interp.executeStatement(stmt); err != nil {
			return nil, err
		}
	}

	// Build result
	return interp.buildResult(), nil
}

// initMarketDataVariables initializes built-in market data variables
func (interp *Interpreter) initMarketDataVariables() {
	if len(interp.marketData) == 0 {
		return
	}

	n := len(interp.marketData)
	open := make([]float64, n)
	close := make([]float64, n)
	high := make([]float64, n)
	low := make([]float64, n)
	volume := make([]float64, n)
	amount := make([]float64, n)

	for i, data := range interp.marketData {
		open[i] = data.Open
		close[i] = data.Close
		high[i] = data.High
		low[i] = data.Low
		volume[i] = data.Volume
		amount[i] = data.Amount
	}

	interp.variables["OPEN"] = NewArrayValue(open)
	interp.variables["CLOSE"] = NewArrayValue(close)
	interp.variables["HIGH"] = NewArrayValue(high)
	interp.variables["LOW"] = NewArrayValue(low)
	interp.variables["VOLUME"] = NewArrayValue(volume)
	interp.variables["VOL"] = interp.variables["VOLUME"]
	interp.variables["V"] = interp.variables["VOLUME"]
	interp.variables["AMOUNT"] = NewArrayValue(amount)
	interp.variables["AMO"] = interp.variables["AMOUNT"]
	interp.variables["O"] = interp.variables["OPEN"]
	interp.variables["C"] = interp.variables["CLOSE"]
	interp.variables["H"] = interp.variables["HIGH"]
	interp.variables["L"] = interp.variables["LOW"]
}

// executeStatement executes a single statement
func (interp *Interpreter) executeStatement(stmt ast.Statement) error {
	switch s := stmt.(type) {
	case *ast.VariableDeclaration:
		return interp.executeVariableDeclaration(s)
	case *ast.OutputDeclaration:
		return interp.executeOutputDeclaration(s)
	case *ast.ExpressionStatement:
		// For standalone expressions, evaluate and add to output with temp name
		value, err := interp.evaluateExpression(s.Expr)
		if err != nil {
			return err
		}
		// Generate a stable unique name for standalone expressions.
		name := fmt.Sprintf("__expr__%d", interp.exprCount)
		interp.exprCount++
		if ident, ok := s.Expr.(*ast.Identifier); ok {
			name = ident.Name
		}
		interp.variables[name] = value
		interp.userVars = append(interp.userVars, name)
		return nil
	default:
		return errors.NewRuntimeError(fmt.Sprintf("unknown statement type: %T", stmt))
	}
}

// executeVariableDeclaration executes a variable declaration
func (interp *Interpreter) executeVariableDeclaration(decl *ast.VariableDeclaration) error {
	value, err := interp.evaluateExpression(decl.Value)
	if err != nil {
		return err
	}
	interp.variables[decl.Name] = value
	interp.userVars = append(interp.userVars, decl.Name) // Preserve order
	return nil
}

// executeOutputDeclaration executes an output declaration
func (interp *Interpreter) executeOutputDeclaration(decl *ast.OutputDeclaration) error {
	value, err := interp.evaluateExpression(decl.Value)
	if err != nil {
		return err
	}
	interp.variables[decl.Name] = value
	interp.outputs = append(interp.outputs, outputDeclaration{
		Name:  decl.Name,
		Value: value,
		Style: lineStyleFromAST(decl.Style),
	})
	return nil
}

// evaluateExpression evaluates an expression and returns a value
func (interp *Interpreter) evaluateExpression(expr ast.Expression) (*Value, error) {
	switch e := expr.(type) {
	case *ast.NumberLiteral:
		return interp.evaluateNumberLiteral(e)
	case *ast.StringLiteral:
		return interp.evaluateStringLiteral(e)
	case *ast.Identifier:
		return interp.evaluateIdentifier(e)
	case *ast.BinaryExpression:
		return interp.evaluateBinaryExpression(e)
	case *ast.UnaryExpression:
		return interp.evaluateUnaryExpression(e)
	case *ast.FunctionCall:
		return interp.evaluateFunctionCall(e)
	default:
		return nil, errors.NewRuntimeError(fmt.Sprintf("unknown expression type: %T", expr))
	}
}

// evaluateNumberLiteral evaluates a number literal
func (interp *Interpreter) evaluateNumberLiteral(lit *ast.NumberLiteral) (*Value, error) {
	return NewSingleValue(lit.Value), nil
}

// evaluateStringLiteral evaluates quoted external references when already available as variables.
func (interp *Interpreter) evaluateStringLiteral(lit *ast.StringLiteral) (*Value, error) {
	if lit.External {
		value, exists := interp.variables[lit.Value]
		if !exists {
			return nil, errors.NewRuntimeError(fmt.Sprintf("undefined external reference: %s", lit.Value))
		}
		return value, nil
	}
	return NewStringValue(lit.Value), nil
}

// evaluateIdentifier evaluates an identifier
func (interp *Interpreter) evaluateIdentifier(id *ast.Identifier) (*Value, error) {
	value, exists := interp.variables[id.Name]
	if !exists {
		return nil, errors.NewRuntimeError(fmt.Sprintf("undefined variable: %s", id.Name))
	}
	return value, nil
}

// evaluateBinaryExpression evaluates a binary expression
func (interp *Interpreter) evaluateBinaryExpression(expr *ast.BinaryExpression) (*Value, error) {
	left, err := interp.evaluateExpression(expr.Left)
	if err != nil {
		return nil, err
	}

	right, err := interp.evaluateExpression(expr.Right)
	if err != nil {
		return nil, err
	}

	// Handle array operations
	if left.IsArray && right.IsArray {
		return interp.binaryOpArrayArray(expr.Operator, left.Array, right.Array)
	} else if left.IsArray {
		return interp.binaryOpArrayScalar(expr.Operator, left.Array, right.Single)
	} else if right.IsArray {
		return interp.binaryOpScalarArray(expr.Operator, left.Single, right.Array)
	} else {
		return interp.binaryOpScalarScalar(expr.Operator, left.Single, right.Single)
	}
}

// binaryOpScalarScalar performs binary operation on two scalars
func (interp *Interpreter) binaryOpScalarScalar(op ast.BinaryOperator, a, b float64) (*Value, error) {
	var result float64
	switch op {
	case ast.OpPlus:
		result = a + b
	case ast.OpMinus:
		result = a - b
	case ast.OpMultiply:
		result = a * b
	case ast.OpDivide:
		if b == 0 {
			return nil, errors.NewRuntimeError("division by zero")
		}
		result = a / b
	case ast.OpGreaterThan:
		if a > b {
			result = 1
		} else {
			result = 0
		}
	case ast.OpLessThan:
		if a < b {
			result = 1
		} else {
			result = 0
		}
	case ast.OpGreaterThanOrEqual:
		if a >= b {
			result = 1
		} else {
			result = 0
		}
	case ast.OpLessThanOrEqual:
		if a <= b {
			result = 1
		} else {
			result = 0
		}
	case ast.OpEqual:
		if math.Abs(a-b) < 1e-10 {
			result = 1
		} else {
			result = 0
		}
	case ast.OpNotEqual:
		if math.Abs(a-b) >= 1e-10 {
			result = 1
		} else {
			result = 0
		}
	case ast.OpAnd:
		if a != 0 && b != 0 {
			result = 1
		} else {
			result = 0
		}
	case ast.OpOr:
		if a != 0 || b != 0 {
			result = 1
		} else {
			result = 0
		}
	default:
		return nil, errors.NewRuntimeError(fmt.Sprintf("unknown binary operator: %s", op))
	}

	return NewSingleValue(result), nil
}

// binaryOpArrayArray performs binary operation on two arrays
func (interp *Interpreter) binaryOpArrayArray(op ast.BinaryOperator, a, b []float64) (*Value, error) {
	if len(a) != len(b) {
		return nil, errors.NewRuntimeError("array length mismatch")
	}

	result := make([]float64, len(a))
	for i := range a {
		val, err := interp.binaryOpScalarScalar(op, a[i], b[i])
		if err != nil {
			return nil, err
		}
		result[i] = val.Single
	}

	return NewArrayValue(result), nil
}

// binaryOpArrayScalar performs binary operation on array and scalar
func (interp *Interpreter) binaryOpArrayScalar(op ast.BinaryOperator, arr []float64, scalar float64) (*Value, error) {
	result := make([]float64, len(arr))
	for i, v := range arr {
		val, err := interp.binaryOpScalarScalar(op, v, scalar)
		if err != nil {
			return nil, err
		}
		result[i] = val.Single
	}
	return NewArrayValue(result), nil
}

// binaryOpScalarArray performs binary operation on scalar and array
func (interp *Interpreter) binaryOpScalarArray(op ast.BinaryOperator, scalar float64, arr []float64) (*Value, error) {
	result := make([]float64, len(arr))
	for i, v := range arr {
		val, err := interp.binaryOpScalarScalar(op, scalar, v)
		if err != nil {
			return nil, err
		}
		result[i] = val.Single
	}
	return NewArrayValue(result), nil
}

// evaluateUnaryExpression evaluates a unary expression
func (interp *Interpreter) evaluateUnaryExpression(expr *ast.UnaryExpression) (*Value, error) {
	operand, err := interp.evaluateExpression(expr.Operand)
	if err != nil {
		return nil, err
	}

	if expr.Operator == ast.OpUnaryMinus {
		if operand.IsArray {
			result := make([]float64, len(operand.Array))
			for i, v := range operand.Array {
				result[i] = -v
			}
			return NewArrayValue(result), nil
		}
		return NewSingleValue(-operand.Single), nil
	}

	return nil, errors.NewRuntimeError(fmt.Sprintf("unknown unary operator: %s", expr.Operator))
}

// evaluateFunctionCall evaluates a function call
func (interp *Interpreter) evaluateFunctionCall(call *ast.FunctionCall) (*Value, error) {
	// Evaluate arguments
	args := make([]*Value, len(call.Arguments))
	for i, arg := range call.Arguments {
		val, err := interp.evaluateExpression(arg)
		if err != nil {
			return nil, err
		}
		args[i] = val
	}

	// Call function
	return interp.functions.Call(call.Name, args, interp.marketData)
}

// buildResult builds the final formula result
func (interp *Interpreter) buildResult() *types.FormulaResult {
	result := types.NewFormulaResult()
	addValueToResult := func(name string, value *Value, style *types.LineStyle) {
		if value.IsDraw {
			for _, drawing := range value.Drawings {
				result.AddDrawing(drawing)
			}
		} else if value.IsArray {
			result.AddOutput(name, value.Array, style)
		} else if !value.IsString {
			result.SetVariable(name, value.Single)
		}
	}

	if len(interp.outputs) > 0 {
		for _, output := range interp.outputs {
			addValueToResult(output.Name, output.Value, output.Style)
		}
		for _, name := range interp.userVars {
			value := interp.variables[name]
			if value.IsDraw {
				addValueToResult(name, value, nil)
			}
		}
		return result
	}

	// Add user-defined variables to result in order
	for _, name := range interp.userVars {
		addValueToResult(name, interp.variables[name], nil)
	}

	return result
}

func lineStyleFromAST(style *ast.DrawingStyle) *types.LineStyle {
	if style == nil {
		return nil
	}

	lineStyle := &types.LineStyle{
		Hidden: style.Hidden,
	}
	if style.Color != nil {
		lineStyle.Color = *style.Color
	}
	if style.LineWidth != nil {
		lineStyle.LineWidth = *style.LineWidth
	}
	if style.LineStyle != nil {
		lineStyle.LineStyle = *style.LineStyle
	}
	if style.DrawMethod != nil {
		lineStyle.DrawMethod = *style.DrawMethod
	}

	return lineStyle
}
