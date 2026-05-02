package helpers

import (
	"fmt"
	"time"

	"github.com/DTrader-store/formula-go/types"
	"github.com/injoyai/tdx"
	"github.com/injoyai/tdx/protocol"
)

// TDXClient wraps the TDX connection for easy data fetching
type TDXClient struct {
	client *tdx.Client
}

// DailyKlines contains formula-ready market data with the original K-line timestamps.
type DailyKlines struct {
	Data  []*types.MarketData
	Times []time.Time
}

// NewTDXClient creates a new TDX client with auto-reconnect
func NewTDXClient(addr string) (*TDXClient, error) {
	client, err := tdx.Dial(addr, tdx.WithRedial())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to TDX server: %w", err)
	}
	return &TDXClient{client: client}, nil
}

// Close closes the TDX connection
func (c *TDXClient) Close() error {
	return c.client.Close()
}

// GetMarketData fetches K-line data and converts to MarketData format
func (c *TDXClient) GetMarketData(code string, start uint16, count uint16) ([]*types.MarketData, error) {
	dailyKlines, err := c.GetDailyKlines(code, start, count)
	if err != nil {
		return nil, err
	}
	return dailyKlines.Data, nil
}

// GetDailyKlines fetches daily K-line data and keeps timestamps aligned with MarketData indexes.
func (c *TDXClient) GetDailyKlines(code string, start uint16, count uint16) (*DailyKlines, error) {
	// Fetch K-line data from TDX
	klines, err := c.client.GetKlineDay(code, start, count)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch daily K-line data: %w", err)
	}

	// Convert to MarketData format
	marketData := make([]*types.MarketData, len(klines.List))
	times := make([]time.Time, len(klines.List))
	for i, kline := range klines.List {
		marketData[i] = KlineToMarketData(kline)
		times[i] = kline.Time
	}

	return &DailyKlines{Data: marketData, Times: times}, nil
}

// KlineToMarketData converts TDX Kline to formula MarketData
func KlineToMarketData(kline *protocol.Kline) *types.MarketData {
	return &types.MarketData{
		Open:   kline.Open.Float64(),
		High:   kline.High.Float64(),
		Low:    kline.Low.Float64(),
		Close:  kline.Close.Float64(),
		Volume: float64(kline.Volume),
		Amount: kline.Amount.Float64(),
	}
}

// DefaultTDXServer returns a reliable TDX server address
func DefaultTDXServer() string {
	return "124.71.187.122:7709" // Shanghai Huawei server
}
