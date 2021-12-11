package client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/monzo/slog"
	"github.com/monzo/terrors"

	"swallowtail/libraries/gerrors"
	"swallowtail/libraries/transport"
	"swallowtail/s.binance/domain"
)

type binanceClient struct {
	http transport.HttpClient
}

func (c *binanceClient) GetLatestPrice(ctx context.Context, req *GetLatestPriceRequest) (*GetLatestPriceResponse, error) {
	url := fmt.Sprintf("%s/%s", binanceFuturesURL, "ticker/price")
	rspBody := &GetLatestPriceResponse{}

	qs := fmt.Sprintf("symbol=%s", req.Symbol)

	if err := c.do(ctx, http.MethodGet, url, qs, nil, rspBody, defaultCredentials); err != nil {
		return nil, gerrors.Augment(err, "failed_to_get_latest_price", nil)
	}

	return rspBody, nil
}

func (c *binanceClient) ListAllAssetPairs(ctx context.Context) (*ListAllAssetPairsResponse, error) {
	url := fmt.Sprintf("%s/%s", binanceAPIUrl, "exchangeInfo")
	rspBody := &ListAllAssetPairsResponse{}

	if err := c.http.Do(ctx, http.MethodGet, url, nil, rspBody); err != nil {
		return nil, terrors.Augment(err, "Failed to list all asset pairs", nil)
	}

	return rspBody, nil
}

func (c *binanceClient) ExecuteSpotOrder(ctx context.Context, trade *domain.Trade) error {
	return gerrors.Unimplemented("unimplemented.execute_spot_trade", nil)
}

func (c *binanceClient) ReadSpotAccount(ctx context.Context, in *ReadSpotAccountRequest) (*ReadSpotAccountResponse, error) {
	return nil, gerrors.Unimplemented("unimplemented.read_spot_account", nil)
}

func (c *binanceClient) ReadPerpetualFuturesAccount(ctx context.Context, _ *ReadPerpetualFuturesAccountRequest, credentials *Credentials) (*ReadPerpetualFuturesAccountResponse, error) {
	url := fmt.Sprintf("%s/%s", binanceFuturesURLV2, "balance")
	rspBody := &ReadPerpetualFuturesAccountResponse{}

	if err := c.doWithSignature(ctx, http.MethodGet, url, "", nil, rspBody, credentials); err != nil {
		return nil, gerrors.Augment(err, "failed_to_read_perpetual_futures_account.client", nil)
	}

	return rspBody, nil
}

func (c *binanceClient) ExecutePerpetualFuturesOrder(ctx context.Context, req *ExecutePerpetualFuturesOrderRequest, credentials *Credentials) (*ExecutePerpetualFuturesOrderResponse, error) {
	url := fmt.Sprintf("%s/%s", binanceFuturesURL, "order")
	rspBody := &ExecutePerpetualFuturesOrderResponse{}

	qs := buildQueryStringFromFuturesPerpetualTrade(req)

	// Execute request.
	if err := c.doWithSignature(ctx, http.MethodPost, url, qs, nil, rspBody, credentials); err != nil {
		slog.Warn(ctx, "Binance Perpetuals futures trade failed: %v", qs)
		return nil, gerrors.Augment(err, "failed_to_execute_perpetual_futures_trade.client", map[string]string{
			"query_string": qs,
		})
	}

	return rspBody, nil
}

func (c *binanceClient) Ping(ctx context.Context) error {
	endpoint := fmt.Sprintf("%s/ping", binanceAPIUrl)
	rspBody := &PingResponse{}

	if err := c.http.Do(ctx, http.MethodGet, endpoint, nil, rspBody); err != nil {
		return terrors.Augment(err, "Failed to connect to the Binance API.", nil)
	}

	return nil
}

func (c *binanceClient) VerifyCredentials(ctx context.Context, credentials *Credentials) (*VerifyCredentialsResponse, error) {
	endpoint := fmt.Sprintf("%s/account/apiRestrictions", binanceSpotURL)
	rspBody := &VerifyCredentialsResponse{}

	if err := c.doWithSignature(ctx, http.MethodGet, endpoint, "", nil, rspBody, credentials); err != nil {
		return nil, gerrors.Augment(err, "client_request_failed.verify_credentials", map[string]string{
			"endpoint": endpoint,
		})
	}

	return rspBody, nil
}

func (c *binanceClient) GetFuturesExchangeInfo(ctx context.Context, req *GetFuturesExchangeInfoRequest) (*GetFuturesExchangeInfoResponse, error) {
	endpoint := fmt.Sprintf("%s/exchangeInfo", binanceFuturesURL)
	rspBody := &GetFuturesExchangeInfoResponse{}

	if err := c.do(ctx, http.MethodGet, endpoint, "", nil, rspBody, nil); err != nil {
		return nil, gerrors.Augment(err, "client_request_failed.verify_credentials", map[string]string{
			"endpoint": endpoint,
		})
	}

	return rspBody, nil
}

func (c *binanceClient) GetStatus(ctx context.Context) (*GetStatusResponse, error) {
	endpoint := fmt.Sprintf("%s/time", binanceAPIUrl)

	rspBody := &GetStatusResponse{}
	if err := c.http.Do(ctx, http.MethodGet, endpoint, nil, rspBody); err != nil {
		return nil, gerrors.Augment(err, "client_request_failed.get_status.time", map[string]string{
			"endpoint": endpoint,
		})
	}

	// Convert to millisecond time
	rspBody.ServerTime /= 1_000

	return rspBody, nil
}

func (c *binanceClient) GetFundingRate(ctx context.Context, req *GetFundingRateRequest) (*GetFundingRateResponse, error) {
	endpoint := fmt.Sprintf("%s/fundingRate?symbol=%s", binanceFuturesURL, req.Symbol)

	if req.StartTime != 0 {
		endpoint = fmt.Sprintf("%s&startTime=%d", endpoint, req.StartTime)
	}
	if req.EndTime != 0 {
		endpoint = fmt.Sprintf("%s&endTime=%d", endpoint, req.EndTime)
	}

	if req.Limit != 0 {
		endpoint = fmt.Sprintf("%s&limit=%d", endpoint, req.Limit)
	}

	rspBody := &GetFundingRateResponse{}
	if err := c.http.Do(ctx, http.MethodGet, endpoint, nil, rspBody); err != nil {
		return nil, gerrors.Augment(err, "client_request_failed.get_funding_rates", nil)
	}

	return rspBody, nil
}
