package engine

import (
	"fmt"
	"testing"

	"github.com/DTrader-store/formula-go/parser/ast"
	"github.com/DTrader-store/formula-go/types"
)

var benchmarkResult *types.FormulaResult
var benchmarkProgram *ast.Program

const benchmarkMACDFormula = `
	EMA12 := EMA(CLOSE, 12)
	EMA26 := EMA(CLOSE, 26)
	DIF := EMA12 - EMA26
	DEA := EMA(DIF, 9)
	MACD := (DIF - DEA) * 2
`

const benchmarkRollingFormula = `
	MA20 := MA(CLOSE, 20)
	SUM20 := SUM(CLOSE, 20)
	STD20 := STD(CLOSE, 20)
`

const benchmarkExtendedRollingFormula = `
	DEV20 := DEVSQ(CLOSE, 20)
	SAMPLE_STD20 := STDDEV(CLOSE, 20)
	COV20 := COVAR(CLOSE, OPEN, 20)
	REL20 := RELATE(CLOSE, OPEN, 20)
	BETA20 := BETA(CLOSE, OPEN, 20)
`

const benchmarkDrawingFormula = `
	UP := CLOSE > OPEN
	TEXT_MARK := DRAWTEXT(UP, LOW, 'UP')
	ICON_MARK := DRAWICON(UP, HIGH, 1)
	NUMBER_MARK := DRAWNUMBER(UP, CLOSE, CLOSE)
	STICK_MARK := STICKLINE(UP, OPEN, CLOSE, 2, 0)
`

func BenchmarkCompileMACD(b *testing.B) {
	engine := NewFormulaEngine()
	var program *ast.Program
	for i := 0; i < b.N; i++ {
		var err error
		program, err = engine.Compile(benchmarkMACDFormula)
		if err != nil {
			b.Fatal(err)
		}
	}
	benchmarkProgram = program
}

func BenchmarkRunMACD(b *testing.B) {
	for _, size := range []int{1000, 10000} {
		b.Run(fmt.Sprintf("%dBars", size), func(b *testing.B) {
			benchmarkRunFormula(b, benchmarkMACDFormula, createBenchmarkMarketData(size))
		})
	}
}

func BenchmarkExecuteCompiledMACD(b *testing.B) {
	for _, size := range []int{1000, 10000} {
		b.Run(fmt.Sprintf("%dBars", size), func(b *testing.B) {
			benchmarkExecuteCompiledFormula(b, benchmarkMACDFormula, createBenchmarkMarketData(size))
		})
	}
}

func BenchmarkRunRollingFunctions(b *testing.B) {
	for _, size := range []int{1000, 10000} {
		b.Run(fmt.Sprintf("%dBars", size), func(b *testing.B) {
			benchmarkRunFormula(b, benchmarkRollingFormula, createBenchmarkMarketData(size))
		})
	}
}

func BenchmarkRunExtendedRollingFunctions(b *testing.B) {
	for _, size := range []int{1000, 10000} {
		b.Run(fmt.Sprintf("%dBars", size), func(b *testing.B) {
			benchmarkRunFormula(b, benchmarkExtendedRollingFormula, createBenchmarkMarketData(size))
		})
	}
}

func BenchmarkRunDrawingEvents(b *testing.B) {
	for _, size := range []int{1000, 10000} {
		b.Run(fmt.Sprintf("%dBars", size), func(b *testing.B) {
			benchmarkRunFormula(b, benchmarkDrawingFormula, createBenchmarkMarketData(size))
		})
	}
}

func benchmarkRunFormula(b *testing.B, formula string, marketData []*types.MarketData) {
	b.Helper()
	engine := NewFormulaEngine()
	b.ReportAllocs()
	b.ResetTimer()

	var result *types.FormulaResult
	for i := 0; i < b.N; i++ {
		var err error
		result, err = engine.Run(formula, marketData)
		if err != nil {
			b.Fatal(err)
		}
	}
	benchmarkResult = result
}

func benchmarkExecuteCompiledFormula(b *testing.B, formula string, marketData []*types.MarketData) {
	b.Helper()
	engine := NewFormulaEngine()
	program, err := engine.Compile(formula)
	if err != nil {
		b.Fatal(err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	var result *types.FormulaResult
	for i := 0; i < b.N; i++ {
		result, err = engine.Execute(program, marketData)
		if err != nil {
			b.Fatal(err)
		}
	}
	benchmarkResult = result
}

func createBenchmarkMarketData(count int) []*types.MarketData {
	data := make([]*types.MarketData, count)
	for i := range data {
		open := 100.0 + float64(i%200)*0.3 + float64(i/200)*0.05
		close := open + float64((i%7)-3)*0.4
		high := maxFloat(open, close) + 1.5 + float64(i%5)*0.1
		low := minFloat(open, close) - 1.5 - float64(i%3)*0.1
		volume := 1000.0 + float64((i*37)%1000)
		amount := close * volume
		data[i] = types.NewMarketData(open, close, high, low, volume, amount)
	}
	return data
}

func maxFloat(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func minFloat(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
