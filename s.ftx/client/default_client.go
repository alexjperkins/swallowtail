package client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/monzo/slog"

	"swallowtail/libraries/gerrors"
	"swallowtail/libraries/transport"
)

const (
	// FTX API Version
	APIVersion = "/api/"
)

type ftxClient struct {
	http     transport.HttpClient
	hostname string
}

func (f *ftxClient) Ping(ctx context.Context) error {
	err := f.do(ctx, http.MethodGet, "/api/stats/latency_stats", nil, nil, nil, nil)
	switch {
	case gerrors.Is(err, gerrors.ErrUnauthenticated):
		return nil
	case err != nil:
		return gerrors.Augment(err, "failed_to_establish_ftx_connection", nil)
	default:
		return nil
	}
}

func (f *ftxClient) GetStatus(ctx context.Context, req *GetStatusRequest) (*GetStatusResponse, error) {
	rsp := &GetStatusResponse{}
	if err := f.signBeforeDo(ctx, http.MethodGet, "/api/stats/latency_stats", req, rsp, nil, nil); err != nil {
		return nil, gerrors.Augment(err, "failed_get_ftx_status", nil)
	}

	return rsp, nil
}

func (f *ftxClient) ExecuteOrder(ctx context.Context, req *ExecuteOrderRequest, credentials *Credentials) (*ExecuteOrderResponse, error) {
	var endpoint = "orders"
	switch req.Type {
	case "stop", "trailingStop", "takeProfit":
		endpoint = "condition_order"
	}

	rsp := &ExecuteOrderResponse{}
	if err := f.signBeforeDo(ctx, http.MethodPost, fmt.Sprintf("%s%s", APIVersion, endpoint), req, rsp, nil, credentials); err != nil {
		slog.Warn(ctx, "FTX order failed: %+v %v", err, req)
		return nil, gerrors.Augment(err, "failed_to_post_order", nil)
	}

	slog.Info(ctx, "FTX order executed: %v")

	return rsp, nil
}

func (f *ftxClient) ListAccountDeposits(ctx context.Context, req *ListAccountDepositsRequest, pagination *PaginationFilter) (*ListAccountDepositsResponse, error) {
	rsp := &ListAccountDepositsResponse{}
	if err := f.signBeforeDo(ctx, http.MethodGet, "/api/wallet/deposits", req, rsp, pagination, depositAccountCredentials); err != nil {
		return nil, gerrors.Augment(err, "failed_to_list_account_deposits", map[string]string{
			"subaccount": depositAccountCredentials.Subaccount,
		})
	}

	return rsp, nil
}

func (f *ftxClient) VerifyCredentials(ctx context.Context, req *VerifyCredentialsRequest, credentials *Credentials) (*VerifyCredentialsResponse, error) {
	rsp := &VerifyCredentialsResponse{}
	if err := f.signBeforeDo(ctx, http.MethodGet, "/api/account", req, rsp, nil, credentials); err != nil {
		return nil, gerrors.Augment(err, "failed_to_verify_credentials", map[string]string{
			"subaccount": credentials.Subaccount,
		})
	}

	return rsp, nil
}

func (f *ftxClient) GetFundingRate(ctx context.Context, req *GetFundingRateRequest) (*GetFundingRateResponse, error) {
	var pagination *PaginationFilter
	switch {
	case req.StartTime != 0, req.EndTime != 0:
		pagination = &PaginationFilter{
			Start: int64(req.StartTime),
			End:   int64(req.EndTime),
		}
	}

	rsp := &GetFundingRateResponse{}
	if err := f.do(ctx, http.MethodGet, fmt.Sprintf("/api/funding_rates?future=%s", req.Instrument), req, rsp, pagination, nil); err != nil {
		return nil, gerrors.Augment(err, "failed_to_get_funding_rate", nil)
	}

	return rsp, nil
}

func (f *ftxClient) ListInstruments(ctx context.Context, req *ListInstrumentsRequest, futuresOnly bool) (*ListInstrumentsResponse, error) {
	// Determine the correct endpoint based on whether the caller requires `futuresOnly`.
	var (
		endpoint string
		rsp      interface{}
	)
	switch {
	case futuresOnly:
		endpoint = "futures"
		rsp = &ListFuturesInstrumentsResponse{}
	default:
		endpoint = "markets"
		rsp = &ListMarketsInstrumentsResponse{}
	}

	if err := f.signBeforeDo(ctx, http.MethodGet, fmt.Sprintf("%s%s", APIVersion, endpoint), req, rsp, nil, nil); err != nil {
		return nil, gerrors.Augment(err, "failed_to_list_instruments", nil)
	}

	// Marshal into generic response.
	switch {
	case futuresOnly:
		return &ListInstrumentsResponse{}, nil
	default:
		return &ListInstrumentsResponse{}, nil
	}
}
