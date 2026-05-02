package engine

import (
	"math"
	"strings"
	"testing"

	"github.com/DTrader-store/formula-go/interpreter"
)

type builtinArityCase struct {
	validArgs []string
	arity     []int
}

func TestRegisteredBuiltinsHaveArityCases(t *testing.T) {
	registryNames := interpreter.NewFunctionRegistry().Names()
	cases := builtinArityCases()

	for _, name := range registryNames {
		if _, ok := cases[name]; !ok {
			t.Errorf("registered builtin %s has no arity case", name)
		}
	}
	for name := range cases {
		if !stringInSlice(name, registryNames) {
			t.Errorf("arity case %s is not a registered builtin", name)
		}
	}
}

func TestBuiltinArityErrors(t *testing.T) {
	engine := NewFormulaEngine()
	marketData := createTestData()

	for _, name := range interpreter.NewFunctionRegistry().Names() {
		tc := builtinArityCases()[name]
		if len(tc.arity) == 0 {
			t.Fatalf("builtin %s has no valid arity", name)
		}
		for _, argCount := range invalidArityCounts(tc.arity) {
			t.Run(name+"_argc_"+argCountName(argCount), func(t *testing.T) {
				formula := "X: " + name + "(" + strings.Join(argsForCount(tc.validArgs, argCount), ", ") + ")"
				_, err := engine.Run(formula, marketData)
				if err == nil {
					t.Fatalf("expected arity error for %s with %d args", name, argCount)
				}
				if !strings.Contains(err.Error(), name+" requires ") {
					t.Fatalf("expected %s arity error, got %q", name, err.Error())
				}
			})
		}
	}
}

func TestBuiltinRuntimeErrors(t *testing.T) {
	tests := []struct {
		name        string
		formula     string
		wantMessage string
	}{
		{
			name:        "undefined variable",
			formula:     "X: UNKNOWN_VAR",
			wantMessage: "undefined variable: UNKNOWN_VAR",
		},
		{
			name:        "undefined function",
			formula:     "X: NO_SUCH_FUNC(C)",
			wantMessage: "undefined function: NO_SUCH_FUNC",
		},
		{
			name:        "external reference must exist",
			formula:     `X: "MACD.DIF#WEEK"`,
			wantMessage: "undefined external reference: MACD.DIF#WEEK",
		},
		{
			name:        "division by zero",
			formula:     "X: 1 / 0",
			wantMessage: "division by zero",
		},
		{
			name:        "rolling period above data length",
			formula:     "X: MA(C, 11)",
			wantMessage: "MA period must be between 1 and 10",
		},
		{
			name:        "rolling period must be positive",
			formula:     "X: SUM(C, 0)",
			wantMessage: "SUM period must be between 1 and 10",
		},
		{
			name:        "wrong argument count",
			formula:     "X: REF(C)",
			wantMessage: "REF requires 2 arguments",
		},
		{
			name:        "first argument must be array",
			formula:     "X: MA(1, 2)",
			wantMessage: "MA first argument must be an array",
		},
		{
			name:        "period argument must be scalar",
			formula:     "X: MA(C, C)",
			wantMessage: "MA second argument must be a number",
		},
		{
			name:        "future reference period must be non-negative",
			formula:     "X: REFX(C, -1)",
			wantMessage: "REFX period must be non-negative",
		},
		{
			name:        "last window must be ordered",
			formula:     "X: LAST(C > O, 1, 2)",
			wantMessage: "LAST requires from >= to >= 0",
		},
		{
			name:        "drawing condition must be array",
			formula:     "X := DRAWTEXT(1, C, 'UP')",
			wantMessage: "DRAWTEXT first argument must be an array",
		},
		{
			name:        "zero-argument function rejects arguments",
			formula:     "X: CURRBARSCOUNT(C)",
			wantMessage: "CURRBARSCOUNT requires 0 arguments",
		},
		{
			name:        "string rejected by numeric function",
			formula:     "X: ABS('UP')",
			wantMessage: "ABS argument must be numeric",
		},
	}

	engine := NewFormulaEngine()
	marketData := createTestData()
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := engine.Run(tc.formula, marketData)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if !strings.Contains(err.Error(), tc.wantMessage) {
				t.Fatalf("expected error containing %q, got %q", tc.wantMessage, err.Error())
			}
		})
	}
}

func TestPeriodFunctionBoundaryErrors(t *testing.T) {
	tests := []struct {
		name        string
		formula     string
		wantMessage string
	}{
		{name: "EMA period above data length", formula: "X: EMA(C, 11)", wantMessage: "EMA period must be between 1 and 10"},
		{name: "HHV period zero", formula: "X: HHV(C, 0)", wantMessage: "HHV period must be between 1 and 10"},
		{name: "LLV period scalar required", formula: "X: LLV(C, O)", wantMessage: "LLV second argument must be a number"},
		{name: "STD period above data length", formula: "X: STD(C, 11)", wantMessage: "STD period must be between 1 and 10"},
		{name: "VAR period zero", formula: "X: VAR(C, 0)", wantMessage: "VAR period must be between 1 and 10"},
		{name: "WMA period above data length", formula: "X: WMA(C, 11)", wantMessage: "WMA period must be between 1 and 10"},
		{name: "COUNT period zero", formula: "X: COUNT(C > O, 0)", wantMessage: "COUNT period must be between 1 and 10"},
		{name: "EVERY period scalar required", formula: "X: EVERY(C > O, C)", wantMessage: "EVERY second argument must be a number"},
		{name: "EXIST period above data length", formula: "X: EXIST(C > O, 11)", wantMessage: "EXIST period must be between 1 and 10"},
		{name: "HHVBARS period zero", formula: "X: HHVBARS(H, 0)", wantMessage: "HHVBARS period must be between 1 and 10"},
		{name: "LLVBARS period scalar required", formula: "X: LLVBARS(L, C)", wantMessage: "LLVBARS second argument must be a number"},
		{name: "AVEDEV period above data length", formula: "X: AVEDEV(C, 11)", wantMessage: "AVEDEV period must be between 1 and 10"},
		{name: "FILTER period zero", formula: "X: FILTER(C > O, 0)", wantMessage: "FILTER period must be positive"},
		{name: "UPNDAY period zero", formula: "X: UPNDAY(C, 0)", wantMessage: "UPNDAY period must be positive"},
		{name: "DOWNNDAY period scalar required", formula: "X: DOWNNDAY(C, C)", wantMessage: "DOWNNDAY second argument must be a number"},
		{name: "NDAY period zero", formula: "X: NDAY(C, O, 0)", wantMessage: "NDAY period must be positive"},
		{name: "LONGCROSS period zero", formula: "X: LONGCROSS(C, O, 0)", wantMessage: "LONGCROSS period must be positive"},
		{name: "SMA period zero", formula: "X: SMA(C, 0, 1)", wantMessage: "SMA period must be positive"},
		{name: "SMA weight above period", formula: "X: SMA(C, 3, 4)", wantMessage: "SMA weight must be between 1 and period"},
		{name: "DEVSQ period above data length", formula: "X: DEVSQ(C, 11)", wantMessage: "DEVSQ period must be between 1 and 10"},
		{name: "STDDEV period scalar required", formula: "X: STDDEV(C, C)", wantMessage: "STDDEV second argument must be a number"},
		{name: "COVAR period above data length", formula: "X: COVAR(C, O, 11)", wantMessage: "COVAR period must be between 1 and 10"},
		{name: "RELATE period scalar required", formula: "X: RELATE(C, O, C)", wantMessage: "RELATE third argument must be a number"},
	}

	runFormulaErrorCases(t, tests)
}

func TestBuiltinTypeBoundaryErrors(t *testing.T) {
	tests := []struct {
		name        string
		formula     string
		wantMessage string
	}{
		{name: "MAX rejects strings", formula: "X: MAX('UP', 1)", wantMessage: "MAX arguments must be numeric"},
		{name: "MIN rejects strings", formula: "X: MIN(1, 'UP')", wantMessage: "MIN arguments must be numeric"},
		{name: "ROUND2 digit argument scalar required", formula: "X: ROUND2(C, C)", wantMessage: "ROUND2 second argument must be a number"},
		{name: "BETWEEN array bound requires value array", formula: "X: BETWEEN(1, C, 110)", wantMessage: "BETWEEN: value must be array when using array bounds"},
		{name: "CROSS requires arrays", formula: "X: CROSS(1, C)", wantMessage: "CROSS requires array arguments"},
		{name: "LONGCROSS requires arrays", formula: "X: LONGCROSS(1, C, 1)", wantMessage: "LONGCROSS first two arguments must be arrays"},
		{name: "LAST condition requires array", formula: "X: LAST(1, 2, 1)", wantMessage: "LAST first argument must be an array"},
		{name: "EXISTR window args require scalars", formula: "X: EXISTR(C > O, C, 1)", wantMessage: "EXISTR second and third arguments must be numbers"},
		{name: "DMA first arg requires array", formula: "X: DMA(1, 0.5)", wantMessage: "DMA first argument must be an array"},
		{name: "VALUEWHEN condition requires array", formula: "X: VALUEWHEN(1, C)", wantMessage: "VALUEWHEN first argument must be an array"},
		{name: "SUMBARS data requires array", formula: "X: SUMBARS(1, 10)", wantMessage: "SUMBARS first argument must be an array"},
		{name: "STICKLINE condition requires array", formula: "X := STICKLINE(1, O, C, 2, 0)", wantMessage: "STICKLINE first argument must be an array"},
		{name: "STICKLINE rejects string", formula: "X := STICKLINE(C > O, O, 'C', 2, 0)", wantMessage: "STICKLINE arguments must be numeric"},
		{name: "DRAWLINE first condition requires array", formula: "X := DRAWLINE(1, L, C > O, H, 0)", wantMessage: "DRAWLINE first argument must be an array"},
		{name: "DRAWLINE second condition requires array", formula: "X := DRAWLINE(C > O, L, 1, H, 0)", wantMessage: "DRAWLINE third argument must be an array"},
		{name: "POLYLINE condition requires array", formula: "X := POLYLINE(1, C)", wantMessage: "POLYLINE first argument must be an array"},
		{name: "POLYLINE rejects string price", formula: "X := POLYLINE(C > O, 'C')", wantMessage: "POLYLINE arguments must be numeric"},
		{name: "DRAWBAND rejects string", formula: "X := DRAWBAND(H, 1, 'LOW', 2)", wantMessage: "DRAWBAND arguments must be numeric"},
		{name: "DRAWKLINE rejects string", formula: "X := DRAWKLINE(H, O, L, 'C')", wantMessage: "DRAWKLINE arguments must be numeric"},
	}

	runFormulaErrorCases(t, tests)
}

func TestBuiltinNaNPropagationAndTruthiness(t *testing.T) {
	engine := NewFormulaEngine()
	marketData := createTestData()

	formula := `
		NULL_VALUE: DRAWNULL()
		NULL_LINE: IF(C > O, C, DRAWNULL())
		NULL_PLUS_ONE: NULL_LINE + 1
		NULL_NOT: NOT(NULL_LINE)
		NULL_COUNT: COUNT(NULL_LINE, 2)
		NULL_EVERY: EVERY(NULL_LINE, 2)
		NULL_EXIST: EXIST(NULL_LINE, 2)
		NULL_BARS: BARSLAST(NULL_LINE)
	`

	result, err := engine.Run(formula, marketData)
	if err != nil {
		t.Fatalf("Run() error: %v", err)
	}

	assertVariableIsNaN(t, result, "NULL_VALUE")
	outputs := outputsByName(result)
	assertSeries(t, outputs["NULL_LINE"], []float64{105, math.NaN(), 107, 110, math.NaN(), 111, math.NaN(), 112, 115, math.NaN()})
	assertSeries(t, outputs["NULL_PLUS_ONE"], []float64{106, math.NaN(), 108, 111, math.NaN(), 112, math.NaN(), 113, 116, math.NaN()})
	assertSeries(t, outputs["NULL_NOT"], []float64{0, 1, 0, 0, 1, 0, 1, 0, 0, 1})
	assertSeries(t, outputs["NULL_COUNT"], []float64{math.NaN(), 1, 1, 2, 1, 1, 1, 1, 2, 1})
	assertSeries(t, outputs["NULL_EVERY"], []float64{0, 0, 0, 1, 0, 0, 0, 0, 1, 0})
	assertSeries(t, outputs["NULL_EXIST"], []float64{0, 1, 1, 1, 1, 1, 1, 1, 1, 1})
	assertSeries(t, outputs["NULL_BARS"], []float64{0, 1, 0, 0, 1, 0, 1, 0, 0, 1})
}

func TestRollingFunctionsPropagateNaNWindow(t *testing.T) {
	engine := NewFormulaEngine()
	marketData := createTestData()

	formula := `
		NULL_LINE: IF(C > O, C, DRAWNULL())
		MA_NULL: MA(NULL_LINE, 2)
		STD_NULL: STD(NULL_LINE, 2)
		DEV_NULL: DEVSQ(NULL_LINE, 2)
	`

	result, err := engine.Run(formula, marketData)
	if err != nil {
		t.Fatalf("Run() error: %v", err)
	}

	outputs := outputsByName(result)
	assertSeries(t, outputs["MA_NULL"], []float64{math.NaN(), math.NaN(), math.NaN(), 108.5, math.NaN(), math.NaN(), math.NaN(), math.NaN(), 113.5, math.NaN()})
	assertSeries(t, outputs["STD_NULL"], []float64{math.NaN(), math.NaN(), math.NaN(), 1.5, math.NaN(), math.NaN(), math.NaN(), math.NaN(), 1.5, math.NaN()})
	assertSeries(t, outputs["DEV_NULL"], []float64{math.NaN(), math.NaN(), math.NaN(), 4.5, math.NaN(), math.NaN(), math.NaN(), math.NaN(), 4.5, math.NaN()})
}

func runFormulaErrorCases(t *testing.T, tests []struct {
	name        string
	formula     string
	wantMessage string
}) {
	t.Helper()
	engine := NewFormulaEngine()
	marketData := createTestData()
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := engine.Run(tc.formula, marketData)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if !strings.Contains(err.Error(), tc.wantMessage) {
				t.Fatalf("expected error containing %q, got %q", tc.wantMessage, err.Error())
			}
		})
	}
}

func builtinArityCases() map[string]builtinArityCase {
	unaryNumeric := []string{"1"}
	binaryNumeric := []string{"1", "2"}
	seriesPeriod := []string{"C", "2"}
	pairPeriod := []string{"C", "O", "2"}
	conditionPeriod := []string{"C > O", "2"}
	conditionWindow := []string{"C > O", "2", "1"}

	return map[string]builtinArityCase{
		"ABS":            {validArgs: unaryNumeric, arity: []int{1}},
		"ACOS":           {validArgs: unaryNumeric, arity: []int{1}},
		"ASIN":           {validArgs: unaryNumeric, arity: []int{1}},
		"ATAN":           {validArgs: unaryNumeric, arity: []int{1}},
		"AVEDEV":         {validArgs: seriesPeriod, arity: []int{2}},
		"BARSCOUNT":      {validArgs: []string{"C"}, arity: []int{1}},
		"BARSLAST":       {validArgs: []string{"C > O"}, arity: []int{1}},
		"BARSLASTCOUNT":  {validArgs: []string{"C > O"}, arity: []int{1}},
		"BARSSINCE":      {validArgs: []string{"C > O"}, arity: []int{1}},
		"BARSTATUS":      {validArgs: nil, arity: []int{0}},
		"BETA":           {validArgs: pairPeriod, arity: []int{3}},
		"BETWEEN":        {validArgs: []string{"C", "100", "110"}, arity: []int{3}},
		"CEILING":        {validArgs: unaryNumeric, arity: []int{1}},
		"CONST":          {validArgs: []string{"C"}, arity: []int{1}},
		"COS":            {validArgs: unaryNumeric, arity: []int{1}},
		"COUNT":          {validArgs: conditionPeriod, arity: []int{2}},
		"COVAR":          {validArgs: pairPeriod, arity: []int{3}},
		"CROSS":          {validArgs: []string{"C", "O"}, arity: []int{2}},
		"CURRBARSCOUNT":  {validArgs: nil, arity: []int{0}},
		"DEVSQ":          {validArgs: seriesPeriod, arity: []int{2}},
		"DMA":            {validArgs: []string{"C", "0.5"}, arity: []int{2}},
		"DOWNNDAY":       {validArgs: seriesPeriod, arity: []int{2}},
		"DRAWBAND":       {validArgs: []string{"H", "1", "L", "2"}, arity: []int{4}},
		"DRAWICON":       {validArgs: []string{"C > O", "H", "1"}, arity: []int{3}},
		"DRAWKLINE":      {validArgs: []string{"H", "O", "L", "C"}, arity: []int{4}},
		"DRAWLINE":       {validArgs: []string{"C > O", "L", "O > C", "H", "0"}, arity: []int{5}},
		"DRAWNULL":       {validArgs: nil, arity: []int{0}},
		"DRAWNUMBER":     {validArgs: []string{"C > O", "C", "C"}, arity: []int{3}},
		"DRAWTEXT":       {validArgs: []string{"C > O", "L", "'UP'"}, arity: []int{3}},
		"EMA":            {validArgs: seriesPeriod, arity: []int{2}},
		"EVERY":          {validArgs: conditionPeriod, arity: []int{2}},
		"EXIST":          {validArgs: conditionPeriod, arity: []int{2}},
		"EXISTR":         {validArgs: conditionWindow, arity: []int{3}},
		"EXP":            {validArgs: unaryNumeric, arity: []int{1}},
		"FILTER":         {validArgs: conditionPeriod, arity: []int{2}},
		"FLOOR":          {validArgs: unaryNumeric, arity: []int{1}},
		"FORCAST":        {validArgs: seriesPeriod, arity: []int{2}},
		"FRACPART":       {validArgs: unaryNumeric, arity: []int{1}},
		"HHV":            {validArgs: seriesPeriod, arity: []int{2}},
		"HHVBARS":        {validArgs: seriesPeriod, arity: []int{2}},
		"IF":             {validArgs: []string{"C > O", "H", "L"}, arity: []int{3}},
		"IFF":            {validArgs: []string{"C > O", "H", "L"}, arity: []int{3}},
		"IFN":            {validArgs: []string{"C > O", "H", "L"}, arity: []int{3}},
		"INTPART":        {validArgs: unaryNumeric, arity: []int{1}},
		"ISLASTBAR":      {validArgs: nil, arity: []int{0}},
		"LAST":           {validArgs: conditionWindow, arity: []int{3}},
		"LLV":            {validArgs: seriesPeriod, arity: []int{2}},
		"LLVBARS":        {validArgs: seriesPeriod, arity: []int{2}},
		"LN":             {validArgs: unaryNumeric, arity: []int{1}},
		"LOG":            {validArgs: unaryNumeric, arity: []int{1}},
		"LONGCROSS":      {validArgs: pairPeriod, arity: []int{3}},
		"MA":             {validArgs: seriesPeriod, arity: []int{2}},
		"MAX":            {validArgs: binaryNumeric, arity: []int{2}},
		"MIN":            {validArgs: binaryNumeric, arity: []int{2}},
		"MOD":            {validArgs: binaryNumeric, arity: []int{2}},
		"NDAY":           {validArgs: pairPeriod, arity: []int{3}},
		"NOT":            {validArgs: unaryNumeric, arity: []int{1}},
		"POW":            {validArgs: binaryNumeric, arity: []int{2}},
		"POLYLINE":       {validArgs: []string{"C > O", "C"}, arity: []int{2}},
		"RANGE":          {validArgs: []string{"C", "100", "110"}, arity: []int{3}},
		"REF":            {validArgs: seriesPeriod, arity: []int{2}},
		"REFV":           {validArgs: seriesPeriod, arity: []int{2}},
		"REFX":           {validArgs: seriesPeriod, arity: []int{2}},
		"REFXV":          {validArgs: seriesPeriod, arity: []int{2}},
		"RELATE":         {validArgs: pairPeriod, arity: []int{3}},
		"ROUND":          {validArgs: unaryNumeric, arity: []int{1}},
		"ROUND2":         {validArgs: binaryNumeric, arity: []int{2}},
		"SIGN":           {validArgs: unaryNumeric, arity: []int{1}},
		"SIN":            {validArgs: unaryNumeric, arity: []int{1}},
		"SLOPE":          {validArgs: seriesPeriod, arity: []int{2}},
		"SMA":            {validArgs: []string{"C", "3", "1"}, arity: []int{2, 3}},
		"SQRT":           {validArgs: unaryNumeric, arity: []int{1}},
		"STD":            {validArgs: seriesPeriod, arity: []int{2}},
		"STDDEV":         {validArgs: seriesPeriod, arity: []int{2}},
		"STDP":           {validArgs: seriesPeriod, arity: []int{2}},
		"STICKLINE":      {validArgs: []string{"C > O", "O", "C", "2", "0"}, arity: []int{5}},
		"SUM":            {validArgs: seriesPeriod, arity: []int{2}},
		"SUMBARS":        {validArgs: []string{"VOL", "2500"}, arity: []int{2}},
		"TAN":            {validArgs: unaryNumeric, arity: []int{1}},
		"TOTALBARSCOUNT": {validArgs: nil, arity: []int{0}},
		"UPNDAY":         {validArgs: seriesPeriod, arity: []int{2}},
		"VALUEWHEN":      {validArgs: []string{"C > O", "C"}, arity: []int{2}},
		"VAR":            {validArgs: seriesPeriod, arity: []int{2}},
		"VARP":           {validArgs: seriesPeriod, arity: []int{2}},
		"WMA":            {validArgs: seriesPeriod, arity: []int{2}},
	}
}

func invalidArityCounts(valid []int) []int {
	validSet := make(map[int]bool, len(valid))
	maxValid := 0
	for _, count := range valid {
		validSet[count] = true
		if count > maxValid {
			maxValid = count
		}
	}

	counts := make([]int, 0, 2)
	if !validSet[0] {
		counts = append(counts, 0)
	}
	tooMany := maxValid + 1
	if !validSet[tooMany] {
		counts = append(counts, tooMany)
	}
	return counts
}

func argsForCount(validArgs []string, count int) []string {
	args := make([]string, count)
	for i := range args {
		if i < len(validArgs) {
			args[i] = validArgs[i]
		} else {
			args[i] = "1"
		}
	}
	return args
}

func argCountName(count int) string {
	if count == 0 {
		return "0"
	}
	return string(rune('0' + count))
}

func assertVariableIsNaN(t *testing.T, result interface {
	GetVariable(string) (float64, bool)
}, name string) {
	t.Helper()
	got, ok := result.GetVariable(name)
	if !ok {
		t.Fatalf("expected variable %s", name)
	}
	if !math.IsNaN(got) {
		t.Fatalf("expected variable %s to be NaN, got %f", name, got)
	}
}
