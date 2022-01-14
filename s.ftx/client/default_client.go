package client

import (
	"context"
	"fmt"
	"net/http"

	"swallowtail/libraries/gerrors"
	"swallowtail/libraries/transport"
	"swallowtail/s.ftx/client/auth"
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
	if err := f.do(ctx, http.MethodGet, "/api/stats/latency_stats", req, rsp, nil, nil); err != nil {
		return nil, gerrors.Augment(err, "failed_get_ftx_status", nil)
	}

	return rsp, nil
}

func (f *ftxClient) ExecuteOrder(ctx context.Context, req *ExecuteOrderRequest, credentials *auth.Credentials) (*ExecuteOrderResponse, error) {
	var endpoint = "orders"
	switch req.Type {
	case "stop", "trailingStop", "takeProfit":
		endpoint = "conditional_orders"
	}

	rsp := &ExecuteOrderResponse{}
	if err := f.signBeforeDo(ctx, http.MethodPost, fmt.Sprintf("/api/%s", endpoint), req, rsp, nil, credentials); err != nil {
		return nil, gerrors.Augment(err, "failed_to_post_order", nil)
	}

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

func (f *ftxClient) VerifyCredentials(ctx context.Context, req *VerifyCredentialsRequest, credentials *auth.Credentials) (*VerifyCredentialsResponse, error) {
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
	rsp := &ListInstrumentsResponse{}
	if err := f.do(ctx, http.MethodGet, "/api/markets", req, rsp, nil, nil); err != nil {
		return nil, gerrors.Augment(err, "failed_to_list_instruments", nil)
	}

	return rsp, nil
}

func (f *ftxClient) ReadAccountInformation(ctx context.Context, credentials *auth.Credentials) (*ReadAccountInformationResponse, error) {
	rsp := &ReadAccountInformationResponse{}
	if err := f.signBeforeDo(ctx, http.MethodGet, "/api/account", nil, rsp, nil, credentials); err != nil {
		return nil, gerrors.Augment(err, "failed_to_read_account_information", nil)
	}

	return rsp, nil
}

func (f *ftxClient) ListAccountBalances(ctx context.Context, credentials *auth.Credentials) (*ListAccountBalancesResponse, error) {
	rsp := &ListAccountBalancesResponse{}
	if err := f.signBeforeDo(ctx, http.MethodGet, "/api/wallet/balances", nil, rsp, nil, credentials); err != nil {
		return nil, gerrors.Augment(err, "failed_to_list_account_balances", nil)
	}

	return rsp, nil
}
