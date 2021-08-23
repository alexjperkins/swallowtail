package client

import (
	"context"
	"fmt"
	"net/http"
	"swallowtail/libraries/gerrors"
	"swallowtail/libraries/transport"
	"swallowtail/s.binance/domain"
	"time"

	"github.com/monzo/terrors"
)

type defaultClient struct {
	c transport.HttpClient
}

// NewDefaultClient ...
func NewDefaultClient(ctx context.Context) (*defaultClient, error) {
	c := &defaultClient{
		transport.NewHTTPClient(ctx, time.Duration(30*time.Second)),
	}
	return c, c.Ping(ctx)
}

func (c *defaultClient) ListAllAssetPairs(ctx context.Context) (*ListAllAssetPairsResponse, error) {
	url := fmt.Sprintf("%s/%s", binanceAPIUrl, "exchangeInfo")
	rspBody := &ListAllAssetPairsResponse{}
	if err := c.c.Do(ctx, http.MethodGet, url, nil, rspBody); err != nil {
		return nil, terrors.Augment(err, "Failed to list all asset pairs", nil)
	}
	return rspBody, nil
}

func (c *defaultClient) ExecuteSpotTrade(ctx context.Context, trade *domain.Trade) error {
	return nil
}

func (c *defaultClient) ReadSpotAccount(ctx context.Context, in *ReadSpotAccountRequest) (*ReadSpotAccountResponse, error) {
	return nil, nil
}

func (c *defaultClient) ReadPerpetualFuturesAccount(ctx context.Context, in *ReadPerptualFuturesAccountRequest) (*ReadPerptualFuturesAccountResponse, error) {
	return nil, nil
}

func (c *defaultClient) Ping(ctx context.Context) error {
	endpoint := fmt.Sprintf("%s/ping", binanceAPIUrl)
	rspBody := &PingResponse{}
	if err := c.c.Do(ctx, http.MethodGet, endpoint, nil, rspBody); err != nil {
		return terrors.Augment(err, "Failed to connect to the Binance API.", nil)
	}
	return nil
}

func (c *defaultClient) VerifyCredentials(ctx context.Context, credentials *Credentials) (*VerifyCredentialsResponse, error) {
	endpoint := fmt.Sprintf("%s/sapi/v1/account/apiRestrictions", binanceAPIUrl)
	rspBody := &VerifyCredentialsResponse{}
	if err := c.c.DoWithEphemeralHeaders(ctx, http.MethodGet, endpoint, nil, rspBody, credentials.AsHeaders()); err != nil {
		return nil, gerrors.Augment(err, "client_request_failed.verify_credentials", map[string]string{
			"endpoint": endpoint,
		})
	}
	return rspBody, nil
}
