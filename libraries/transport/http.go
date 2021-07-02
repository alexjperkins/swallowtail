package transport

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/monzo/slog"
	"github.com/monzo/terrors"
)

type HttpClient interface {
	DoRequest(ctx context.Context, method, endpoint string, reqBody, rspBody interface{}) error
}

func NewHTTPClient(ctx context.Context, timeout time.Duration, headers map[string]string) HttpClient {
	return &httpClient{
		headers: headers,
		c: &http.Client{
			Timeout: timeout,
		},
	}
}

type httpClient struct {
	c       *http.Client
	headers map[string]string
}

func (h *httpClient) DoRequest(ctx context.Context, method, url string, reqBody, rspBody interface{}) error {
	errParams := map[string]string{
		"method": method,
		"url":    url,
	}

	var body io.Reader
	if reqBody != nil {
		reqBodyBytes, err := json.Marshal(reqBody)
		if err != nil {
			return terrors.Augment(err, "Failed to marshal API request body", errParams)
		}
		body = bytes.NewReader(reqBodyBytes)
	}

	rsp, err := h.doRawRequest(ctx, method, url, body, errParams)
	if err != nil {
		return err
	}

	defer rsp.Body.Close()
	rspBodyBytes, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return terrors.Augment(err, "Failed to read API response body", errParams)
	}

	if err := json.Unmarshal(rspBodyBytes, rspBody); err != nil {
		return terrors.BadRequest("unmarshal_error", "Failed to unmarshal API response", errParams)
	}
	return nil
}

func (h *httpClient) doRawRequest(ctx context.Context, method, url string, body io.Reader, errParams map[string]string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, terrors.Augment(err, "Failed to execute API request", errParams)
	}

	for k, v := range h.headers {
		h.authorize(req, k, v)
	}

	slog.Debug(ctx, "Sending API request", errParams)

	rsp, err := h.c.Do(req)
	if err != nil {
		return nil, err
	}

	slog.Debug(ctx, "Got API response", errParams)
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
	var code string
	switch rsp.StatusCode {
	case 401:
		code = terrors.ErrUnauthorized
	case 404:
		code = terrors.ErrNotFound
	case 429:
		code = terrors.ErrRateLimited
	default:
		code = terrors.ErrBadRequest + ".bad_status"
	}

	return terrors.New(code, msg, errParams)
}
