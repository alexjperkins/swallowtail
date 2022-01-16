package transport

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/monzo/slog"
	"google.golang.org/grpc/codes"

	"swallowtail/libraries/gerrors"
)

const (
	RequestErrorMessageDetailKey = "http_request_error_body"
)

type HTTPRateLimiter interface {
	RefreshWait(header http.Header, statusCode int)
	Wait()
}

type HttpClient interface {
	Do(ctx context.Context, method, endpoint string, reqBody, rspBody interface{}) error
	DoWithEphemeralHeaders(ctx context.Context, method, endpoint string, reqBody, rspBody interface{}, headers map[string]string) error
}

func NewHTTPClient(timeout time.Duration, rateLimiter HTTPRateLimiter) HttpClient {
	return &httpClient{
		c: &http.Client{
			Timeout: timeout,
		},
		rateLimiter: rateLimiter,
	}
}

type httpClient struct {
	c           *http.Client
	headers     map[string]string
	rateLimiter HTTPRateLimiter
}

func (h *httpClient) WithHeaders(headers map[string]string) {
	h.headers = headers
}

func (h *httpClient) Do(ctx context.Context, method, url string, reqBody, rspBody interface{}) error {
	errParams := map[string]string{
		"method": method,
		"url":    url,
	}

	var body io.Reader
	if reqBody != nil {
		reqBodyBytes, err := json.Marshal(reqBody)
		if err != nil {
			return gerrors.Augment(err, "failed_to_marshal_request_body", errParams)
		}
		body = bytes.NewReader(reqBodyBytes)
	}

	rsp, err := h.doRawRequest(ctx, method, url, body, h.headers)
	if err != nil {
		return err
	}

	defer rsp.Body.Close()
	rspBodyBytes, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return gerrors.Augment(err, "failed_to_read_response_body", errParams)
	}

	if err := json.Unmarshal(rspBodyBytes, rspBody); err != nil {
		return gerrors.Augment(err, "bad_request.unmarshal_error", errParams)
	}

	return nil
}

func (h *httpClient) DoWithEphemeralHeaders(ctx context.Context, method, url string, reqBody, rspBody interface{}, headers map[string]string) error {
	errParams := map[string]string{
		"method": method,
		"url":    url,
	}

	var body io.Reader
	if reqBody != nil {
		reqBodyBytes, err := json.Marshal(reqBody)
		if err != nil {
			return gerrors.Augment(err, "failed_to_marshal_request_body", errParams)
		}
		body = bytes.NewReader(reqBodyBytes)
	}

	rsp, err := h.doRawRequestWithRateLimit(ctx, method, url, body, headers)
	if err != nil {
		return err
	}

	defer rsp.Body.Close()
	rspBodyBytes, err := ioutil.ReadAll(rsp.Body)

	if err != nil {
		return gerrors.Augment(err, "failed_to_read_response_body", errParams)
	}

	if err := json.Unmarshal(rspBodyBytes, rspBody); err != nil {
		slog.Error(ctx, "Response body for marshaling failure: %v", string(rspBodyBytes))
		return gerrors.Augment(err, "bad_request.unmarshal_error", errParams)
	}

	return nil
}

func (h *httpClient) doRawRequestWithRateLimit(ctx context.Context, method, url string, body io.Reader, headers map[string]string) (*http.Response, error) {
	if h.rateLimiter == nil {
		return h.doRawRequest(ctx, method, url, body, headers)
	}

	h.rateLimiter.Wait()
	rsp, err := h.doRawRequest(ctx, method, url, body, headers)
	if err != nil {
		return nil, err
	}

	h.rateLimiter.RefreshWait(rsp.Header, rsp.StatusCode)
	return rsp, err
}

func (h *httpClient) doRawRequest(ctx context.Context, method, url string, body io.Reader, headers map[string]string) (*http.Response, error) {
	errParams := map[string]string{
		"method": method,
		"url":    url,
	}

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_execute_request", errParams)
	}

	for k, v := range headers {
		h.authorize(req, k, v)
	}

	rsp, err := h.c.Do(req)
	if err != nil {
		return nil, err
	}

	if err := validateStatusCode(rsp); err != nil {
		// Best effort; here we attempt to read the response bytes.
		defer rsp.Body.Close()
		rspBodyBytes, _ := ioutil.ReadAll(rsp.Body)
		slog.Error(ctx, "Failed request: %s %s Response: %+v, %s", method, url, rsp, string(rspBodyBytes))

		errParams[RequestErrorMessageDetailKey] = string(rspBodyBytes)
		return nil, gerrors.Augment(err, "failed_to_execute_request.status_code", errParams)
	}

	return rsp, err
}

func (h *httpClient) authorize(req *http.Request, key, value string) {
	req.Header.Set(key, value)

}

func validateStatusCode(rsp *http.Response) error {
	if rsp.StatusCode >= 200 && rsp.StatusCode < 300 {
		return nil
	}

	msg := fmt.Sprintf("API request failed with status: %s", rsp.Status)
	var code codes.Code
	switch rsp.StatusCode {
	case 401:
		code = gerrors.ErrUnauthenticated
	case 404:
		code = gerrors.ErrNotFound
	case 429:
		code = gerrors.ErrRateLimited
	default:
		return gerrors.FailedPrecondition("bad_request", map[string]string{
			"status_code": strconv.Itoa(rsp.StatusCode),
		})
	}

	return gerrors.New(code, msg, map[string]string{
		"status_code": strconv.Itoa(rsp.StatusCode),
	})
}
