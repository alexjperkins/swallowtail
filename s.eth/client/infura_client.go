package client

import (
	"context"
	"fmt"
	"net/http"
	"swallowtail/libraries/transport"
	"swallowtail/s.eth/domain"
	"time"

	"gopkg.in/tomb.v2"

	"github.com/monzo/slog"
	"github.com/monzo/terrors"
)

const (
	infuraWSEndpoint = "wss://ropseten.infura.io/ws/v3"
	projectID        = "1dde09bdeda14412a9d0ac522172ae98"
	infuraClientID   = "infura-client"
)

func NewInfuraClient(ctx context.Context) (EthClient, error) {
	cfg := &transport.WsConfig{
		Endpoint: fmt.Sprintf("%s/%s", infuraWSEndpoint, projectID),
		BufSize:  32,
	}
	i := &infuraClient{
		&http.Client{
			Timeout:   time.Duration(30 * time.Second),
			Transport: &metricsRoundTripper{},
		},
		transport.NewWebsocket(ctx, cfg),
	}
	go i.keepAlive()
	return i, nil
}

type infuraClient struct {
	*http.Client
	*transport.Websocket
}

func (i *infuraClient) SubscribePendingTransactions(ctx context.Context) (<-chan *domain.EthMempoolTxEvent, error) {
	i.SendJSON(ctx, &domain.SubscribeToInfuraRequest{
		Id:      0,
		JSONRPC: "2.0",
		Method:  "eth_subscribe",
		Params: []string{
			"newPendingTransactions",
		},
	}, time.Duration(30*time.Second))

	// The first message received should be a subscribed message.
	rsp := &domain.SubscribeToInfuraResponse{}
	err := i.ReceiveJSON(rsp)
	if err != nil {
		// If we have something expected then return; we may want to put this in a retry loop.
		return nil, terrors.Augment(err, "Failed to subscribe to pending transactions", map[string]string{
			"client_id": infuraClientID,
		})
	}

	receiverCh := make(chan *domain.EthMempoolTxEvent, 32)
	t, _ := tomb.WithContext(ctx)
	t.Go(func() error {
		for {
			msg := &domain.EthMempoolTxEvent{}
			err := i.ReceiveJSON(msg)
			if err != nil {
				return err
			}
			receiverCh <- msg
		}
	})
	return receiverCh, err
}

func (i *infuraClient) StopReceiver() {
}

func (i *infuraClient) keepAlive() {
}

type metricsRoundTripper struct{}

func (m *metricsRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	start := time.Now()
	// TODO replace with metrics
	defer slog.Info(r.Context(), "Request: %+v took %v", r, time.Now().Sub(start))
	return http.DefaultTransport.RoundTrip(r)
}
