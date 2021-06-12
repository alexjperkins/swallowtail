package client

import (
	"context"
	"sync"
	"testing"
)

// mockClient is a mock client for binance.
type mockClient struct {
	called map[string]int
	sync.Mutex
}

func (m *mockClient) ListAllAssetPairs(context.Context) (*ListAllAssetPairsResponse, error) {
	m.Lock()
	defer m.Unlock()
	if m.called == nil {
		m.called = map[string]int{}
	}
	m.called["ListAllAssetPairs"]++
	return &ListAllAssetPairsResponse{
		Symbols: []*BinanceAssetItem{
			{
				Symbol:            "BTCUSDT",
				BaseAsset:         "BTC",
				WithMarginTrading: true,
				WithSpotTrading:   true,
			},
			{
				Symbol:            "ETHUSDT",
				BaseAsset:         "ETH",
				WithMarginTrading: true,
				WithSpotTrading:   true,
			},
			{
				Symbol:            "SOLUSDT",
				BaseAsset:         "SOL",
				WithMarginTrading: true,
				WithSpotTrading:   true,
			},
			{
				Symbol:            "USDTUSD",
				BaseAsset:         "USDT",
				WithMarginTrading: true,
				WithSpotTrading:   true,
			},
		},
	}, nil
}

func (m *mockClient) Ping(ctx context.Context) error {
	return nil
}

func (m *mockClient) AssertListAllAssetPairs(t *testing.T, expectedNumberOfCalls int) {
	m.Lock()
	defer m.Unlock()
	howManyCalls, ok := m.called["ListAllAssetPairs"]
	if !ok {
		t.Fatalf("ListAllAssetPairs not found in mock called")
	}
	if howManyCalls != expectedNumberOfCalls {
		t.Errorf("ListAllAssetPairs: Expecting %v calls, got %v", expectedNumberOfCalls, howManyCalls)
	}
}
