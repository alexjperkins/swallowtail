package client

import (
	"context"
	"fmt"
	"strconv"
	"swallowtail/libraries/gerrors"
	"swallowtail/s.ftx/client/auth"
	"time"
)

func (f *ftxClient) do(ctx context.Context, method, endpoint string, req, rsp interface{}, pagination *PaginationFilter, credentials *auth.Credentials) error {
	url := fmt.Sprintf("%s%s", f.hostname, buildEndpoint(endpoint, pagination))

	var creds = credentials
	if creds == nil {
		creds = &auth.Credentials{}
	}

	return f.http.DoWithEphemeralHeaders(ctx, method, url, req, rsp, creds.SubaccountAsHeaders())
}

func (f *ftxClient) signBeforeDo(ctx context.Context, method, endpoint string, req, rsp interface{}, pagination *PaginationFilter, credentials *auth.Credentials) error {
	ts := strconv.FormatInt(time.Now().UTC().Unix()*1000, 10)
	preparedEndpoint := buildEndpoint(endpoint, pagination)

	signature, err := auth.SignRequest(preparedEndpoint, method, ts, req, credentials)
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
