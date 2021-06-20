package client

import (
	"context"
	"fmt"
	"swallowtail/libraries/transport"
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

func NewInfuraClient(ctx context.Context) (Web3Client, error) {
	cfg := &transport.WsConfig{
		Endpoint: fmt.Sprintf("%s/%s", infuraWSEndpoint, projectID),
		BufSize:  32,
	}
	headers := map[string]string{}
	i := &infuraClient{
		c:  transport.NewHTTPClient(ctx, time.Duration(30*time.Second), headers),
		ws: transport.NewWebsocket(ctx, cfg),
		t:  tomb.Tomb{},
	}
	go i.keepAlive()
	return i, nil
}

type infuraClient struct {
	c  transport.HttpClient
	ws *transport.Websocket
	t  tomb.Tomb
}

func (i *infuraClient) SubscribePendingTransactions(ctx context.Context) (<-chan *PendingTransactionEvent, error) {
	i.ws.SendJSON(ctx, SubscribeToInfuraRequest{
		Id:      0,
		JSONRPC: "2.0",
		Method:  "eth_subscribe",
		Params: []string{
			"newPendingTransactions",
		},
	}, time.Duration(30*time.Second))

	// The first message received should be a subscribed message.
	rsp := &SubscribeToInfuraResponse{}
	err := i.ws.ReceiveJSON(rsp)
	if err != nil {
		// If we have something expected then return; we may want to put this in a retry loop.
		return nil, terrors.Augment(err, "Failed to subscribe to pending transactions", map[string]string{
			"client_id": infuraClientID,
		})
	}

	receiverCh := make(chan *PendingTransactionEvent, 32)
	i.t.Go(func() error {
		for {
			msg := &PendingTransactionEvent{}
			err := i.ws.ReceiveJSON(msg)
			if err != nil {
				return err
			}
			receiverCh <- msg
		}
	})
	return receiverCh, err
}

func (i *infuraClient) StopReceiver() {
	slog.Info(context.TODO(), "Stop receiver call received.")
	i.t.Kill(tomb.ErrDying)
}

func (i *infuraClient) keepAlive() {
}
