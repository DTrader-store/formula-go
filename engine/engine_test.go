package engine

import (
	"math"
	"testing"

	"github.com/DTrader-store/formula-go/types"
)

// createTestData creates sample market data for testing
func createTestData() []*types.MarketData {
	return []*types.MarketData{
		types.NewMarketData(100, 105, 107, 99, 1000, 100000),
		types.NewMarketData(105, 103, 108, 102, 1100, 110000),
		types.NewMarketData(103, 107, 109, 101, 1200, 120000),
		types.NewMarketData(107, 110, 112, 106, 1300, 130000),
		types.NewMarketData(110, 108, 113, 107, 1400, 140000),
		types.NewMarketData(108, 111, 114, 107, 1500, 150000),
		types.NewMarketData(111, 109, 115, 108, 1600, 160000),
		types.NewMarketData(109, 112, 116, 108, 1700, 170000),
		types.NewMarketData(112, 115, 117, 110, 1800, 180000),
		types.NewMarketData(115, 113, 118, 112, 1900, 190000),
	}
}

func TestEngineSimpleExpression(t *testing.T) {
	engine := NewFormulaEngine()
	marketData := createTestData()

	formula := "CLOSE"
	result, err := engine.Run(formula, marketData)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	if len(result.Outputs) != 1 {
		t.Fatalf("Expected 1 output, got %d", len(result.Outputs))
	}

	if len(result.Outputs[0].Data) != len(marketData) {
		t.Errorf("Expected %d data points, got %d", len(marketData), len(result.Outputs[0].Data))
	}
}

func TestEngineMA(t *testing.T) {
	engine := NewFormulaEngine()
	marketData := createTestData()

	formula := "MA5 := MA(CLOSE, 5)"
	result, err := engine.Run(formula, marketData)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	if len(result.Outputs) != 1 {
		t.Fatalf("Expected 1 output, got %d", len(result.Outputs))
	}

	ma5 := result.Outputs[0]
	if ma5.Name != "MA5" {
		t.Errorf("Expected output name 'MA5', got '%s'", ma5.Name)
	}

	// First 4 values should be NaN
	for i := 0; i < 4; i++ {
		if !math.IsNaN(ma5.Data[i]) {
			t.Errorf("Expected NaN at index %d, got %f", i, ma5.Data[i])
		}
	}

	// 5th value should be average of first 5 closes
	expected := (105.0 + 103.0 + 107.0 + 110.0 + 108.0) / 5.0
	if math.Abs(ma5.Data[4]-expected) > 0.01 {
		t.Errorf("Expected MA5[4] = %.2f, got %.2f", expected, ma5.Data[4])
	}
}

func TestEngineArithmetic(t *testing.T) {
	engine := NewFormulaEngine()
	marketData := createTestData()

	formula := "DIFF := HIGH - LOW"
	result, err := engine.Run(formula, marketData)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	if len(result.Outputs) != 1 {
		t.Fatalf("Expected 1 output, got %d", len(result.Outputs))
	}

	diff := result.Outputs[0]
	for i := range marketData {
		expected := marketData[i].High - marketData[i].Low
		if math.Abs(diff.Data[i]-expected) > 0.01 {
			t.Errorf("Index %d: expected %.2f, got %.2f", i, expected, diff.Data[i])
		}
	}
}

func TestEngineComparison(t *testing.T) {
	engine := NewFormulaEngine()
	marketData := createTestData()

	formula := "SIGNAL := CLOSE > OPEN"
	result, err := engine.Run(formula, marketData)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	if len(result.Outputs) != 1 {
		t.Fatalf("Expected 1 output, got %d", len(result.Outputs))
	}

	signal := result.Outputs[0]
	for i := range marketData {
		expected := 0.0
		if marketData[i].Close > marketData[i].Open {
			expected = 1.0
		}
		if signal.Data[i] != expected {
			t.Errorf("Index %d: expected %.0f, got %.0f", i, expected, signal.Data[i])
		}
	}
}

func TestEngineMultipleStatements(t *testing.T) {
	engine := NewFormulaEngine()
	marketData := createTestData()

	formula := `
		MA5 := MA(CLOSE, 5)
		MA10 := MA(CLOSE, 10)
		CROSS_SIGNAL := CROSS(MA5, MA10)
	`

	result, err := engine.Run(formula, marketData)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	if len(result.Outputs) != 3 {
		t.Fatalf("Expected 3 outputs, got %d", len(result.Outputs))
	}

	names := []string{"MA5", "MA10", "CROSS_SIGNAL"}
	for i, output := range result.Outputs {
		if output.Name != names[i] {
			t.Errorf("Expected output name '%s', got '%s'", names[i], output.Name)
		}
	}
}

func TestEngineREF(t *testing.T) {
	engine := NewFormulaEngine()
	marketData := createTestData()

	formula := "PREV_CLOSE := REF(CLOSE, 1)"
	result, err := engine.Run(formula, marketData)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	if len(result.Outputs) != 1 {
		t.Fatalf("Expected 1 output, got %d", len(result.Outputs))
	}

	prevClose := result.Outputs[0]

	// First value should be NaN
	if !math.IsNaN(prevClose.Data[0]) {
		t.Errorf("Expected NaN at index 0, got %f", prevClose.Data[0])
	}

	// Check other values
	for i := 1; i < len(marketData); i++ {
		expected := marketData[i-1].Close
		if math.Abs(prevClose.Data[i]-expected) > 0.01 {
			t.Errorf("Index %d: expected %.2f, got %.2f", i, expected, prevClose.Data[i])
		}
	}
}

func TestEngineIF(t *testing.T) {
	engine := NewFormulaEngine()
	marketData := createTestData()

	formula := "RESULT := IF(CLOSE > OPEN, HIGH, LOW)"
	result, err := engine.Run(formula, marketData)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	if len(result.Outputs) != 1 {
		t.Fatalf("Expected 1 output, got %d", len(result.Outputs))
	}

	output := result.Outputs[0]
	for i := range marketData {
		expected := marketData[i].Low
		if marketData[i].Close > marketData[i].Open {
			expected = marketData[i].High
		}
		if math.Abs(output.Data[i]-expected) > 0.01 {
			t.Errorf("Index %d: expected %.2f, got %.2f", i, expected, output.Data[i])
		}
	}
}

func TestEngineHHVLLV(t *testing.T) {
	engine := NewFormulaEngine()
	marketData := createTestData()

	formula := `
		HIGHEST := HHV(HIGH, 5)
		LOWEST := LLV(LOW, 5)
	`

	result, err := engine.Run(formula, marketData)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	if len(result.Outputs) != 2 {
		t.Fatalf("Expected 2 outputs, got %d", len(result.Outputs))
	}

	// Verify HHV
	highest := result.Outputs[0]
	for i := 4; i < len(marketData); i++ {
		maxVal := marketData[i].High
		for j := 1; j < 5; j++ {
			if marketData[i-j].High > maxVal {
				maxVal = marketData[i-j].High
			}
		}
		if math.Abs(highest.Data[i]-maxVal) > 0.01 {
			t.Errorf("HHV at index %d: expected %.2f, got %.2f", i, maxVal, highest.Data[i])
		}
	}
}

func TestEngineEMA(t *testing.T) {
	engine := NewFormulaEngine()
	marketData := createTestData()

	formula := "EMA5 := EMA(CLOSE, 5)"
	result, err := engine.Run(formula, marketData)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	if len(result.Outputs) != 1 {
		t.Fatalf("Expected 1 output, got %d", len(result.Outputs))
	}

	ema5 := result.Outputs[0]
	if ema5.Name != "EMA5" {
		t.Errorf("Expected output name 'EMA5', got '%s'", ema5.Name)
	}

	// First value should equal first close
	if math.Abs(ema5.Data[0]-marketData[0].Close) > 0.01 {
		t.Errorf("Expected EMA5[0] = %.2f, got %.2f", marketData[0].Close, ema5.Data[0])
	}
}

func TestEngineComplexFormula(t *testing.T) {
	engine := NewFormulaEngine()
	marketData := createTestData()

	// MACD-like formula
	formula := `
		EMA12 := EMA(CLOSE, 5)
		EMA26 := EMA(CLOSE, 8)
		DIF := EMA12 - EMA26
		DEA := EMA(DIF, 3)
		MACD := (DIF - DEA) * 2
	`

	result, err := engine.Run(formula, marketData)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	if len(result.Outputs) != 5 {
		t.Fatalf("Expected 5 outputs, got %d", len(result.Outputs))
	}

	expectedNames := []string{"EMA12", "EMA26", "DIF", "DEA", "MACD"}
	for i, name := range expectedNames {
		if result.Outputs[i].Name != name {
			t.Errorf("Output %d: expected name '%s', got '%s'", i, name, result.Outputs[i].Name)
		}
	}
}

func TestEngineTDXSMA(t *testing.T) {
	engine := NewFormulaEngine()
	marketData := createTestData()

	result, err := engine.Run("SMA3 := SMA(CLOSE, 3, 1)", marketData)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	got := result.Outputs[0].Data
	expected := []float64{105, 104.3333333333, 105.2222222222, 106.8148148148, 107.2098765432}
	for i, want := range expected {
		if math.Abs(got[i]-want) > 0.0001 {
			t.Errorf("SMA3[%d]: expected %.6f, got %.6f", i, want, got[i])
		}
	}
}

func TestEngineBarsFunctions(t *testing.T) {
	engine := NewFormulaEngine()
	marketData := createTestData()

	formula := `
		UP := CLOSE > OPEN
		BAR_COUNT := BARSCOUNT(CLOSE)
		SINCE := BARSSINCE(UP)
		LAST_COUNT := BARSLASTCOUNT(UP)
	`
	result, err := engine.Run(formula, marketData)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	outputs := outputsByName(result)
	assertSeries(t, outputs["BAR_COUNT"], []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	assertSeries(t, outputs["SINCE"], []float64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9})
	assertSeries(t, outputs["LAST_COUNT"], []float64{1, 0, 1, 2, 0, 1, 0, 1, 2, 0})
}

func TestEngineHHVBarsLLVBars(t *testing.T) {
	engine := NewFormulaEngine()
	marketData := createTestData()

	formula := `
		HH := HHVBARS(HIGH, 5)
		LL := LLVBARS(LOW, 5)
	`
	result, err := engine.Run(formula, marketData)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	outputs := outputsByName(result)
	assertSeries(t, outputs["HH"], []float64{math.NaN(), math.NaN(), math.NaN(), math.NaN(), 0, 0, 0, 0, 0, 0})
	assertSeries(t, outputs["LL"], []float64{math.NaN(), math.NaN(), math.NaN(), math.NaN(), 4, 3, 4, 4, 3, 4})
}

func TestEngineDMAConstValueWhen(t *testing.T) {
	engine := NewFormulaEngine()
	marketData := createTestData()

	formula := `
		DYNAMIC := DMA(CLOSE, 0.5)
		FINAL_CLOSE := CONST(CLOSE)
		LAST_UP_CLOSE := VALUEWHEN(CLOSE > OPEN, CLOSE)
	`
	result, err := engine.Run(formula, marketData)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	outputs := outputsByName(result)
	assertSeries(t, outputs["DYNAMIC"], []float64{105, 104, 105.5, 107.75, 107.875, 109.4375, 109.21875, 110.609375, 112.8046875, 112.90234375})
	assertSeries(t, outputs["FINAL_CLOSE"], []float64{113, 113, 113, 113, 113, 113, 113, 113, 113, 113})
	assertSeries(t, outputs["LAST_UP_CLOSE"], []float64{105, 105, 107, 110, 110, 111, 111, 112, 115, 115})
}

func TestEngineTDXOutputDeclarationStyles(t *testing.T) {
	engine := NewFormulaEngine()
	marketData := createTestData()

	formula := `
		DIF: EMA(CLOSE, 3), COLORWHITE, LINETHICK2
		DEA: MA(DIF, 3), COLORSTICK, NODRAW
	`
	result, err := engine.Run(formula, marketData)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	if len(result.Outputs) != 2 {
		t.Fatalf("Expected 2 outputs, got %d", len(result.Outputs))
	}

	dif := result.Outputs[0]
	if dif.Name != "DIF" {
		t.Fatalf("Expected first output DIF, got %s", dif.Name)
	}
	if dif.Style == nil {
		t.Fatal("Expected DIF style")
	}
	if dif.Style.Color != "WHITE" {
		t.Errorf("Expected WHITE color, got %s", dif.Style.Color)
	}
	if dif.Style.LineWidth != 2 {
		t.Errorf("Expected line width 2, got %d", dif.Style.LineWidth)
	}

	dea := result.Outputs[1]
	if dea.Name != "DEA" {
		t.Fatalf("Expected second output DEA, got %s", dea.Name)
	}
	if dea.Style == nil {
		t.Fatal("Expected DEA style")
	}
	if dea.Style.DrawMethod != "colorstick" {
		t.Errorf("Expected colorstick draw method, got %s", dea.Style.DrawMethod)
	}
	if !dea.Style.Hidden {
		t.Error("Expected DEA to be hidden")
	}
}

func TestEngineTDXAliases(t *testing.T) {
	engine := NewFormulaEngine()
	marketData := createTestData()

	formula := `
		SHORT_CLOSE: C
		RANGE_VALUE: H - L
		TURNOVER: VOL + V + AMO * 0
		OPEN_ALIAS: O
	`
	result, err := engine.Run(formula, marketData)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	outputs := outputsByName(result)
	assertSeries(t, outputs["SHORT_CLOSE"], []float64{105, 103, 107, 110, 108, 111, 109, 112, 115, 113})
	assertSeries(t, outputs["RANGE_VALUE"], []float64{8, 6, 8, 6, 6, 7, 7, 8, 7, 6})
	assertSeries(t, outputs["TURNOVER"], []float64{2000, 2200, 2400, 2600, 2800, 3000, 3200, 3400, 3600, 3800})
	assertSeries(t, outputs["OPEN_ALIAS"], []float64{100, 105, 103, 107, 110, 108, 111, 109, 112, 115})
}

func TestEngineTDXLogicFunctions(t *testing.T) {
	engine := NewFormulaEngine()
	marketData := createTestData()

	formula := `
		UP := C > O
		NOT_UP: NOT(UP)
		UP2: UPNDAY(C, 2)
		DOWN2: DOWNNDAY(C, 2)
		ABOVE2: NDAY(C, O, 2)
		LAST_UP: LAST(UP, 2, 1)
		EXIST_RECENT: EXISTR(UP, 2, 1)
		IN_RANGE: RANGE(C, 103, 110)
		CROSS_LONG: LONGCROSS(C, O, 1)
		IF_ALIAS: IFF(UP, H, L)
		INVERSE_IF: IFN(UP, H, L)
		NULL_IF: IF(UP, C, DRAWNULL())
	`
	result, err := engine.Run(formula, marketData)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	outputs := outputsByName(result)
	assertSeries(t, outputs["NOT_UP"], []float64{0, 1, 0, 0, 1, 0, 1, 0, 0, 1})
	assertSeries(t, outputs["UP2"], []float64{0, 0, 0, 1, 0, 0, 0, 0, 1, 0})
	assertSeries(t, outputs["DOWN2"], []float64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
	assertSeries(t, outputs["ABOVE2"], []float64{0, 0, 0, 1, 0, 0, 0, 0, 1, 0})
	assertSeries(t, outputs["LAST_UP"], []float64{0, 0, 0, 0, 1, 0, 0, 0, 0, 1})
	assertSeries(t, outputs["EXIST_RECENT"], []float64{0, 0, 1, 1, 1, 1, 1, 1, 1, 1})
	assertSeries(t, outputs["IN_RANGE"], []float64{1, 1, 1, 1, 1, 0, 1, 0, 0, 0})
	assertSeries(t, outputs["CROSS_LONG"], []float64{0, 0, 1, 0, 0, 1, 0, 1, 0, 0})
	assertSeries(t, outputs["IF_ALIAS"], []float64{107, 102, 109, 112, 107, 114, 108, 116, 117, 112})
	assertSeries(t, outputs["INVERSE_IF"], []float64{99, 108, 101, 106, 113, 107, 115, 108, 110, 118})
	assertSeries(t, outputs["NULL_IF"], []float64{105, math.NaN(), 107, 110, math.NaN(), 111, math.NaN(), 112, 115, math.NaN()})
}

func TestEngineTDXDrawingFunctions(t *testing.T) {
	engine := NewFormulaEngine()
	marketData := createTestData()

	formula := `
		UP := C > O
		TEXT_MARK := DRAWTEXT(UP, L, 'UP')
		ICON_MARK := DRAWICON(UP, H, 1)
		NUMBER_MARK := DRAWNUMBER(UP, C, C)
		STICK_MARK := STICKLINE(UP, O, C, 2, 0)
	`
	result, err := engine.Run(formula, marketData)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	if len(result.Drawings) != 24 {
		t.Fatalf("Expected 24 drawing events, got %d", len(result.Drawings))
	}

	first := result.Drawings[0]
	if first.Function != "DRAWTEXT" || first.BarIndex != 0 {
		t.Fatalf("Expected first DRAWTEXT at bar 0, got %#v", first)
	}
	if first.Text != "UP" {
		t.Errorf("Expected text UP, got %s", first.Text)
	}
	if first.Values["price"] != 99 {
		t.Errorf("Expected text price 99, got %f", first.Values["price"])
	}

	icon := result.Drawings[6]
	if icon.Function != "DRAWICON" || icon.BarIndex != 0 || icon.Values["price"] != 107 || icon.Values["value"] != 1 {
		t.Fatalf("Unexpected DRAWICON event: %#v", icon)
	}

	number := result.Drawings[12]
	if number.Function != "DRAWNUMBER" || number.BarIndex != 0 || number.Values["price"] != 105 || number.Values["value"] != 105 {
		t.Fatalf("Unexpected DRAWNUMBER event: %#v", number)
	}

	firstStick := result.Drawings[18]
	if firstStick.Function != "STICKLINE" || firstStick.BarIndex != 0 {
		t.Fatalf("Expected first STICKLINE at bar 0, got %#v", firstStick)
	}
	if firstStick.Values["price1"] != 100 || firstStick.Values["price2"] != 105 || firstStick.Values["width"] != 2 || firstStick.Values["empty"] != 0 {
		t.Errorf("Unexpected first STICKLINE values: %#v", firstStick.Values)
	}

	stick := result.Drawings[len(result.Drawings)-1]
	if stick.Function != "STICKLINE" || stick.BarIndex != 8 {
		t.Fatalf("Expected last STICKLINE at bar 8, got %#v", stick)
	}
	if stick.Values["price1"] != 112 || stick.Values["price2"] != 115 || stick.Values["width"] != 2 || stick.Values["empty"] != 0 {
		t.Errorf("Unexpected STICKLINE values: %#v", stick.Values)
	}
}

func TestEngineTDXDrawingEventPayloads(t *testing.T) {
	engine := NewFormulaEngine()
	marketData := createTestData()

	formula := `
		LINE_MARK := DRAWLINE(BARSTATUS() = 1, L, ISLASTBAR(), H, 0)
		POLY_MARK := POLYLINE(C > O, C)
		BAND_MARK := DRAWBAND(H, 1, L, 2)
		KLINE_MARK := DRAWKLINE(H, O, L, C)
	`
	result, err := engine.Run(formula, marketData)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	if len(result.Drawings) != 27 {
		t.Fatalf("Expected 27 drawing events, got %d", len(result.Drawings))
	}

	line := result.Drawings[0]
	if line.Function != "DRAWLINE" || line.BarIndex != 0 {
		t.Fatalf("Unexpected DRAWLINE event: %#v", line)
	}
	if line.Values["startBar"] != 0 || line.Values["startPrice"] != 99 || line.Values["endBar"] != 9 || line.Values["endPrice"] != 118 || line.Values["expand"] != 0 {
		t.Errorf("Unexpected DRAWLINE values: %#v", line.Values)
	}

	poly := result.Drawings[1]
	if poly.Function != "POLYLINE" || poly.BarIndex != 0 || poly.Values["price"] != 105 {
		t.Fatalf("Unexpected first POLYLINE event: %#v", poly)
	}

	band := result.Drawings[7]
	if band.Function != "DRAWBAND" || band.BarIndex != 0 {
		t.Fatalf("Unexpected first DRAWBAND event: %#v", band)
	}
	if band.Values["upper"] != 107 || band.Values["upperColor"] != 1 || band.Values["lower"] != 99 || band.Values["lowerColor"] != 2 {
		t.Errorf("Unexpected DRAWBAND values: %#v", band.Values)
	}

	kline := result.Drawings[17]
	if kline.Function != "DRAWKLINE" || kline.BarIndex != 0 {
		t.Fatalf("Unexpected first DRAWKLINE event: %#v", kline)
	}
	if kline.Values["high"] != 107 || kline.Values["open"] != 100 || kline.Values["low"] != 99 || kline.Values["close"] != 105 {
		t.Errorf("Unexpected DRAWKLINE values: %#v", kline.Values)
	}
}

func TestEngineOutputDeclarationsKeepStandaloneDrawingEvents(t *testing.T) {
	engine := NewFormulaEngine()
	marketData := createTestData()

	formula := `
		PRICE:C, COLORBLACK;
		DRAWTEXT(C > O, C, 'B');
	`
	result, err := engine.Run(formula, marketData)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	if len(result.Outputs) != 1 {
		t.Fatalf("Expected 1 output, got %d", len(result.Outputs))
	}
	if len(result.Drawings) != 6 {
		t.Fatalf("Expected 6 drawing events, got %d", len(result.Drawings))
	}
	if result.Drawings[0].Function != "DRAWTEXT" || result.Drawings[0].Text != "B" {
		t.Fatalf("Unexpected first drawing event: %#v", result.Drawings[0])
	}
}

func TestEngineKeepsMultipleStandaloneDrawingEvents(t *testing.T) {
	engine := NewFormulaEngine()
	marketData := createTestData()

	formula := `
		DRAWTEXT(C > O, C, 'B');
		DRAWTEXT(O > C, C, 'S');
	`
	result, err := engine.Run(formula, marketData)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	if len(result.Drawings) != 10 {
		t.Fatalf("Expected 10 drawing events, got %d", len(result.Drawings))
	}
	if result.Drawings[0].Text != "B" {
		t.Fatalf("Expected first drawing text B, got %#v", result.Drawings[0])
	}
	if result.Drawings[6].Text != "S" {
		t.Fatalf("Expected seventh drawing text S, got %#v", result.Drawings[6])
	}
}

func TestEngineTDXMathTrigFunctions(t *testing.T) {
	engine := NewFormulaEngine()
	marketData := createTestData()

	formula := `
		MAX_VALUE: MAX(2, 3)
		MIN_VALUE: MIN(2, 3)
		ABS_VALUE: ABS(-5)
		SQRT_VALUE: SQRT(9)
		POW_VALUE: POW(2, 3)
		EXP_VALUE: EXP(0)
		LN_VALUE: LN(EXP(1))
		LOG_VALUE: LOG(100)
		MOD_VALUE: MOD(10, 3)
		CEILING_VALUE: CEILING(1.2)
		FLOOR_VALUE: FLOOR(1.8)
		INTPART_VALUE: INTPART(-1.8)
		FRACPART_VALUE: FRACPART(1.25)
		ROUND_VALUE: ROUND(1.5)
		ROUND2_VALUE: ROUND2(1.234, 2)
		SIGN_VALUE: SIGN(-5)
		SIN_VALUE: SIN(0)
		COS_VALUE: COS(0)
		TAN_VALUE: TAN(0)
		ASIN_VALUE: ASIN(0)
		ACOS_VALUE: ACOS(1)
		ATAN_VALUE: ATAN(0)
	`
	result, err := engine.Run(formula, marketData)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	assertVariable(t, result, "MAX_VALUE", 3)
	assertVariable(t, result, "MIN_VALUE", 2)
	assertVariable(t, result, "ABS_VALUE", 5)
	assertVariable(t, result, "SQRT_VALUE", 3)
	assertVariable(t, result, "POW_VALUE", 8)
	assertVariable(t, result, "EXP_VALUE", 1)
	assertVariable(t, result, "LN_VALUE", 1)
	assertVariable(t, result, "LOG_VALUE", 2)
	assertVariable(t, result, "MOD_VALUE", 1)
	assertVariable(t, result, "CEILING_VALUE", 2)
	assertVariable(t, result, "FLOOR_VALUE", 1)
	assertVariable(t, result, "INTPART_VALUE", -1)
	assertVariable(t, result, "FRACPART_VALUE", 0.25)
	assertVariable(t, result, "ROUND_VALUE", 2)
	assertVariable(t, result, "ROUND2_VALUE", 1.23)
	assertVariable(t, result, "SIGN_VALUE", -1)
	assertVariable(t, result, "SIN_VALUE", 0)
	assertVariable(t, result, "COS_VALUE", 1)
	assertVariable(t, result, "TAN_VALUE", 0)
	assertVariable(t, result, "ASIN_VALUE", 0)
	assertVariable(t, result, "ACOS_VALUE", 0)
	assertVariable(t, result, "ATAN_VALUE", 0)
}

func TestEngineTDXRollingFunctionPreciseSeries(t *testing.T) {
	engine := NewFormulaEngine()
	marketData := createTestData()

	formula := `
		UP := C > O
		SUM_VALUE: SUM(C, 3)
		STD_VALUE: STD(C, 3)
		VAR_VALUE: VAR(C, 3)
		AVEDEV_VALUE: AVEDEV(C, 3)
		WMA_VALUE: WMA(C, 3)
		COUNT_VALUE: COUNT(UP, 3)
		EVERY_VALUE: EVERY(UP, 2)
		EXIST_VALUE: EXIST(UP, 2)
		BARSLAST_VALUE: BARSLAST(UP)
		FILTER_VALUE: FILTER(UP, 3)
		BETWEEN_VALUE: BETWEEN(C, 105, 112)
	`
	result, err := engine.Run(formula, marketData)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	outputs := outputsByName(result)
	assertSeries(t, outputs["SUM_VALUE"], []float64{math.NaN(), math.NaN(), 315, 320, 325, 329, 328, 332, 336, 340})
	assertSeries(t, outputs["STD_VALUE"], []float64{math.NaN(), math.NaN(), 1.6329931619, 2.8674417557, 1.2472191289, 1.2472191289, 1.2472191289, 1.2472191289, 2.4494897428, 1.2472191289})
	assertSeries(t, outputs["VAR_VALUE"], []float64{math.NaN(), math.NaN(), 2.6666666667, 8.2222222222, 1.5555555556, 1.5555555556, 1.5555555556, 1.5555555556, 6, 1.5555555556})
	assertSeries(t, outputs["AVEDEV_VALUE"], []float64{math.NaN(), math.NaN(), 1.3333333333, 2.4444444444, 1.1111111111, 1.1111111111, 1.1111111111, 1.1111111111, 2, 1.1111111111})
	assertSeries(t, outputs["WMA_VALUE"], []float64{math.NaN(), math.NaN(), 105.3333333333, 107.8333333333, 108.5, 109.8333333333, 109.5, 110.8333333333, 113, 113.5})
	assertSeries(t, outputs["COUNT_VALUE"], []float64{math.NaN(), math.NaN(), 2, 2, 2, 2, 1, 2, 2, 2})
	assertSeries(t, outputs["EVERY_VALUE"], []float64{0, 0, 0, 1, 0, 0, 0, 0, 1, 0})
	assertSeries(t, outputs["EXIST_VALUE"], []float64{0, 1, 1, 1, 1, 1, 1, 1, 1, 1})
	assertSeries(t, outputs["BARSLAST_VALUE"], []float64{0, 1, 0, 0, 1, 0, 1, 0, 0, 1})
	assertSeries(t, outputs["FILTER_VALUE"], []float64{1, 0, 0, 1, 0, 0, 0, 1, 0, 0})
	assertSeries(t, outputs["BETWEEN_VALUE"], []float64{1, 0, 1, 1, 1, 1, 1, 1, 0, 0})
}

func TestEngineTDXStatRegressionFunctions(t *testing.T) {
	engine := NewFormulaEngine()
	marketData := createTestData()

	formula := `
		DEV: DEVSQ(C, 3)
		FORECAST_VALUE: FORCAST(C, 3)
		SLOPE_VALUE: SLOPE(C, 3)
		STDP_VALUE: STDP(C, 3)
		STDDEV_VALUE: STDDEV(C, 3)
		VARP_VALUE: VARP(C, 3)
		COV_VALUE: COVAR(C, O, 3)
		REL_VALUE: RELATE(C, O, 3)
		BETA_VALUE: BETA(C, O, 3)
	`
	result, err := engine.Run(formula, marketData)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	outputs := outputsByName(result)
	assertClose(t, outputs["DEV"][2], 8)
	assertClose(t, outputs["FORECAST_VALUE"][2], 106)
	assertClose(t, outputs["SLOPE_VALUE"][2], 1)
	assertClose(t, outputs["STDP_VALUE"][2], 1.6329931619)
	assertClose(t, outputs["STDDEV_VALUE"][2], 2)
	assertClose(t, outputs["VARP_VALUE"][2], 2.6666666667)
	assertClose(t, outputs["COV_VALUE"][2], -1.3333333333)
	assertClose(t, outputs["REL_VALUE"][2], -0.3973597071)
	assertClose(t, outputs["BETA_VALUE"][2], -0.3157894737)
}

func TestEngineTDXReferenceBarFunctions(t *testing.T) {
	engine := NewFormulaEngine()
	marketData := createTestData()

	formula := `
		PAST: REFV(C, 1)
		FUTURE: REFX(C, 1)
		FUTURE_V: REFXV(C, 2)
		CURR: CURRBARSCOUNT()
		TOTAL: TOTALBARSCOUNT()
		LAST_BAR: ISLASTBAR()
		STATUS: BARSTATUS()
		SUM_BARS: SUMBARS(VOL, 2500)
	`
	result, err := engine.Run(formula, marketData)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	outputs := outputsByName(result)
	assertSeries(t, outputs["PAST"], []float64{math.NaN(), 105, 103, 107, 110, 108, 111, 109, 112, 115})
	assertSeries(t, outputs["FUTURE"], []float64{103, 107, 110, 108, 111, 109, 112, 115, 113, math.NaN()})
	assertSeries(t, outputs["FUTURE_V"], []float64{107, 110, 108, 111, 109, 112, 115, 113, math.NaN(), math.NaN()})
	assertSeries(t, outputs["CURR"], []float64{10, 9, 8, 7, 6, 5, 4, 3, 2, 1})
	assertSeries(t, outputs["TOTAL"], []float64{10, 10, 10, 10, 10, 10, 10, 10, 10, 10})
	assertSeries(t, outputs["LAST_BAR"], []float64{0, 0, 0, 0, 0, 0, 0, 0, 0, 1})
	assertSeries(t, outputs["STATUS"], []float64{1, 2, 2, 2, 2, 2, 2, 2, 2, 3})
	assertSeries(t, outputs["SUM_BARS"], []float64{math.NaN(), math.NaN(), 3, 2, 2, 2, 2, 2, 2, 2})
}

func outputsByName(result *types.FormulaResult) map[string][]float64 {
	outputs := make(map[string][]float64)
	for _, output := range result.Outputs {
		outputs[output.Name] = output.Data
	}
	return outputs
}

func assertClose(t *testing.T, got, want float64) {
	t.Helper()
	if math.Abs(got-want) > 0.0001 {
		t.Errorf("expected %.6f, got %.6f", want, got)
	}
}

func assertVariable(t *testing.T, result *types.FormulaResult, name string, want float64) {
	t.Helper()
	got, ok := result.Variables[name]
	if !ok {
		t.Fatalf("expected variable %s", name)
	}
	assertClose(t, got, want)
}

func assertSeries(t *testing.T, got, want []float64) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("expected %d data points, got %d", len(want), len(got))
	}
	for i := range want {
		if math.IsNaN(want[i]) {
			if !math.IsNaN(got[i]) {
				t.Errorf("index %d: expected NaN, got %f", i, got[i])
			}
			continue
		}
		if math.Abs(got[i]-want[i]) > 0.0001 {
			t.Errorf("index %d: expected %.6f, got %.6f", i, want[i], got[i])
		}
	}
}
