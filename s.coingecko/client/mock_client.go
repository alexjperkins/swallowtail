package client

import "context"

type MockCoingeckoClient struct{}

func (m *MockCoingeckoClient) Ping(ctx context.Context) error {
	return nil
}

func (m *MockCoingeckoClient) GetCurrentPriceFromSymbol(ctx context.Context, symbol, assetPair string) (float64, error) {
	return 0, nil
}

func (m *MockCoingeckoClient) GetCurrentPriceFromID(ctx context.Context, id, assetPair string) (float64, error) {
	return 0, nil
}

func (m *MockCoingeckoClient) GetATHFromSymbol(ctx context.Context, symbol string) (float64, error) {
	return 0, nil
}

func (m *MockCoingeckoClient) GetATHFromID(ctx context.Context, id string) (float64, error) {
	return 0, nil
}

func (m *MockCoingeckoClient) RefreshCoins(ctx context.Context) {
}
