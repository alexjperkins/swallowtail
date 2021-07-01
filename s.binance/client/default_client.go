package client

import (
	"context"
	"fmt"
	"net/http"
	"swallowtail/libraries/transport"
	"swallowtail/s.binance/domain"
	"time"

	"github.com/monzo/terrors"
)

type defaultClient struct {
	c transport.HttpClient
}

func NewDefaultClient(ctx context.Context, apiKey string) (*defaultClient, error) {
	headers := map[string]string{
		"X-MBX-APIKEY": apiKey,
	}
	c := &defaultClient{
		transport.NewHTTPClient(ctx, time.Duration(30*time.Second), headers),
	}
	return c, c.Ping(ctx)
}

func (d *defaultClient) ListAllAssetPairs(ctx context.Context) (*ListAllAssetPairsResponse, error) {
	url := fmt.Sprintf("%s/%s", binanceAPIUrl, "exchangeInfo")
	rspBody := &ListAllAssetPairsResponse{}
	if err := d.c.DoRequest(ctx, http.MethodGet, url, nil, rspBody); err != nil {
		return nil, terrors.Augment(err, "Failed to list all asset pairs", nil)
	}
	return rspBody, nil
}

func (d *defaultClient) ExecuteSpotTrade(ctx context.Context, trade *domain.Trade) error {
	return nil
}

func (d *defaultClient) Ping(ctx context.Context) error {
	url := fmt.Sprintf("%s/%s", binanceAPIUrl, "ping")
	rspBody := &PingRequest{}
	if err := d.c.DoRequest(ctx, http.MethodGet, url, nil, rspBody); err != nil {
		return terrors.Augment(err, "Failed to connect to the Binance API.", nil)
	}
	return nil
}
