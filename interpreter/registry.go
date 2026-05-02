package interpreter

import (
	"fmt"
	"sort"
	"strings"

	"github.com/DTrader-store/formula-go/errors"
	"github.com/DTrader-store/formula-go/types"
)

// Function represents a built-in function
type Function func(args []*Value, marketData []*types.MarketData) (*Value, error)

// FunctionRegistry manages built-in functions
type FunctionRegistry struct {
	functions map[string]Function
}

// NewFunctionRegistry creates a new function registry
func NewFunctionRegistry() *FunctionRegistry {
	reg := &FunctionRegistry{
		functions: make(map[string]Function),
	}
	reg.registerBuiltinFunctions()
	return reg
}

// Register registers a function
func (r *FunctionRegistry) Register(name string, fn Function) {
	r.functions[strings.ToUpper(name)] = fn
}

// Names returns all registered function names in sorted uppercase order.
func (r *FunctionRegistry) Names() []string {
	names := make([]string, 0, len(r.functions))
	for name := range r.functions {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// Call calls a registered function
func (r *FunctionRegistry) Call(name string, args []*Value, marketData []*types.MarketData) (*Value, error) {
	fn, exists := r.functions[strings.ToUpper(name)]
	if !exists {
		return nil, errors.NewRuntimeError(fmt.Sprintf("undefined function: %s", name))
	}
	return fn(args, marketData)
}

// registerBuiltinFunctions registers all built-in functions
func (r *FunctionRegistry) registerBuiltinFunctions() {
	// Mathematical functions
	r.Register("MA", fnMA)
	r.Register("EMA", fnEMA)
	r.Register("SUM", fnSUM)
	r.Register("MAX", fnMAX)
	r.Register("MIN", fnMIN)
	r.Register("ABS", fnABS)
	r.Register("SQRT", fnSQRT)
	r.Register("POW", fnPOW)
	r.Register("EXP", fnEXP)
	r.Register("LN", fnLN)
	r.Register("LOG", fnLOG)
	r.Register("MOD", fnMOD)
	r.Register("CEILING", fnCEILING)
	r.Register("FLOOR", fnFLOOR)
	r.Register("INTPART", fnINTPART)
	r.Register("FRACPART", fnFRACPART)
	r.Register("ROUND", fnROUND)
	r.Register("ROUND2", fnROUND2)
	r.Register("SIGN", fnSIGN)
	r.Register("SIN", fnSIN)
	r.Register("COS", fnCOS)
	r.Register("TAN", fnTAN)
	r.Register("ASIN", fnASIN)
	r.Register("ACOS", fnACOS)
	r.Register("ATAN", fnATAN)

	// Reference functions
	r.Register("REF", fnREF)
	r.Register("REFV", fnREFV)
	r.Register("REFX", fnREFX)
	r.Register("REFXV", fnREFXV)
	r.Register("HHV", fnHHV)
	r.Register("LLV", fnLLV)

	// Conditional functions
	r.Register("IF", fnIF)
	r.Register("CROSS", fnCROSS)

	// Phase 4: Additional functions
	r.Register("STD", fnSTD)
	r.Register("STDP", fnSTDP)
	r.Register("STDDEV", fnSTDDEV)
	r.Register("VAR", fnVAR)
	r.Register("VARP", fnVARP)
	r.Register("DEVSQ", fnDEVSQ)
	r.Register("FORCAST", fnFORCAST)
	r.Register("SLOPE", fnSLOPE)
	r.Register("COVAR", fnCOVAR)
	r.Register("RELATE", fnRELATE)
	r.Register("BETA", fnBETA)
	r.Register("SMA", fnSMA)
	r.Register("WMA", fnWMA)
	r.Register("COUNT", fnCOUNT)
	r.Register("EVERY", fnEVERY)
	r.Register("EXIST", fnEXIST)
	r.Register("BARSLAST", fnBARSLAST)
	r.Register("HHVBARS", fnHHVBARS)
	r.Register("LLVBARS", fnLLVBARS)
	r.Register("BARSCOUNT", fnBARSCOUNT)
	r.Register("BARSSINCE", fnBARSSINCE)
	r.Register("BARSLASTCOUNT", fnBARSLASTCOUNT)
	r.Register("CURRBARSCOUNT", fnCURRBARSCOUNT)
	r.Register("TOTALBARSCOUNT", fnTOTALBARSCOUNT)
	r.Register("ISLASTBAR", fnISLASTBAR)
	r.Register("BARSTATUS", fnBARSTATUS)
	r.Register("SUMBARS", fnSUMBARS)
	r.Register("AVEDEV", fnAVEDEV)
	r.Register("FILTER", fnFILTER)
	r.Register("BETWEEN", fnBETWEEN)
	r.Register("RANGE", fnRANGE)
	r.Register("NOT", fnNOT)
	r.Register("IFN", fnIFN)
	r.Register("IFF", fnIFF)
	r.Register("DRAWNULL", fnDRAWNULL)
	r.Register("LONGCROSS", fnLONGCROSS)
	r.Register("UPNDAY", fnUPNDAY)
	r.Register("DOWNNDAY", fnDOWNNDAY)
	r.Register("NDAY", fnNDAY)
	r.Register("LAST", fnLAST)
	r.Register("EXISTR", fnEXISTR)
	r.Register("DMA", fnDMA)
	r.Register("CONST", fnCONST)
	r.Register("VALUEWHEN", fnVALUEWHEN)

	// Drawing functions
	r.Register("DRAWTEXT", fnDRAWTEXT)
	r.Register("DRAWICON", fnDRAWICON)
	r.Register("DRAWNUMBER", fnDRAWNUMBER)
	r.Register("STICKLINE", fnSTICKLINE)
	r.Register("DRAWLINE", fnDRAWLINE)
	r.Register("POLYLINE", fnPOLYLINE)
	r.Register("DRAWBAND", fnDRAWBAND)
	r.Register("DRAWKLINE", fnDRAWKLINE)
}
