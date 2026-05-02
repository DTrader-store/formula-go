package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"sort"
	"time"

	"github.com/DTrader-store/formula-go"
	"github.com/DTrader-store/formula-go/examples/helpers"
)

const resistanceSupportFormula = `
Z1 := 20;
M1 := 3;
阻力1:MA(REF(HHV(H,Z1),1),M1),COLORRED;
阻力2:=MA(REF(HHV(H,15*Z1),1),M1),COLORCYAN;
支撑1:MA(REF(LLV(L,Z1),1),M1),COLORGREEN;
支撑2:=MA(REF(LLV(L,15*Z1),1),M1),COLORBLUE;
现价:C, COLORBLACK;
最高:H;
最低:L;
DRAWTEXT(CROSS(阻力1,H),C,'S');
DRAWTEXT(CROSS(L,支撑1),C,'B');
`

func main() {
	code := flag.String("code", "sz000001", "TDX stock code, for example sz000001 or sh600000")
	count := flag.Uint("count", 360, "daily K-line count to fetch")
	start := flag.Uint("start", 0, "TDX K-line start offset")
	flag.Parse()

	client, err := helpers.NewTDXClient(helpers.DefaultTDXServer())
	if err != nil {
		log.Fatalf("create TDX client: %v", err)
	}
	defer client.Close()

	klines, err := client.GetDailyKlines(*code, uint16(*start), uint16(*count))
	if err != nil {
		log.Fatalf("fetch TDX market data: %v", err)
	}
	marketData := klines.Data
	if len(marketData) == 0 {
		log.Fatalf("TDX returned no market data for %s", *code)
	}

	engine := formula.NewFormulaEngine()
	result, err := engine.Run(resistanceSupportFormula, marketData)
	if err != nil {
		log.Fatalf("execute formula: %v", err)
	}

	fmt.Println("=== TDX Resistance / Support Demo ===")
	fmt.Println("Source: TDX real daily K-line data via github.com/injoyai/tdx")
	fmt.Printf("Server: %s\n", helpers.DefaultTDXServer())
	fmt.Printf("Code: %s, bars: %d, start: %d\n\n", *code, len(marketData), *start)

	fmt.Println("Formula:")
	fmt.Print(resistanceSupportFormula)
	fmt.Println()

	printLatestOutputs(result, klines.Times)
	printRecentSignals(result, marketData, klines.Times, 10)
}

func printLatestOutputs(result *formula.FormulaResult, times []time.Time) {
	fmt.Println("Latest output values:")
	if len(times) > 0 {
		fmt.Printf("  %-8s: %s\n", "日期", formatKlineTime(times[len(times)-1]))
	}
	for _, output := range result.Outputs {
		value := lastFinite(output.Data)
		if math.IsNaN(value) {
			fmt.Printf("  %-8s: NaN\n", output.Name)
			continue
		}
		fmt.Printf("  %-8s: %.3f\n", output.Name, value)
	}
	fmt.Println()
}

func printRecentSignals(result *formula.FormulaResult, data []*formula.MarketData, times []time.Time, limit int) {
	signals := make([]*formula.DrawingEvent, 0, len(result.Drawings))
	for _, drawing := range result.Drawings {
		if drawing.Function == "DRAWTEXT" && (drawing.Text == "S" || drawing.Text == "B") {
			signals = append(signals, drawing)
		}
	}
	sort.Slice(signals, func(i, j int) bool {
		return signals[i].BarIndex < signals[j].BarIndex
	})

	fmt.Printf("DRAWTEXT signals: %d total\n", len(signals))
	if len(signals) == 0 {
		return
	}

	start := len(signals) - limit
	if start < 0 {
		start = 0
	}
	fmt.Printf("Recent %d signals:\n", len(signals)-start)
	for _, signal := range signals[start:] {
		bar := signal.BarIndex
		if bar < 0 || bar >= len(data) || bar >= len(times) {
			continue
		}
		fmt.Printf("  date=%s signal=%s close=%.3f high=%.3f low=%.3f eventPrice=%.3f\n",
			formatKlineTime(times[bar]),
			signal.Text,
			data[bar].Close,
			data[bar].High,
			data[bar].Low,
			signal.Values["price"],
		)
	}
}

func formatKlineTime(t time.Time) string {
	if t.IsZero() {
		return "unknown"
	}
	return t.Format("2006-01-02")
}

func lastFinite(values []float64) float64 {
	for i := len(values) - 1; i >= 0; i-- {
		if !math.IsNaN(values[i]) {
			return values[i]
		}
	}
	return math.NaN()
}
