package client

import "context"

type cgClient struct{}

func (c *cgClient) GetAllCoinIDs(ctx context.Context) ([]*CoingeckoListCoinItem, error) {
	return nil, nil
}

func (c *cgClient) GetCurrentPriceFromSymbol(ctx context.Context, symbol, assetPair string) (float64, error) {
	return 0, nil
}

func (c *cgClient) GetCurrentPriceFromID(ctx context.Context, id, assetPair string) (float64, error) {
	return 0, nil
}

func (c *cgClient) GetATHFromSymbol(ctx context.Context, symbol string) (float64, error) {
	return 0, nil
}

func (c *cgClient) GetATHFromID(ctx context.Context, id string) (float64, error) {
	return 0, nil
}

func (c *cgClient) Ping(ctx context.Context) bool {
	return false
}
