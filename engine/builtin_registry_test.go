package engine

import (
	"math"
	"os"
	"regexp"
	"sort"
	"strings"
	"testing"

	"github.com/DTrader-store/formula-go/interpreter"
	"github.com/DTrader-store/formula-go/types"
)

type builtinCoverageCase struct {
	formula string
	assert  func(t *testing.T, result *types.FormulaResult)
}

func TestRegisteredBuiltinsHaveCoverageCases(t *testing.T) {
	registryNames := interpreter.NewFunctionRegistry().Names()
	cases := builtinCoverageCases()

	for _, name := range registryNames {
		if _, ok := cases[name]; !ok {
			t.Errorf("registered builtin %s has no coverage case", name)
		}
	}

	for name := range cases {
		if !stringInSlice(name, registryNames) {
			t.Errorf("coverage case %s is not a registered builtin", name)
		}
	}

	engine := NewFormulaEngine()
	marketData := createTestData()
	for _, name := range registryNames {
		tc := cases[name]
		t.Run(name, func(t *testing.T) {
			result, err := engine.Run(tc.formula, marketData)
			if err != nil {
				t.Fatalf("Run() error: %v", err)
			}
			tc.assert(t, result)
		})
	}
}

func TestReadmeBuiltinListMatchesRegistry(t *testing.T) {
	content, err := os.ReadFile("../README.md")
	if err != nil {
		t.Fatalf("read README.md: %v", err)
	}

	readmeNames, err := readmeBuiltinNames(string(content))
	if err != nil {
		t.Fatal(err)
	}
	registryNames := interpreter.NewFunctionRegistry().Names()

	if diff := missingStrings(registryNames, readmeNames); len(diff) > 0 {
		t.Errorf("README missing registered builtins: %s", strings.Join(diff, ", "))
	}
	if diff := missingStrings(readmeNames, registryNames); len(diff) > 0 {
		t.Errorf("README lists unregistered builtins: %s", strings.Join(diff, ", "))
	}
}

func builtinCoverageCases() map[string]builtinCoverageCase {
	return map[string]builtinCoverageCase{
		"ABS":            variableCase("X: ABS(-5)", "X", 5),
		"ACOS":           variableCase("X: ACOS(1)", "X", 0),
		"ASIN":           variableCase("X: ASIN(0)", "X", 0),
		"ATAN":           variableCase("X: ATAN(0)", "X", 0),
		"AVEDEV":         outputIndexCase("X: AVEDEV(C, 3)", "X", 2, 1.3333333333),
		"BARSCOUNT":      outputSeriesCase("X: BARSCOUNT(C)", "X", []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
		"BARSLAST":       outputSeriesCase("UP := C > O\nX: BARSLAST(UP)", "X", []float64{0, 1, 0, 0, 1, 0, 1, 0, 0, 1}),
		"BARSLASTCOUNT":  outputSeriesCase("UP := C > O\nX: BARSLASTCOUNT(UP)", "X", []float64{1, 0, 1, 2, 0, 1, 0, 1, 2, 0}),
		"BARSSINCE":      outputSeriesCase("UP := C > O\nX: BARSSINCE(UP)", "X", []float64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}),
		"BARSTATUS":      outputSeriesCase("X: BARSTATUS()", "X", []float64{1, 2, 2, 2, 2, 2, 2, 2, 2, 3}),
		"BETA":           outputIndexCase("X: BETA(C, O, 3)", "X", 2, -0.3157894737),
		"BETWEEN":        outputSeriesCase("X: BETWEEN(C, 105, 112)", "X", []float64{1, 0, 1, 1, 1, 1, 1, 1, 0, 0}),
		"CEILING":        variableCase("X: CEILING(1.2)", "X", 2),
		"CONST":          outputSeriesCase("X: CONST(C)", "X", []float64{113, 113, 113, 113, 113, 113, 113, 113, 113, 113}),
		"COS":            variableCase("X: COS(0)", "X", 1),
		"COUNT":          outputSeriesCase("UP := C > O\nX: COUNT(UP, 3)", "X", []float64{math.NaN(), math.NaN(), 2, 2, 2, 2, 1, 2, 2, 2}),
		"COVAR":          outputIndexCase("X: COVAR(C, O, 3)", "X", 2, -1.3333333333),
		"CROSS":          outputSeriesCase("A := REF(C, 1)\nX: CROSS(C, A)", "X", []float64{0, 0, 1, 0, 0, 1, 0, 1, 0, 0}),
		"CURRBARSCOUNT":  outputSeriesCase("X: CURRBARSCOUNT()", "X", []float64{10, 9, 8, 7, 6, 5, 4, 3, 2, 1}),
		"DEVSQ":          outputIndexCase("X: DEVSQ(C, 3)", "X", 2, 8),
		"DMA":            outputSeriesCase("X: DMA(C, 0.5)", "X", []float64{105, 104, 105.5, 107.75, 107.875, 109.4375, 109.21875, 110.609375, 112.8046875, 112.90234375}),
		"DOWNNDAY":       outputSeriesCase("X: DOWNNDAY(C, 2)", "X", []float64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}),
		"DRAWBAND":       drawingCase("X := DRAWBAND(H, 1, L, 2)", "DRAWBAND", 10),
		"DRAWICON":       drawingCase("X := DRAWICON(C > O, H, 1)", "DRAWICON", 6),
		"DRAWKLINE":      drawingCase("X := DRAWKLINE(H, O, L, C)", "DRAWKLINE", 10),
		"DRAWLINE":       drawingCase("X := DRAWLINE(BARSTATUS() = 1, L, ISLASTBAR(), H, 0)", "DRAWLINE", 1),
		"DRAWNULL":       outputSeriesCase("X: IF(C > O, C, DRAWNULL())", "X", []float64{105, math.NaN(), 107, 110, math.NaN(), 111, math.NaN(), 112, 115, math.NaN()}),
		"DRAWNUMBER":     drawingCase("X := DRAWNUMBER(C > O, C, C)", "DRAWNUMBER", 6),
		"DRAWTEXT":       drawingCase("X := DRAWTEXT(C > O, L, 'UP')", "DRAWTEXT", 6),
		"EMA":            outputIndexCase("X: EMA(C, 5)", "X", 0, 105),
		"EVERY":          outputSeriesCase("UP := C > O\nX: EVERY(UP, 2)", "X", []float64{0, 0, 0, 1, 0, 0, 0, 0, 1, 0}),
		"EXIST":          outputSeriesCase("UP := C > O\nX: EXIST(UP, 2)", "X", []float64{0, 1, 1, 1, 1, 1, 1, 1, 1, 1}),
		"EXISTR":         outputSeriesCase("UP := C > O\nX: EXISTR(UP, 2, 1)", "X", []float64{0, 0, 1, 1, 1, 1, 1, 1, 1, 1}),
		"EXP":            variableCase("X: EXP(0)", "X", 1),
		"FILTER":         outputSeriesCase("UP := C > O\nX: FILTER(UP, 3)", "X", []float64{1, 0, 0, 1, 0, 0, 0, 1, 0, 0}),
		"FLOOR":          variableCase("X: FLOOR(1.8)", "X", 1),
		"FORCAST":        outputIndexCase("X: FORCAST(C, 3)", "X", 2, 106),
		"FRACPART":       variableCase("X: FRACPART(1.25)", "X", 0.25),
		"HHV":            outputIndexCase("X: HHV(H, 5)", "X", 4, 113),
		"HHVBARS":        outputSeriesCase("X: HHVBARS(H, 5)", "X", []float64{math.NaN(), math.NaN(), math.NaN(), math.NaN(), 0, 0, 0, 0, 0, 0}),
		"IF":             outputSeriesCase("X: IF(C > O, H, L)", "X", []float64{107, 102, 109, 112, 107, 114, 108, 116, 117, 112}),
		"IFF":            outputSeriesCase("X: IFF(C > O, H, L)", "X", []float64{107, 102, 109, 112, 107, 114, 108, 116, 117, 112}),
		"IFN":            outputSeriesCase("X: IFN(C > O, H, L)", "X", []float64{99, 108, 101, 106, 113, 107, 115, 108, 110, 118}),
		"INTPART":        variableCase("X: INTPART(-1.8)", "X", -1),
		"ISLASTBAR":      outputSeriesCase("X: ISLASTBAR()", "X", []float64{0, 0, 0, 0, 0, 0, 0, 0, 0, 1}),
		"LAST":           outputSeriesCase("UP := C > O\nX: LAST(UP, 2, 1)", "X", []float64{0, 0, 0, 0, 1, 0, 0, 0, 0, 1}),
		"LLV":            outputIndexCase("X: LLV(L, 5)", "X", 4, 99),
		"LLVBARS":        outputSeriesCase("X: LLVBARS(L, 5)", "X", []float64{math.NaN(), math.NaN(), math.NaN(), math.NaN(), 4, 3, 4, 4, 3, 4}),
		"LN":             variableCase("X: LN(EXP(1))", "X", 1),
		"LOG":            variableCase("X: LOG(100)", "X", 2),
		"LONGCROSS":      outputSeriesCase("X: LONGCROSS(C, O, 1)", "X", []float64{0, 0, 1, 0, 0, 1, 0, 1, 0, 0}),
		"MA":             outputIndexCase("X: MA(C, 5)", "X", 4, 106.6),
		"MAX":            variableCase("X: MAX(2, 3)", "X", 3),
		"MIN":            variableCase("X: MIN(2, 3)", "X", 2),
		"MOD":            variableCase("X: MOD(10, 3)", "X", 1),
		"NDAY":           outputSeriesCase("X: NDAY(C, O, 2)", "X", []float64{0, 0, 0, 1, 0, 0, 0, 0, 1, 0}),
		"NOT":            outputSeriesCase("X: NOT(C > O)", "X", []float64{0, 1, 0, 0, 1, 0, 1, 0, 0, 1}),
		"POW":            variableCase("X: POW(2, 3)", "X", 8),
		"POLYLINE":       drawingCase("X := POLYLINE(C > O, C)", "POLYLINE", 6),
		"RANGE":          outputSeriesCase("X: RANGE(C, 103, 110)", "X", []float64{1, 1, 1, 1, 1, 0, 1, 0, 0, 0}),
		"REF":            outputSeriesCase("X: REF(C, 1)", "X", []float64{math.NaN(), 105, 103, 107, 110, 108, 111, 109, 112, 115}),
		"REFV":           outputSeriesCase("X: REFV(C, 1)", "X", []float64{math.NaN(), 105, 103, 107, 110, 108, 111, 109, 112, 115}),
		"REFX":           outputSeriesCase("X: REFX(C, 1)", "X", []float64{103, 107, 110, 108, 111, 109, 112, 115, 113, math.NaN()}),
		"REFXV":          outputSeriesCase("X: REFXV(C, 2)", "X", []float64{107, 110, 108, 111, 109, 112, 115, 113, math.NaN(), math.NaN()}),
		"RELATE":         outputIndexCase("X: RELATE(C, O, 3)", "X", 2, -0.3973597071),
		"ROUND":          variableCase("X: ROUND(1.5)", "X", 2),
		"ROUND2":         variableCase("X: ROUND2(1.234, 2)", "X", 1.23),
		"SIGN":           variableCase("X: SIGN(-5)", "X", -1),
		"SIN":            variableCase("X: SIN(0)", "X", 0),
		"SLOPE":          outputIndexCase("X: SLOPE(C, 3)", "X", 2, 1),
		"SMA":            outputIndexCase("X: SMA(C, 3, 1)", "X", 1, 104.3333333333),
		"SQRT":           variableCase("X: SQRT(9)", "X", 3),
		"STD":            outputIndexCase("X: STD(C, 3)", "X", 2, 1.6329931619),
		"STDDEV":         outputIndexCase("X: STDDEV(C, 3)", "X", 2, 2),
		"STDP":           outputIndexCase("X: STDP(C, 3)", "X", 2, 1.6329931619),
		"STICKLINE":      drawingCase("X := STICKLINE(C > O, O, C, 2, 0)", "STICKLINE", 6),
		"SUM":            outputIndexCase("X: SUM(C, 3)", "X", 2, 315),
		"SUMBARS":        outputSeriesCase("X: SUMBARS(VOL, 2500)", "X", []float64{math.NaN(), math.NaN(), 3, 2, 2, 2, 2, 2, 2, 2}),
		"TAN":            variableCase("X: TAN(0)", "X", 0),
		"TOTALBARSCOUNT": outputSeriesCase("X: TOTALBARSCOUNT()", "X", []float64{10, 10, 10, 10, 10, 10, 10, 10, 10, 10}),
		"UPNDAY":         outputSeriesCase("X: UPNDAY(C, 2)", "X", []float64{0, 0, 0, 1, 0, 0, 0, 0, 1, 0}),
		"VALUEWHEN":      outputSeriesCase("X: VALUEWHEN(C > O, C)", "X", []float64{105, 105, 107, 110, 110, 111, 111, 112, 115, 115}),
		"VAR":            outputIndexCase("X: VAR(C, 3)", "X", 2, 2.6666666667),
		"VARP":           outputIndexCase("X: VARP(C, 3)", "X", 2, 2.6666666667),
		"WMA":            outputIndexCase("X: WMA(C, 3)", "X", 2, 105.3333333333),
	}
}

func variableCase(formula, name string, want float64) builtinCoverageCase {
	return builtinCoverageCase{
		formula: formula,
		assert: func(t *testing.T, result *types.FormulaResult) {
			t.Helper()
			assertVariable(t, result, name, want)
		},
	}
}

func outputIndexCase(formula, name string, index int, want float64) builtinCoverageCase {
	return builtinCoverageCase{
		formula: formula,
		assert: func(t *testing.T, result *types.FormulaResult) {
			t.Helper()
			got := outputByName(t, result, name)
			if index >= len(got) {
				t.Fatalf("output %s has %d values, cannot check index %d", name, len(got), index)
			}
			assertClose(t, got[index], want)
		},
	}
}

func outputSeriesCase(formula, name string, want []float64) builtinCoverageCase {
	return builtinCoverageCase{
		formula: formula,
		assert: func(t *testing.T, result *types.FormulaResult) {
			t.Helper()
			assertSeries(t, outputByName(t, result, name), want)
		},
	}
}

func drawingCase(formula, function string, wantCount int) builtinCoverageCase {
	return builtinCoverageCase{
		formula: formula,
		assert: func(t *testing.T, result *types.FormulaResult) {
			t.Helper()
			if len(result.Drawings) != wantCount {
				t.Fatalf("expected %d drawing events, got %d", wantCount, len(result.Drawings))
			}
			if result.Drawings[0].Function != function {
				t.Fatalf("expected first drawing function %s, got %s", function, result.Drawings[0].Function)
			}
		},
	}
}

func outputByName(t *testing.T, result *types.FormulaResult, name string) []float64 {
	t.Helper()
	output, ok := outputsByName(result)[name]
	if !ok {
		t.Fatalf("expected output %s", name)
	}
	return output
}

func readmeBuiltinNames(content string) ([]string, error) {
	startMarker := "### 2. 内置函数"
	endMarker := "### 3. 内置变量"
	start := strings.Index(content, startMarker)
	if start < 0 {
		return nil, os.ErrNotExist
	}
	end := strings.Index(content[start:], endMarker)
	if end < 0 {
		return nil, os.ErrNotExist
	}

	section := content[start : start+end]
	signaturePattern := regexp.MustCompile("`([A-Z][A-Z0-9]*)\\(")
	matches := signaturePattern.FindAllStringSubmatch(section, -1)

	seen := make(map[string]bool)
	for _, match := range matches {
		seen[match[1]] = true
	}

	names := make([]string, 0, len(seen))
	for name := range seen {
		names = append(names, name)
	}
	sort.Strings(names)
	return names, nil
}

func missingStrings(want, got []string) []string {
	gotSet := make(map[string]bool, len(got))
	for _, value := range got {
		gotSet[value] = true
	}

	var missing []string
	for _, value := range want {
		if !gotSet[value] {
			missing = append(missing, value)
		}
	}
	return missing
}

func stringInSlice(value string, values []string) bool {
	for _, item := range values {
		if item == value {
			return true
		}
	}
	return false
}
