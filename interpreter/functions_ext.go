package interpreter

import (
	"fmt"
	"math"

	"github.com/DTrader-store/formula-go/errors"
	"github.com/DTrader-store/formula-go/types"
)

// Additional built-in functions for Phase 4

// fnSTD implements Standard Deviation: STD(data, period)
func fnSTD(args []*Value, _ []*types.MarketData) (*Value, error) {
	if len(args) != 2 {
		return nil, errors.NewRuntimeError("STD requires 2 arguments")
	}

	data := args[0]
	period := args[1]

	if !data.IsArray {
		return nil, errors.NewRuntimeError("STD first argument must be an array")
	}
	if period.IsArray {
		return nil, errors.NewRuntimeError("STD second argument must be a number")
	}

	n := int(period.Single)
	if n <= 0 || n > len(data.Array) {
		return nil, errors.NewRuntimeError(fmt.Sprintf("STD period must be between 1 and %d", len(data.Array)))
	}

	result := make([]float64, len(data.Array))

	// Fill first n-1 values with NaN
	for i := 0; i < n-1; i++ {
		result[i] = math.NaN()
	}

	// Calculate STD
	for i := n - 1; i < len(data.Array); i++ {
		// Calculate mean
		mean := 0.0
		for j := 0; j < n; j++ {
			mean += data.Array[i-j]
		}
		mean /= float64(n)

		// Calculate variance
		variance := 0.0
		for j := 0; j < n; j++ {
			diff := data.Array[i-j] - mean
			variance += diff * diff
		}
		variance /= float64(n)

		result[i] = math.Sqrt(variance)
	}

	return NewArrayValue(result), nil
}

// fnVAR implements Variance: VAR(data, period)
func fnVAR(args []*Value, _ []*types.MarketData) (*Value, error) {
	if len(args) != 2 {
		return nil, errors.NewRuntimeError("VAR requires 2 arguments")
	}

	data := args[0]
	period := args[1]

	if !data.IsArray {
		return nil, errors.NewRuntimeError("VAR first argument must be an array")
	}
	if period.IsArray {
		return nil, errors.NewRuntimeError("VAR second argument must be a number")
	}

	n := int(period.Single)
	if n <= 0 || n > len(data.Array) {
		return nil, errors.NewRuntimeError(fmt.Sprintf("VAR period must be between 1 and %d", len(data.Array)))
	}

	result := make([]float64, len(data.Array))

	// Fill first n-1 values with NaN
	for i := 0; i < n-1; i++ {
		result[i] = math.NaN()
	}

	// Calculate VAR
	for i := n - 1; i < len(data.Array); i++ {
		// Calculate mean
		mean := 0.0
		for j := 0; j < n; j++ {
			mean += data.Array[i-j]
		}
		mean /= float64(n)

		// Calculate variance
		variance := 0.0
		for j := 0; j < n; j++ {
			diff := data.Array[i-j] - mean
			variance += diff * diff
		}
		result[i] = variance / float64(n)
	}

	return NewArrayValue(result), nil
}

// fnSMA implements SMA(data, period) as a MA alias for backward compatibility,
// and SMA(data, period, weight) as the TDX recursive smoothing formula.
func fnSMA(args []*Value, data []*types.MarketData) (*Value, error) {
	if len(args) == 2 {
		return fnMA(args, data)
	}
	if len(args) != 3 {
		return nil, errors.NewRuntimeError("SMA requires 2 or 3 arguments")
	}

	source := args[0]
	period := args[1]
	weight := args[2]

	if !source.IsArray {
		return nil, errors.NewRuntimeError("SMA first argument must be an array")
	}
	if period.IsArray {
		return nil, errors.NewRuntimeError("SMA second argument must be a number")
	}
	if weight.IsArray {
		return nil, errors.NewRuntimeError("SMA third argument must be a number")
	}

	n := int(period.Single)
	m := int(weight.Single)
	if n <= 0 {
		return nil, errors.NewRuntimeError("SMA period must be positive")
	}
	if m <= 0 || m > n {
		return nil, errors.NewRuntimeError("SMA weight must be between 1 and period")
	}
	if len(source.Array) == 0 {
		return NewArrayValue([]float64{}), nil
	}

	result := make([]float64, len(source.Array))
	result[0] = source.Array[0]
	for i := 1; i < len(source.Array); i++ {
		result[i] = (float64(m)*source.Array[i] + float64(n-m)*result[i-1]) / float64(n)
	}

	return NewArrayValue(result), nil
}

// fnWMA implements Weighted Moving Average: WMA(data, period)
func fnWMA(args []*Value, _ []*types.MarketData) (*Value, error) {
	if len(args) != 2 {
		return nil, errors.NewRuntimeError("WMA requires 2 arguments")
	}

	data := args[0]
	period := args[1]

	if !data.IsArray {
		return nil, errors.NewRuntimeError("WMA first argument must be an array")
	}
	if period.IsArray {
		return nil, errors.NewRuntimeError("WMA second argument must be a number")
	}

	n := int(period.Single)
	if n <= 0 || n > len(data.Array) {
		return nil, errors.NewRuntimeError(fmt.Sprintf("WMA period must be between 1 and %d", len(data.Array)))
	}

	result := make([]float64, len(data.Array))

	// Calculate weight sum
	weightSum := float64(n * (n + 1) / 2)

	// Fill first n-1 values with NaN
	for i := 0; i < n-1; i++ {
		result[i] = math.NaN()
	}

	// Calculate WMA
	for i := n - 1; i < len(data.Array); i++ {
		weightedSum := 0.0
		for j := 0; j < n; j++ {
			weight := float64(n - j)
			weightedSum += data.Array[i-j] * weight
		}
		result[i] = weightedSum / weightSum
	}

	return NewArrayValue(result), nil
}

// fnCOUNT implements Count: COUNT(condition, period)
func fnCOUNT(args []*Value, _ []*types.MarketData) (*Value, error) {
	if len(args) != 2 {
		return nil, errors.NewRuntimeError("COUNT requires 2 arguments")
	}

	condition := args[0]
	period := args[1]

	if !condition.IsArray {
		return nil, errors.NewRuntimeError("COUNT first argument must be an array")
	}
	if period.IsArray {
		return nil, errors.NewRuntimeError("COUNT second argument must be a number")
	}

	n := int(period.Single)
	if n <= 0 || n > len(condition.Array) {
		return nil, errors.NewRuntimeError(fmt.Sprintf("COUNT period must be between 1 and %d", len(condition.Array)))
	}

	result := make([]float64, len(condition.Array))

	// Fill first n-1 values with NaN
	for i := 0; i < n-1; i++ {
		result[i] = math.NaN()
	}

	// Count true conditions
	for i := n - 1; i < len(condition.Array); i++ {
		count := 0.0
		for j := 0; j < n; j++ {
			if isTruthy(condition.Array[i-j]) {
				count++
			}
		}
		result[i] = count
	}

	return NewArrayValue(result), nil
}

// fnEVERY implements Every: EVERY(condition, period) - returns 1 if condition is true for all periods
func fnEVERY(args []*Value, _ []*types.MarketData) (*Value, error) {
	if len(args) != 2 {
		return nil, errors.NewRuntimeError("EVERY requires 2 arguments")
	}

	condition := args[0]
	period := args[1]

	if !condition.IsArray {
		return nil, errors.NewRuntimeError("EVERY first argument must be an array")
	}
	if period.IsArray {
		return nil, errors.NewRuntimeError("EVERY second argument must be a number")
	}

	n := int(period.Single)
	if n <= 0 || n > len(condition.Array) {
		return nil, errors.NewRuntimeError(fmt.Sprintf("EVERY period must be between 1 and %d", len(condition.Array)))
	}

	result := make([]float64, len(condition.Array))

	// Fill first n-1 values with 0
	for i := 0; i < n-1; i++ {
		result[i] = 0
	}

	// Check if every condition is true
	for i := n - 1; i < len(condition.Array); i++ {
		everyCond := true
		for j := 0; j < n; j++ {
			if !isTruthy(condition.Array[i-j]) {
				everyCond = false
				break
			}
		}
		if everyCond {
			result[i] = 1
		} else {
			result[i] = 0
		}
	}

	return NewArrayValue(result), nil
}

// fnEXIST implements Exist: EXIST(condition, period) - returns 1 if condition is true for any period
func fnEXIST(args []*Value, _ []*types.MarketData) (*Value, error) {
	if len(args) != 2 {
		return nil, errors.NewRuntimeError("EXIST requires 2 arguments")
	}

	condition := args[0]
	period := args[1]

	if !condition.IsArray {
		return nil, errors.NewRuntimeError("EXIST first argument must be an array")
	}
	if period.IsArray {
		return nil, errors.NewRuntimeError("EXIST second argument must be a number")
	}

	n := int(period.Single)
	if n <= 0 || n > len(condition.Array) {
		return nil, errors.NewRuntimeError(fmt.Sprintf("EXIST period must be between 1 and %d", len(condition.Array)))
	}

	result := make([]float64, len(condition.Array))

	// Fill first n-1 values with 0
	for i := 0; i < n-1; i++ {
		result[i] = 0
	}

	// Check if any condition is true
	for i := n - 1; i < len(condition.Array); i++ {
		exists := false
		for j := 0; j < n; j++ {
			if isTruthy(condition.Array[i-j]) {
				exists = true
				break
			}
		}
		if exists {
			result[i] = 1
		} else {
			result[i] = 0
		}
	}

	return NewArrayValue(result), nil
}

// fnBARSLAST implements BarsLast: BARSLAST(condition) - returns bars since last true condition
func fnBARSLAST(args []*Value, _ []*types.MarketData) (*Value, error) {
	if len(args) != 1 {
		return nil, errors.NewRuntimeError("BARSLAST requires 1 argument")
	}

	condition := args[0]

	if !condition.IsArray {
		return nil, errors.NewRuntimeError("BARSLAST argument must be an array")
	}

	result := make([]float64, len(condition.Array))
	lastTrueIndex := -1

	for i := 0; i < len(condition.Array); i++ {
		if isTruthy(condition.Array[i]) {
			lastTrueIndex = i
			result[i] = 0
		} else if lastTrueIndex >= 0 {
			result[i] = float64(i - lastTrueIndex)
		} else {
			result[i] = math.NaN()
		}
	}

	return NewArrayValue(result), nil
}

// fnHHVBARS implements HHVBARS(data, period) - bars since the highest value in the window.
func fnHHVBARS(args []*Value, _ []*types.MarketData) (*Value, error) {
	if len(args) != 2 {
		return nil, errors.NewRuntimeError("HHVBARS requires 2 arguments")
	}

	data := args[0]
	period := args[1]

	if !data.IsArray {
		return nil, errors.NewRuntimeError("HHVBARS first argument must be an array")
	}
	if period.IsArray {
		return nil, errors.NewRuntimeError("HHVBARS second argument must be a number")
	}

	n := int(period.Single)
	if n <= 0 || n > len(data.Array) {
		return nil, errors.NewRuntimeError(fmt.Sprintf("HHVBARS period must be between 1 and %d", len(data.Array)))
	}

	result := make([]float64, len(data.Array))
	for i := 0; i < n-1; i++ {
		result[i] = math.NaN()
	}

	for i := n - 1; i < len(data.Array); i++ {
		maxValue := data.Array[i]
		bars := 0
		for j := 1; j < n; j++ {
			if data.Array[i-j] > maxValue {
				maxValue = data.Array[i-j]
				bars = j
			}
		}
		result[i] = float64(bars)
	}

	return NewArrayValue(result), nil
}

// fnLLVBARS implements LLVBARS(data, period) - bars since the lowest value in the window.
func fnLLVBARS(args []*Value, _ []*types.MarketData) (*Value, error) {
	if len(args) != 2 {
		return nil, errors.NewRuntimeError("LLVBARS requires 2 arguments")
	}

	data := args[0]
	period := args[1]

	if !data.IsArray {
		return nil, errors.NewRuntimeError("LLVBARS first argument must be an array")
	}
	if period.IsArray {
		return nil, errors.NewRuntimeError("LLVBARS second argument must be a number")
	}

	n := int(period.Single)
	if n <= 0 || n > len(data.Array) {
		return nil, errors.NewRuntimeError(fmt.Sprintf("LLVBARS period must be between 1 and %d", len(data.Array)))
	}

	result := make([]float64, len(data.Array))
	for i := 0; i < n-1; i++ {
		result[i] = math.NaN()
	}

	for i := n - 1; i < len(data.Array); i++ {
		minValue := data.Array[i]
		bars := 0
		for j := 1; j < n; j++ {
			if data.Array[i-j] < minValue {
				minValue = data.Array[i-j]
				bars = j
			}
		}
		result[i] = float64(bars)
	}

	return NewArrayValue(result), nil
}

// fnBARSCOUNT implements BARSCOUNT(data) - count of valid bars seen so far.
func fnBARSCOUNT(args []*Value, _ []*types.MarketData) (*Value, error) {
	if len(args) != 1 {
		return nil, errors.NewRuntimeError("BARSCOUNT requires 1 argument")
	}

	data := args[0]
	if !data.IsArray {
		return nil, errors.NewRuntimeError("BARSCOUNT argument must be an array")
	}

	result := make([]float64, len(data.Array))
	count := 0
	for i, value := range data.Array {
		if !math.IsNaN(value) {
			count++
		}
		result[i] = float64(count)
	}

	return NewArrayValue(result), nil
}

// fnBARSSINCE implements BARSSINCE(condition) - bars since the first true condition.
func fnBARSSINCE(args []*Value, _ []*types.MarketData) (*Value, error) {
	if len(args) != 1 {
		return nil, errors.NewRuntimeError("BARSSINCE requires 1 argument")
	}

	condition := args[0]
	if !condition.IsArray {
		return nil, errors.NewRuntimeError("BARSSINCE argument must be an array")
	}

	result := make([]float64, len(condition.Array))
	firstTrueIndex := -1
	for i, value := range condition.Array {
		if firstTrueIndex < 0 {
			if isTruthy(value) {
				firstTrueIndex = i
				result[i] = 0
			} else {
				result[i] = math.NaN()
			}
			continue
		}
		result[i] = float64(i - firstTrueIndex)
	}

	return NewArrayValue(result), nil
}

// fnBARSLASTCOUNT implements BARSLASTCOUNT(condition) - consecutive true count ending at each bar.
func fnBARSLASTCOUNT(args []*Value, _ []*types.MarketData) (*Value, error) {
	if len(args) != 1 {
		return nil, errors.NewRuntimeError("BARSLASTCOUNT requires 1 argument")
	}

	condition := args[0]
	if !condition.IsArray {
		return nil, errors.NewRuntimeError("BARSLASTCOUNT argument must be an array")
	}

	result := make([]float64, len(condition.Array))
	count := 0
	for i, value := range condition.Array {
		if isTruthy(value) {
			count++
		} else {
			count = 0
		}
		result[i] = float64(count)
	}

	return NewArrayValue(result), nil
}

// fnAVEDEV implements Average Deviation: AVEDEV(data, period)
func fnAVEDEV(args []*Value, _ []*types.MarketData) (*Value, error) {
	if len(args) != 2 {
		return nil, errors.NewRuntimeError("AVEDEV requires 2 arguments")
	}

	data := args[0]
	period := args[1]

	if !data.IsArray {
		return nil, errors.NewRuntimeError("AVEDEV first argument must be an array")
	}
	if period.IsArray {
		return nil, errors.NewRuntimeError("AVEDEV second argument must be a number")
	}

	n := int(period.Single)
	if n <= 0 || n > len(data.Array) {
		return nil, errors.NewRuntimeError(fmt.Sprintf("AVEDEV period must be between 1 and %d", len(data.Array)))
	}

	result := make([]float64, len(data.Array))

	// Fill first n-1 values with NaN
	for i := 0; i < n-1; i++ {
		result[i] = math.NaN()
	}

	// Calculate AVEDEV
	for i := n - 1; i < len(data.Array); i++ {
		// Calculate mean
		mean := 0.0
		for j := 0; j < n; j++ {
			mean += data.Array[i-j]
		}
		mean /= float64(n)

		// Calculate average deviation
		devSum := 0.0
		for j := 0; j < n; j++ {
			devSum += math.Abs(data.Array[i-j] - mean)
		}
		result[i] = devSum / float64(n)
	}

	return NewArrayValue(result), nil
}

// fnFILTER implements Filter: FILTER(condition, period) - filters signals
func fnFILTER(args []*Value, _ []*types.MarketData) (*Value, error) {
	if len(args) != 2 {
		return nil, errors.NewRuntimeError("FILTER requires 2 arguments")
	}

	condition := args[0]
	period := args[1]

	if !condition.IsArray {
		return nil, errors.NewRuntimeError("FILTER first argument must be an array")
	}
	if period.IsArray {
		return nil, errors.NewRuntimeError("FILTER second argument must be a number")
	}

	n := int(period.Single)
	if n <= 0 {
		return nil, errors.NewRuntimeError("FILTER period must be positive")
	}

	result := make([]float64, len(condition.Array))
	lastSignal := -n - 1 // Initialize to allow first signal

	for i := 0; i < len(condition.Array); i++ {
		if isTruthy(condition.Array[i]) && (i-lastSignal) >= n {
			result[i] = 1
			lastSignal = i
		} else {
			result[i] = 0
		}
	}

	return NewArrayValue(result), nil
}

// fnBETWEEN implements Between: BETWEEN(value, lower, upper)
func fnBETWEEN(args []*Value, _ []*types.MarketData) (*Value, error) {
	if len(args) != 3 {
		return nil, errors.NewRuntimeError("BETWEEN requires 3 arguments")
	}

	value := args[0]
	lower := args[1]
	upper := args[2]

	// Handle scalar case
	if !value.IsArray && !lower.IsArray && !upper.IsArray {
		if value.Single >= lower.Single && value.Single <= upper.Single {
			return NewSingleValue(1), nil
		}
		return NewSingleValue(0), nil
	}

	// Handle array case
	if !value.IsArray {
		return nil, errors.NewRuntimeError("BETWEEN: value must be array when using array bounds")
	}

	result := make([]float64, len(value.Array))
	for i := range value.Array {
		lowerBound := lower.Single
		upperBound := upper.Single
		if lower.IsArray {
			lowerBound = lower.Array[i]
		}
		if upper.IsArray {
			upperBound = upper.Array[i]
		}

		if value.Array[i] >= lowerBound && value.Array[i] <= upperBound {
			result[i] = 1
		} else {
			result[i] = 0
		}
	}

	return NewArrayValue(result), nil
}

// fnRANGE implements RANGE(value, lower, upper) as a TDX-compatible alias of BETWEEN.
func fnRANGE(args []*Value, data []*types.MarketData) (*Value, error) {
	if len(args) != 3 {
		return nil, errors.NewRuntimeError("RANGE requires 3 arguments")
	}
	return fnBETWEEN(args, data)
}

// fnNOT implements NOT(value).
func fnNOT(args []*Value, _ []*types.MarketData) (*Value, error) {
	if len(args) != 1 {
		return nil, errors.NewRuntimeError("NOT requires 1 argument")
	}

	value := args[0]
	if !value.IsArray {
		if isTruthy(value.Single) {
			return NewSingleValue(0), nil
		}
		return NewSingleValue(1), nil
	}

	result := make([]float64, len(value.Array))
	for i, v := range value.Array {
		if isTruthy(v) {
			result[i] = 0
		} else {
			result[i] = 1
		}
	}
	return NewArrayValue(result), nil
}

// fnIFN implements IFN(condition, trueValue, falseValue) using the inverse condition.
func fnIFN(args []*Value, data []*types.MarketData) (*Value, error) {
	if len(args) != 3 {
		return nil, errors.NewRuntimeError("IFN requires 3 arguments")
	}
	return fnIF([]*Value{args[0], args[2], args[1]}, data)
}

// fnIFF implements IFF(condition, trueValue, falseValue) as an IF alias.
func fnIFF(args []*Value, data []*types.MarketData) (*Value, error) {
	if len(args) != 3 {
		return nil, errors.NewRuntimeError("IFF requires 3 arguments")
	}
	return fnIF(args, data)
}

// fnDRAWNULL returns NaN for chart gaps.
func fnDRAWNULL(args []*Value, _ []*types.MarketData) (*Value, error) {
	if len(args) != 0 {
		return nil, errors.NewRuntimeError("DRAWNULL requires 0 arguments")
	}
	return NewSingleValue(math.NaN()), nil
}

// fnPOW implements POW(base, exponent).
func fnPOW(args []*Value, _ []*types.MarketData) (*Value, error) {
	return numericBinaryFunc(args, "POW", math.Pow)
}

// fnEXP implements EXP(value).
func fnEXP(args []*Value, _ []*types.MarketData) (*Value, error) {
	return numericUnaryFunc(args, "EXP", math.Exp)
}

// fnLN implements LN(value).
func fnLN(args []*Value, _ []*types.MarketData) (*Value, error) {
	return numericUnaryFunc(args, "LN", math.Log)
}

// fnLOG implements LOG(value) as base-10 logarithm.
func fnLOG(args []*Value, _ []*types.MarketData) (*Value, error) {
	return numericUnaryFunc(args, "LOG", math.Log10)
}

// fnMOD implements MOD(a, b).
func fnMOD(args []*Value, _ []*types.MarketData) (*Value, error) {
	return numericBinaryFunc(args, "MOD", math.Mod)
}

// fnCEILING implements CEILING(value).
func fnCEILING(args []*Value, _ []*types.MarketData) (*Value, error) {
	return numericUnaryFunc(args, "CEILING", math.Ceil)
}

// fnFLOOR implements FLOOR(value).
func fnFLOOR(args []*Value, _ []*types.MarketData) (*Value, error) {
	return numericUnaryFunc(args, "FLOOR", math.Floor)
}

// fnINTPART implements INTPART(value).
func fnINTPART(args []*Value, _ []*types.MarketData) (*Value, error) {
	return numericUnaryFunc(args, "INTPART", math.Trunc)
}

// fnFRACPART implements FRACPART(value).
func fnFRACPART(args []*Value, _ []*types.MarketData) (*Value, error) {
	return numericUnaryFunc(args, "FRACPART", func(v float64) float64 {
		return v - math.Trunc(v)
	})
}

// fnROUND implements ROUND(value).
func fnROUND(args []*Value, _ []*types.MarketData) (*Value, error) {
	return numericUnaryFunc(args, "ROUND", math.Round)
}

// fnROUND2 implements ROUND2(value, digits).
func fnROUND2(args []*Value, _ []*types.MarketData) (*Value, error) {
	if len(args) != 2 {
		return nil, errors.NewRuntimeError("ROUND2 requires 2 arguments")
	}
	if args[1].IsArray {
		return nil, errors.NewRuntimeError("ROUND2 second argument must be a number")
	}

	scale := math.Pow(10, args[1].Single)
	return numericUnaryFunc([]*Value{args[0]}, "ROUND2", func(v float64) float64 {
		return math.Round(v*scale) / scale
	})
}

// fnSIGN implements SIGN(value).
func fnSIGN(args []*Value, _ []*types.MarketData) (*Value, error) {
	return numericUnaryFunc(args, "SIGN", func(v float64) float64 {
		if v > 0 {
			return 1
		}
		if v < 0 {
			return -1
		}
		return 0
	})
}

// fnSIN implements SIN(value).
func fnSIN(args []*Value, _ []*types.MarketData) (*Value, error) {
	return numericUnaryFunc(args, "SIN", math.Sin)
}

// fnCOS implements COS(value).
func fnCOS(args []*Value, _ []*types.MarketData) (*Value, error) {
	return numericUnaryFunc(args, "COS", math.Cos)
}

// fnTAN implements TAN(value).
func fnTAN(args []*Value, _ []*types.MarketData) (*Value, error) {
	return numericUnaryFunc(args, "TAN", math.Tan)
}

// fnASIN implements ASIN(value).
func fnASIN(args []*Value, _ []*types.MarketData) (*Value, error) {
	return numericUnaryFunc(args, "ASIN", math.Asin)
}

// fnACOS implements ACOS(value).
func fnACOS(args []*Value, _ []*types.MarketData) (*Value, error) {
	return numericUnaryFunc(args, "ACOS", math.Acos)
}

// fnATAN implements ATAN(value).
func fnATAN(args []*Value, _ []*types.MarketData) (*Value, error) {
	return numericUnaryFunc(args, "ATAN", math.Atan)
}

// fnLONGCROSS implements LONGCROSS(a, b, period).
func fnLONGCROSS(args []*Value, _ []*types.MarketData) (*Value, error) {
	if len(args) != 3 {
		return nil, errors.NewRuntimeError("LONGCROSS requires 3 arguments")
	}

	a, b, period := args[0], args[1], args[2]
	if !a.IsArray || !b.IsArray {
		return nil, errors.NewRuntimeError("LONGCROSS first two arguments must be arrays")
	}
	if period.IsArray {
		return nil, errors.NewRuntimeError("LONGCROSS third argument must be a number")
	}
	if len(a.Array) != len(b.Array) {
		return nil, errors.NewRuntimeError("LONGCROSS: array length mismatch")
	}

	n := int(period.Single)
	if n <= 0 {
		return nil, errors.NewRuntimeError("LONGCROSS period must be positive")
	}

	result := make([]float64, len(a.Array))
	for i := 1; i < len(a.Array); i++ {
		if !(a.Array[i-1] <= b.Array[i-1] && a.Array[i] > b.Array[i]) || i < n {
			continue
		}
		ok := true
		for j := 1; j <= n; j++ {
			if a.Array[i-j] >= b.Array[i-j] {
				ok = false
				break
			}
		}
		if ok {
			result[i] = 1
		}
	}
	return NewArrayValue(result), nil
}

// fnUPNDAY implements UPNDAY(data, period).
func fnUPNDAY(args []*Value, _ []*types.MarketData) (*Value, error) {
	if len(args) != 2 {
		return nil, errors.NewRuntimeError("UPNDAY requires 2 arguments")
	}
	return compareConsecutive(args, "UPNDAY", func(curr, prev float64) bool { return curr > prev })
}

// fnDOWNNDAY implements DOWNNDAY(data, period).
func fnDOWNNDAY(args []*Value, _ []*types.MarketData) (*Value, error) {
	if len(args) != 2 {
		return nil, errors.NewRuntimeError("DOWNNDAY requires 2 arguments")
	}
	return compareConsecutive(args, "DOWNNDAY", func(curr, prev float64) bool { return curr < prev })
}

// fnNDAY implements NDAY(a, b, period) - a has been greater than b for period bars.
func fnNDAY(args []*Value, _ []*types.MarketData) (*Value, error) {
	if len(args) != 3 {
		return nil, errors.NewRuntimeError("NDAY requires 3 arguments")
	}

	a, b, period := args[0], args[1], args[2]
	if !a.IsArray || !b.IsArray {
		return nil, errors.NewRuntimeError("NDAY first two arguments must be arrays")
	}
	if period.IsArray {
		return nil, errors.NewRuntimeError("NDAY third argument must be a number")
	}
	if len(a.Array) != len(b.Array) {
		return nil, errors.NewRuntimeError("NDAY: array length mismatch")
	}

	n := int(period.Single)
	if n <= 0 {
		return nil, errors.NewRuntimeError("NDAY period must be positive")
	}

	result := make([]float64, len(a.Array))
	for i := n - 1; i < len(a.Array); i++ {
		ok := true
		for j := 0; j < n; j++ {
			if !(a.Array[i-j] > b.Array[i-j]) {
				ok = false
				break
			}
		}
		if ok {
			result[i] = 1
		}
	}

	return NewArrayValue(result), nil
}

// fnLAST implements LAST(condition, from, to).
func fnLAST(args []*Value, _ []*types.MarketData) (*Value, error) {
	if len(args) != 3 {
		return nil, errors.NewRuntimeError("LAST requires 3 arguments")
	}

	condition, from, to := args[0], args[1], args[2]
	if !condition.IsArray {
		return nil, errors.NewRuntimeError("LAST first argument must be an array")
	}
	if from.IsArray || to.IsArray {
		return nil, errors.NewRuntimeError("LAST second and third arguments must be numbers")
	}

	fromN := int(from.Single)
	toN := int(to.Single)
	if fromN < toN || toN < 0 {
		return nil, errors.NewRuntimeError("LAST requires from >= to >= 0")
	}

	result := make([]float64, len(condition.Array))
	for i := range condition.Array {
		if i < fromN {
			continue
		}
		ok := true
		for j := toN; j <= fromN; j++ {
			if !isTruthy(condition.Array[i-j]) {
				ok = false
				break
			}
		}
		if ok {
			result[i] = 1
		}
	}
	return NewArrayValue(result), nil
}

// fnEXISTR implements EXISTR(condition, from, to).
func fnEXISTR(args []*Value, _ []*types.MarketData) (*Value, error) {
	if len(args) != 3 {
		return nil, errors.NewRuntimeError("EXISTR requires 3 arguments")
	}

	condition, from, to := args[0], args[1], args[2]
	if !condition.IsArray {
		return nil, errors.NewRuntimeError("EXISTR first argument must be an array")
	}
	if from.IsArray || to.IsArray {
		return nil, errors.NewRuntimeError("EXISTR second and third arguments must be numbers")
	}

	fromN := int(from.Single)
	toN := int(to.Single)
	if fromN < toN || toN < 0 {
		return nil, errors.NewRuntimeError("EXISTR requires from >= to >= 0")
	}

	result := make([]float64, len(condition.Array))
	for i := range condition.Array {
		if i < fromN {
			continue
		}
		for j := toN; j <= fromN; j++ {
			if isTruthy(condition.Array[i-j]) {
				result[i] = 1
				break
			}
		}
	}
	return NewArrayValue(result), nil
}

// fnDMA implements DMA(data, alpha) - dynamic moving average.
func fnDMA(args []*Value, _ []*types.MarketData) (*Value, error) {
	if len(args) != 2 {
		return nil, errors.NewRuntimeError("DMA requires 2 arguments")
	}

	data := args[0]
	alpha := args[1]

	if !data.IsArray {
		return nil, errors.NewRuntimeError("DMA first argument must be an array")
	}
	if alpha.IsArray && len(alpha.Array) != len(data.Array) {
		return nil, errors.NewRuntimeError("DMA: array length mismatch")
	}
	if len(data.Array) == 0 {
		return NewArrayValue([]float64{}), nil
	}

	result := make([]float64, len(data.Array))
	result[0] = data.Array[0]
	for i := 1; i < len(data.Array); i++ {
		a := alpha.Single
		if alpha.IsArray {
			a = alpha.Array[i]
		}
		result[i] = a*data.Array[i] + (1-a)*result[i-1]
	}

	return NewArrayValue(result), nil
}

// fnCONST implements CONST(value) - fill all bars with the final value.
func fnCONST(args []*Value, _ []*types.MarketData) (*Value, error) {
	if len(args) != 1 {
		return nil, errors.NewRuntimeError("CONST requires 1 argument")
	}

	value := args[0]
	if !value.IsArray {
		return NewSingleValue(value.Single), nil
	}
	if len(value.Array) == 0 {
		return NewArrayValue([]float64{}), nil
	}

	lastValue := value.Array[len(value.Array)-1]
	result := make([]float64, len(value.Array))
	for i := range result {
		result[i] = lastValue
	}

	return NewArrayValue(result), nil
}

// fnVALUEWHEN implements VALUEWHEN(condition, value) - hold value from the latest true condition.
func fnVALUEWHEN(args []*Value, _ []*types.MarketData) (*Value, error) {
	if len(args) != 2 {
		return nil, errors.NewRuntimeError("VALUEWHEN requires 2 arguments")
	}

	condition := args[0]
	value := args[1]

	if !condition.IsArray {
		return nil, errors.NewRuntimeError("VALUEWHEN first argument must be an array")
	}
	if value.IsArray && len(value.Array) != len(condition.Array) {
		return nil, errors.NewRuntimeError("VALUEWHEN: array length mismatch")
	}

	result := make([]float64, len(condition.Array))
	lastValue := math.NaN()
	hasValue := false
	for i, cond := range condition.Array {
		if isTruthy(cond) {
			if value.IsArray {
				lastValue = value.Array[i]
			} else {
				lastValue = value.Single
			}
			hasValue = true
		}

		if hasValue {
			result[i] = lastValue
		} else {
			result[i] = math.NaN()
		}
	}

	return NewArrayValue(result), nil
}

// fnDEVSQ implements DEVSQ(data, period).
func fnDEVSQ(args []*Value, _ []*types.MarketData) (*Value, error) {
	return rollingStatsFunc(args, "DEVSQ", func(values []float64) float64 {
		mean := average(values)
		sum := 0.0
		for _, v := range values {
			diff := v - mean
			sum += diff * diff
		}
		return sum
	})
}

// fnFORCAST implements FORCAST(data, period) using linear regression projection at the current bar.
func fnFORCAST(args []*Value, _ []*types.MarketData) (*Value, error) {
	return rollingRegressionFunc(args, "FORCAST", func(slope, intercept float64, n int) float64 {
		return intercept + slope*float64(n-1)
	})
}

// fnSLOPE implements SLOPE(data, period).
func fnSLOPE(args []*Value, _ []*types.MarketData) (*Value, error) {
	return rollingRegressionFunc(args, "SLOPE", func(slope, _ float64, _ int) float64 {
		return slope
	})
}

// fnSTDP implements STDP(data, period) as population standard deviation.
func fnSTDP(args []*Value, data []*types.MarketData) (*Value, error) {
	if len(args) != 2 {
		return nil, errors.NewRuntimeError("STDP requires 2 arguments")
	}
	return fnSTD(args, data)
}

// fnSTDDEV implements STDDEV(data, period) as sample standard deviation.
func fnSTDDEV(args []*Value, _ []*types.MarketData) (*Value, error) {
	return rollingStatsFunc(args, "STDDEV", func(values []float64) float64 {
		if len(values) < 2 {
			return 0
		}
		mean := average(values)
		sum := 0.0
		for _, v := range values {
			diff := v - mean
			sum += diff * diff
		}
		return math.Sqrt(sum / float64(len(values)-1))
	})
}

// fnVARP implements VARP(data, period) as population variance.
func fnVARP(args []*Value, data []*types.MarketData) (*Value, error) {
	if len(args) != 2 {
		return nil, errors.NewRuntimeError("VARP requires 2 arguments")
	}
	return fnVAR(args, data)
}

// fnCOVAR implements COVAR(a, b, period) as population covariance.
func fnCOVAR(args []*Value, _ []*types.MarketData) (*Value, error) {
	return rollingPairStatsFunc(args, "COVAR", func(a, b []float64) float64 {
		meanA := average(a)
		meanB := average(b)
		sum := 0.0
		for i := range a {
			sum += (a[i] - meanA) * (b[i] - meanB)
		}
		return sum / float64(len(a))
	})
}

// fnRELATE implements RELATE(a, b, period) as correlation coefficient.
func fnRELATE(args []*Value, _ []*types.MarketData) (*Value, error) {
	return rollingPairStatsFunc(args, "RELATE", func(a, b []float64) float64 {
		cov := covariance(a, b)
		varA := variance(a)
		varB := variance(b)
		if varA == 0 || varB == 0 {
			return math.NaN()
		}
		return cov / math.Sqrt(varA*varB)
	})
}

// fnBETA implements BETA(a, b, period) as covariance(a,b) / variance(b).
func fnBETA(args []*Value, _ []*types.MarketData) (*Value, error) {
	return rollingPairStatsFunc(args, "BETA", func(a, b []float64) float64 {
		varB := variance(b)
		if varB == 0 {
			return math.NaN()
		}
		return covariance(a, b) / varB
	})
}

// fnREFV implements REFV(data, n) without future-function marking.
func fnREFV(args []*Value, data []*types.MarketData) (*Value, error) {
	if len(args) != 2 {
		return nil, errors.NewRuntimeError("REFV requires 2 arguments")
	}
	return fnREF(args, data)
}

// fnREFX implements REFX(data, n) - future reference.
func fnREFX(args []*Value, _ []*types.MarketData) (*Value, error) {
	return futureReference(args, "REFX")
}

// fnREFXV implements REFXV(data, n) - future reference without future-function marking.
func fnREFXV(args []*Value, _ []*types.MarketData) (*Value, error) {
	return futureReference(args, "REFXV")
}

// fnCURRBARSCOUNT implements CURRBARSCOUNT().
func fnCURRBARSCOUNT(args []*Value, data []*types.MarketData) (*Value, error) {
	if len(args) != 0 {
		return nil, errors.NewRuntimeError("CURRBARSCOUNT requires 0 arguments")
	}
	n := len(data)
	result := make([]float64, n)
	for i := 0; i < n; i++ {
		result[i] = float64(n - i)
	}
	return NewArrayValue(result), nil
}

// fnTOTALBARSCOUNT implements TOTALBARSCOUNT().
func fnTOTALBARSCOUNT(args []*Value, data []*types.MarketData) (*Value, error) {
	if len(args) != 0 {
		return nil, errors.NewRuntimeError("TOTALBARSCOUNT requires 0 arguments")
	}
	n := len(data)
	result := make([]float64, n)
	for i := range result {
		result[i] = float64(n)
	}
	return NewArrayValue(result), nil
}

// fnISLASTBAR implements ISLASTBAR().
func fnISLASTBAR(args []*Value, data []*types.MarketData) (*Value, error) {
	if len(args) != 0 {
		return nil, errors.NewRuntimeError("ISLASTBAR requires 0 arguments")
	}
	result := make([]float64, len(data))
	if len(result) > 0 {
		result[len(result)-1] = 1
	}
	return NewArrayValue(result), nil
}

// fnBARSTATUS implements BARSTATUS().
func fnBARSTATUS(args []*Value, data []*types.MarketData) (*Value, error) {
	if len(args) != 0 {
		return nil, errors.NewRuntimeError("BARSTATUS requires 0 arguments")
	}
	result := make([]float64, len(data))
	if len(result) > 0 {
		result[0] = 1
		if len(result) > 1 {
			for i := 1; i < len(result)-1; i++ {
				result[i] = 2
			}
			result[len(result)-1] = 3
		}
	}
	return NewArrayValue(result), nil
}

// fnSUMBARS implements SUMBARS(data, target).
func fnSUMBARS(args []*Value, _ []*types.MarketData) (*Value, error) {
	if len(args) != 2 {
		return nil, errors.NewRuntimeError("SUMBARS requires 2 arguments")
	}

	data := args[0]
	target := args[1]
	if !data.IsArray {
		return nil, errors.NewRuntimeError("SUMBARS first argument must be an array")
	}
	if target.IsArray {
		return nil, errors.NewRuntimeError("SUMBARS second argument must be a number")
	}

	result := make([]float64, len(data.Array))
	for i := range data.Array {
		sum := 0.0
		for j := i; j >= 0; j-- {
			sum += data.Array[j]
			if sum >= target.Single {
				result[i] = float64(i - j + 1)
				break
			}
		}
		if result[i] == 0 {
			result[i] = math.NaN()
		}
	}
	return NewArrayValue(result), nil
}

// fnDRAWTEXT implements DRAWTEXT(condition, price, text).
func fnDRAWTEXT(args []*Value, _ []*types.MarketData) (*Value, error) {
	if len(args) != 3 {
		return nil, errors.NewRuntimeError("DRAWTEXT requires 3 arguments")
	}
	return buildPointDrawings("DRAWTEXT", args[0], args[1], nil, args[2], "price")
}

// fnDRAWICON implements DRAWICON(condition, price, iconType).
func fnDRAWICON(args []*Value, _ []*types.MarketData) (*Value, error) {
	if len(args) != 3 {
		return nil, errors.NewRuntimeError("DRAWICON requires 3 arguments")
	}
	return buildPointDrawings("DRAWICON", args[0], args[1], args[2], nil, "price")
}

// fnDRAWNUMBER implements DRAWNUMBER(condition, price, number).
func fnDRAWNUMBER(args []*Value, _ []*types.MarketData) (*Value, error) {
	if len(args) != 3 {
		return nil, errors.NewRuntimeError("DRAWNUMBER requires 3 arguments")
	}
	return buildPointDrawings("DRAWNUMBER", args[0], args[1], args[2], nil, "price")
}

// fnSTICKLINE implements STICKLINE(condition, price1, price2, width, empty).
func fnSTICKLINE(args []*Value, _ []*types.MarketData) (*Value, error) {
	if len(args) != 5 {
		return nil, errors.NewRuntimeError("STICKLINE requires 5 arguments")
	}

	condition, price1, price2, width, empty := args[0], args[1], args[2], args[3], args[4]
	if !condition.IsArray {
		return nil, errors.NewRuntimeError("STICKLINE first argument must be an array")
	}
	if err := validateDrawingNumericArgs("STICKLINE", len(condition.Array), price1, price2, width, empty); err != nil {
		return nil, err
	}

	drawings := make([]*types.DrawingEvent, 0, truthyCount(condition.Array))
	for i, cond := range condition.Array {
		if !isTruthy(cond) {
			continue
		}
		event := &types.DrawingEvent{
			Function: "STICKLINE",
			BarIndex: i,
			Values:   make(map[string]float64, 4),
		}
		event.Values["price1"] = scalarOrArrayAt(price1, i)
		event.Values["price2"] = scalarOrArrayAt(price2, i)
		event.Values["width"] = scalarOrArrayAt(width, i)
		event.Values["empty"] = scalarOrArrayAt(empty, i)
		drawings = append(drawings, event)
	}

	return NewDrawingValue(drawings), nil
}

// fnDRAWLINE implements DRAWLINE(cond1, price1, cond2, price2, expand).
func fnDRAWLINE(args []*Value, _ []*types.MarketData) (*Value, error) {
	if len(args) != 5 {
		return nil, errors.NewRuntimeError("DRAWLINE requires 5 arguments")
	}

	cond1, price1, cond2, price2, expand := args[0], args[1], args[2], args[3], args[4]
	if !cond1.IsArray {
		return nil, errors.NewRuntimeError("DRAWLINE first argument must be an array")
	}
	if !cond2.IsArray {
		return nil, errors.NewRuntimeError("DRAWLINE third argument must be an array")
	}
	if len(cond1.Array) != len(cond2.Array) {
		return nil, errors.NewRuntimeError("DRAWLINE: condition array length mismatch")
	}
	if err := validateDrawingNumericArgs("DRAWLINE", len(cond1.Array), price1, price2, expand); err != nil {
		return nil, err
	}

	drawings := make([]*types.DrawingEvent, 0)
	startIndex := -1
	startPrice := math.NaN()
	for i := range cond1.Array {
		if isTruthy(cond1.Array[i]) {
			startIndex = i
			startPrice = scalarOrArrayAt(price1, i)
		}
		if startIndex < 0 || !isTruthy(cond2.Array[i]) || i < startIndex {
			continue
		}

		endPrice := scalarOrArrayAt(price2, i)
		drawings = append(drawings, &types.DrawingEvent{
			Function: "DRAWLINE",
			BarIndex: startIndex,
			Values: map[string]float64{
				"startBar":   float64(startIndex),
				"startPrice": startPrice,
				"endBar":     float64(i),
				"endPrice":   endPrice,
				"expand":     scalarOrArrayAt(expand, i),
			},
		})
		startIndex = -1
		startPrice = math.NaN()
	}

	return NewDrawingValue(drawings), nil
}

// fnPOLYLINE implements POLYLINE(condition, price).
func fnPOLYLINE(args []*Value, _ []*types.MarketData) (*Value, error) {
	if len(args) != 2 {
		return nil, errors.NewRuntimeError("POLYLINE requires 2 arguments")
	}
	return buildPointDrawings("POLYLINE", args[0], args[1], nil, nil, "price")
}

// fnDRAWBAND implements DRAWBAND(upper, upperColor, lower, lowerColor).
func fnDRAWBAND(args []*Value, _ []*types.MarketData) (*Value, error) {
	if len(args) != 4 {
		return nil, errors.NewRuntimeError("DRAWBAND requires 4 arguments")
	}

	upper, upperColor, lower, lowerColor := args[0], args[1], args[2], args[3]
	length, err := drawingValueLength("DRAWBAND", upper, upperColor, lower, lowerColor)
	if err != nil {
		return nil, err
	}

	drawings := make([]*types.DrawingEvent, 0, length)
	for i := 0; i < length; i++ {
		drawings = append(drawings, &types.DrawingEvent{
			Function: "DRAWBAND",
			BarIndex: i,
			Values: map[string]float64{
				"upper":      scalarOrArrayAt(upper, i),
				"upperColor": scalarOrArrayAt(upperColor, i),
				"lower":      scalarOrArrayAt(lower, i),
				"lowerColor": scalarOrArrayAt(lowerColor, i),
			},
		})
	}

	return NewDrawingValue(drawings), nil
}

// fnDRAWKLINE implements DRAWKLINE(high, open, low, close).
func fnDRAWKLINE(args []*Value, _ []*types.MarketData) (*Value, error) {
	if len(args) != 4 {
		return nil, errors.NewRuntimeError("DRAWKLINE requires 4 arguments")
	}

	high, open, low, close := args[0], args[1], args[2], args[3]
	length, err := drawingValueLength("DRAWKLINE", high, open, low, close)
	if err != nil {
		return nil, err
	}

	drawings := make([]*types.DrawingEvent, 0, length)
	for i := 0; i < length; i++ {
		drawings = append(drawings, &types.DrawingEvent{
			Function: "DRAWKLINE",
			BarIndex: i,
			Values: map[string]float64{
				"high":  scalarOrArrayAt(high, i),
				"open":  scalarOrArrayAt(open, i),
				"low":   scalarOrArrayAt(low, i),
				"close": scalarOrArrayAt(close, i),
			},
		})
	}

	return NewDrawingValue(drawings), nil
}

func buildPointDrawings(function string, condition, price, numeric *Value, text *Value, priceKey string) (*Value, error) {
	if !condition.IsArray {
		return nil, errors.NewRuntimeError(fmt.Sprintf("%s first argument must be an array", function))
	}
	if err := validateDrawingNumericArgs(function, len(condition.Array), price); err != nil {
		return nil, err
	}
	if numeric != nil {
		if err := validateDrawingNumericArgs(function, len(condition.Array), numeric); err != nil {
			return nil, err
		}
	}

	valueCount := 1
	if numeric != nil {
		valueCount = 2
	}
	drawings := make([]*types.DrawingEvent, 0, truthyCount(condition.Array))
	for i, cond := range condition.Array {
		if !isTruthy(cond) {
			continue
		}

		event := &types.DrawingEvent{
			Function: function,
			BarIndex: i,
			Values:   make(map[string]float64, valueCount),
		}
		event.Values[priceKey] = scalarOrArrayAt(price, i)
		if numeric != nil {
			event.Values["value"] = scalarOrArrayAt(numeric, i)
		}
		if text != nil {
			event.Text = textValueAt(text, i)
		}
		drawings = append(drawings, event)
	}

	return NewDrawingValue(drawings), nil
}

func validateDrawingNumericArgs(function string, length int, values ...*Value) error {
	for _, value := range values {
		if value == nil || value.IsString || value.IsDraw {
			return errors.NewRuntimeError(fmt.Sprintf("%s arguments must be numeric", function))
		}
		if value.IsArray && len(value.Array) != length {
			return errors.NewRuntimeError(fmt.Sprintf("%s: array length mismatch", function))
		}
	}
	return nil
}

func drawingValueLength(function string, values ...*Value) (int, error) {
	length := 0
	for _, value := range values {
		if value == nil || value.IsString || value.IsDraw {
			return 0, errors.NewRuntimeError(fmt.Sprintf("%s arguments must be numeric", function))
		}
		if !value.IsArray {
			continue
		}
		if length == 0 {
			length = len(value.Array)
			continue
		}
		if len(value.Array) != length {
			return 0, errors.NewRuntimeError(fmt.Sprintf("%s: array length mismatch", function))
		}
	}
	if length == 0 {
		length = 1
	}
	return length, nil
}

func scalarOrArrayAt(value *Value, index int) float64 {
	if value == nil {
		return math.NaN()
	}
	if value.IsArray {
		if index >= len(value.Array) {
			return math.NaN()
		}
		return value.Array[index]
	}
	return value.Single
}

func textValueAt(value *Value, index int) string {
	if value == nil {
		return ""
	}
	if value.IsString {
		return value.Text
	}
	if value.IsArray {
		return fmt.Sprintf("%g", scalarOrArrayAt(value, index))
	}
	return fmt.Sprintf("%g", value.Single)
}

func numericUnaryFunc(args []*Value, name string, fn func(float64) float64) (*Value, error) {
	if len(args) != 1 {
		return nil, errors.NewRuntimeError(fmt.Sprintf("%s requires 1 argument", name))
	}
	value := args[0]
	if value.IsString || value.IsDraw {
		return nil, errors.NewRuntimeError(fmt.Sprintf("%s argument must be numeric", name))
	}
	if !value.IsArray {
		return NewSingleValue(fn(value.Single)), nil
	}

	result := make([]float64, len(value.Array))
	for i, v := range value.Array {
		result[i] = fn(v)
	}
	return NewArrayValue(result), nil
}

func numericBinaryFunc(args []*Value, name string, fn func(float64, float64) float64) (*Value, error) {
	if len(args) != 2 {
		return nil, errors.NewRuntimeError(fmt.Sprintf("%s requires 2 arguments", name))
	}

	a, b := args[0], args[1]
	if a.IsString || b.IsString || a.IsDraw || b.IsDraw {
		return nil, errors.NewRuntimeError(fmt.Sprintf("%s arguments must be numeric", name))
	}
	if !a.IsArray && !b.IsArray {
		return NewSingleValue(fn(a.Single, b.Single)), nil
	}

	length := valueLength(a)
	if length == 0 {
		length = valueLength(b)
	}
	if a.IsArray && b.IsArray && len(a.Array) != len(b.Array) {
		return nil, errors.NewRuntimeError(fmt.Sprintf("%s: array length mismatch", name))
	}

	result := make([]float64, length)
	for i := range result {
		result[i] = fn(scalarOrArrayAt(a, i), scalarOrArrayAt(b, i))
	}
	return NewArrayValue(result), nil
}

func rollingStatsFunc(args []*Value, name string, fn func([]float64) float64) (*Value, error) {
	if len(args) != 2 {
		return nil, errors.NewRuntimeError(fmt.Sprintf("%s requires 2 arguments", name))
	}

	data := args[0]
	period := args[1]
	if !data.IsArray {
		return nil, errors.NewRuntimeError(fmt.Sprintf("%s first argument must be an array", name))
	}
	if period.IsArray {
		return nil, errors.NewRuntimeError(fmt.Sprintf("%s second argument must be a number", name))
	}

	n := int(period.Single)
	if n <= 0 || n > len(data.Array) {
		return nil, errors.NewRuntimeError(fmt.Sprintf("%s period must be between 1 and %d", name, len(data.Array)))
	}

	result := make([]float64, len(data.Array))
	for i := 0; i < n-1; i++ {
		result[i] = math.NaN()
	}
	window := make([]float64, n)
	for i := n - 1; i < len(data.Array); i++ {
		copy(window, data.Array[i-n+1:i+1])
		result[i] = fn(window)
	}
	return NewArrayValue(result), nil
}

func rollingRegressionFunc(args []*Value, name string, fn func(float64, float64, int) float64) (*Value, error) {
	return rollingStatsFunc(args, name, func(values []float64) float64 {
		slope, intercept := linearRegression(values)
		return fn(slope, intercept, len(values))
	})
}

func rollingPairStatsFunc(args []*Value, name string, fn func([]float64, []float64) float64) (*Value, error) {
	if len(args) != 3 {
		return nil, errors.NewRuntimeError(fmt.Sprintf("%s requires 3 arguments", name))
	}

	a, b, period := args[0], args[1], args[2]
	if !a.IsArray || !b.IsArray {
		return nil, errors.NewRuntimeError(fmt.Sprintf("%s first two arguments must be arrays", name))
	}
	if len(a.Array) != len(b.Array) {
		return nil, errors.NewRuntimeError(fmt.Sprintf("%s: array length mismatch", name))
	}
	if period.IsArray {
		return nil, errors.NewRuntimeError(fmt.Sprintf("%s third argument must be a number", name))
	}

	n := int(period.Single)
	if n <= 0 || n > len(a.Array) {
		return nil, errors.NewRuntimeError(fmt.Sprintf("%s period must be between 1 and %d", name, len(a.Array)))
	}

	result := make([]float64, len(a.Array))
	for i := 0; i < n-1; i++ {
		result[i] = math.NaN()
	}
	windowA := make([]float64, n)
	windowB := make([]float64, n)
	for i := n - 1; i < len(a.Array); i++ {
		copy(windowA, a.Array[i-n+1:i+1])
		copy(windowB, b.Array[i-n+1:i+1])
		result[i] = fn(windowA, windowB)
	}
	return NewArrayValue(result), nil
}

func futureReference(args []*Value, name string) (*Value, error) {
	if len(args) != 2 {
		return nil, errors.NewRuntimeError(fmt.Sprintf("%s requires 2 arguments", name))
	}

	data := args[0]
	period := args[1]
	if !data.IsArray {
		return nil, errors.NewRuntimeError(fmt.Sprintf("%s first argument must be an array", name))
	}
	if period.IsArray {
		return nil, errors.NewRuntimeError(fmt.Sprintf("%s second argument must be a number", name))
	}

	n := int(period.Single)
	if n < 0 {
		return nil, errors.NewRuntimeError(fmt.Sprintf("%s period must be non-negative", name))
	}

	result := make([]float64, len(data.Array))
	for i := range data.Array {
		futureIndex := i + n
		if futureIndex >= len(data.Array) {
			result[i] = math.NaN()
			continue
		}
		result[i] = data.Array[futureIndex]
	}
	return NewArrayValue(result), nil
}

func average(values []float64) float64 {
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}

func variance(values []float64) float64 {
	mean := average(values)
	sum := 0.0
	for _, v := range values {
		diff := v - mean
		sum += diff * diff
	}
	return sum / float64(len(values))
}

func covariance(a, b []float64) float64 {
	meanA := average(a)
	meanB := average(b)
	sum := 0.0
	for i := range a {
		sum += (a[i] - meanA) * (b[i] - meanB)
	}
	return sum / float64(len(a))
}

func linearRegression(values []float64) (float64, float64) {
	n := float64(len(values))
	sumX := 0.0
	sumY := 0.0
	sumXY := 0.0
	sumXX := 0.0
	for i, y := range values {
		x := float64(i)
		sumX += x
		sumY += y
		sumXY += x * y
		sumXX += x * x
	}
	denominator := n*sumXX - sumX*sumX
	if denominator == 0 {
		return 0, values[len(values)-1]
	}
	slope := (n*sumXY - sumX*sumY) / denominator
	intercept := (sumY - slope*sumX) / n
	return slope, intercept
}

func valueLength(value *Value) int {
	if value != nil && value.IsArray {
		return len(value.Array)
	}
	return 0
}

func compareConsecutive(args []*Value, name string, cmp func(curr, prev float64) bool) (*Value, error) {
	data := args[0]
	period := args[1]

	if !data.IsArray {
		return nil, errors.NewRuntimeError(fmt.Sprintf("%s first argument must be an array", name))
	}
	if period.IsArray {
		return nil, errors.NewRuntimeError(fmt.Sprintf("%s second argument must be a number", name))
	}

	n := int(period.Single)
	if n <= 0 {
		return nil, errors.NewRuntimeError(fmt.Sprintf("%s period must be positive", name))
	}

	result := make([]float64, len(data.Array))
	for i := n; i < len(data.Array); i++ {
		ok := true
		for j := 0; j < n; j++ {
			if !cmp(data.Array[i-j], data.Array[i-j-1]) {
				ok = false
				break
			}
		}
		if ok {
			result[i] = 1
		}
	}

	return NewArrayValue(result), nil
}

func isTruthy(value float64) bool {
	return value != 0 && !math.IsNaN(value)
}

func truthyCount(values []float64) int {
	count := 0
	for _, value := range values {
		if isTruthy(value) {
			count++
		}
	}
	return count
}
