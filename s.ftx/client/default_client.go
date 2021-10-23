package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"swallowtail/libraries/gerrors"
	"swallowtail/libraries/transport"
	"time"
)

type ftxClient struct {
	http     transport.HttpClient
	hostname string
}

func (f *ftxClient) Ping(ctx context.Context) error {
	err := f.do(ctx, http.MethodGet, "/api/stats/latency_stats", PingRequest{}, PingResponse{}, nil, nil)
	switch {
	case gerrors.Is(err, gerrors.ErrUnauthenticated):
		return nil
	case err != nil:
		return gerrors.Augment(err, "failed_to_establish_ftx_connection", nil)
	default:
		return nil
	}
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

func (f *ftxClient) do(ctx context.Context, method, endpoint string, req, rsp interface{}, pagination *PaginationFilter, credentials *Credentials) error {
	url := fmt.Sprintf("%s%s", f.hostname, buildEndpoint(endpoint, pagination))

	var creds = credentials
	if creds == nil {
		creds = &Credentials{}
	}

	return f.http.DoWithEphemeralHeaders(ctx, method, url, req, rsp, creds.SubaccountAsHeaders())
}

func (f *ftxClient) signBeforeDo(ctx context.Context, method, endpoint string, req, rsp interface{}, pagination *PaginationFilter, credentials *Credentials) error {
	ts := strconv.FormatInt(time.Now().UTC().Unix()*1000, 10)
	preparedEndpoint := buildEndpoint(endpoint, pagination)

	reqBody, err := json.Marshal(req)
	if err != nil {
		return gerrors.Augment(err, "failed_to_sign_request.bad_request_body", nil)
	}

	signature, err := credentials.SignRequest(method, preparedEndpoint, ts, reqBody)
	if err != nil {
		return gerrors.Augment(err, "failed_to_send_request.failed_to_sign_request", nil)
	}

	url := fmt.Sprintf("%s%s", f.hostname, preparedEndpoint)

	return f.http.DoWithEphemeralHeaders(ctx, method, url, req, rsp, credentials.AsHeaders(signature, ts))
}

func buildEndpoint(base string, pagination *PaginationFilter) string {
	var endpoint = base
	if pagination != nil {
		endpoint = fmt.Sprintf("%s?%s", endpoint, pagination.ToQueryString())
	}

	return endpoint
}
