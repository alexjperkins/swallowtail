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
	"swallowtail/libraries/gerrors"
	"time"

	"google.golang.org/grpc/codes"
)

type HttpClient interface {
	Do(ctx context.Context, method, endpoint string, reqBody, rspBody interface{}) error
	DoWithEphemeralHeaders(ctx context.Context, method, endpoint string, reqBody, rspBody interface{}, headers map[string]string) error
}

func NewHTTPClient(ctx context.Context, timeout time.Duration) HttpClient {
	return &httpClient{
		c: &http.Client{
			Timeout: timeout,
		},
	}
}

type httpClient struct {
	c       *http.Client
	headers map[string]string
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
		return gerrors.FailedPrecondition("bad_request.unmarshal_error", errParams)
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

	rsp, err := h.doRawRequest(ctx, method, url, body, headers)
	if err != nil {
		return err
	}

	defer rsp.Body.Close()
	rspBodyBytes, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return gerrors.Augment(err, "failed_to_read_response_body", errParams)
	}

	if err := json.Unmarshal(rspBodyBytes, rspBody); err != nil {
		return gerrors.FailedPrecondition("bad_request.unmarshal_error", errParams)
	}
	return nil
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

	if err := validateStatusCode(rsp, errParams); err != nil {
		return nil, err
	}

	return rsp, err
}

func (h *httpClient) authorize(req *http.Request, key, value string) {
	req.Header.Add(key, value)

}

func validateStatusCode(rsp *http.Response, errParams map[string]string) error {
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
		code = gerrors.ErrUnauthenticated
	default:
		return gerrors.FailedPrecondition("bad_request", map[string]string{
			"status_code": strconv.Itoa(rsp.StatusCode),
		})
	}

	return gerrors.New(code, msg, errParams)
}
