package client

import (
	"context"
	"fmt"
	"net/http"

	"swallowtail/libraries/gerrors"
	"swallowtail/libraries/transport"
	"swallowtail/s.bitfinex/dto"
)

const (
	bitfinexURL        = "https://api-pub.bitfinex.com"
	bitfinexAPIVersion = "v2"
)

type bitfinexClient struct {
	http transport.HttpClient
}

func (b *bitfinexClient) Ping(ctx context.Context) error {
	if _, err := b.GetStatus(ctx, &dto.GetStatusRequest{}); err != nil {
		return gerrors.Augment(err, "failed_to_establish_bitfinex_connection", nil)
	}

	return nil
}

func (b *bitfinexClient) GetStatus(ctx context.Context, req *dto.GetStatusRequest) (*dto.GetStatusResponse, error) {
	rsp := &dto.GetStatusProxyResponse{}
	if err := b.http.Do(ctx, http.MethodGet, fmt.Sprintf("%s/%s/platform/status", bitfinexURL, bitfinexAPIVersion), nil, rsp); err != nil {
		return nil, gerrors.Augment(err, "failed_get_status.client", nil)
	}

	return &dto.GetStatusResponse{
		Operative: rsp.Operative(),
	}, nil
}
