package binance

import (
	"context"
	"fmt"
	"swallowtail/clients/binance/subscribers"
	"swallowtail/transport"
	"sync"

	"github.com/monzo/slog"
)

var (
	api       string
	secret    string
	streamUri string = "wss://stream.binance.com:9443"
)

// StreamClient is a client for interacting with Binance
type StreamClient struct {
	ws  transport.StreamingTransport
	cli transport.Transport

	symbol string
	option string

	subscribers []subscribers.Subscriber // move to map?
	subMu       sync.RWMutex
}

func New(ctx context.Context, symbol, option string) *StreamClient {
	endpoint, _ := buildStreamEndpoint(streamUri, symbol, option)
	cfg := &transport.WsConfig{
		Endpoint: endpoint,
	}
	c := &StreamClient{
		ws: transport.NewWebsocket(cfg),
	}
	go c.Start(ctx)
	return c
}

func (c *StreamClient) Subscribe(s subscribers.Subscriber) {
	c.subMu.Lock()
	defer c.subMu.Unlock()
	c.subscribers = append(c.subscribers, s)
}

func (c *StreamClient) GetSubscribers() []subscribers.Subscriber {
	c.subMu.RLock()
	defer c.subMu.RUnlock()
	return c.subscribers
}

func (c *StreamClient) Start(ctx context.Context) {
	ch, errCh := c.ws.Receiver(ctx)
	for {
		select {
		case rmsg := <-ch:
			msg, err := wsMsgToBinanceMsg(rmsg)
			if err != nil {
				// do something
			}
			for _, s := range c.GetSubscribers() {
				// do we want this to block? probs not
				s.Send(ctx, msg)
			}
		case err := <-errCh:
			// Just print to console for now before we have a logger DI.
			slog.Info(ctx, "Message failed: %s", err.Error())
		case <-ctx.Done():
			for _, s := range c.GetSubscribers() {
				s.Close()
			}
			return
		}
	}
}

func buildStreamEndpoint(uri, symbol, option string) (string, error) {
	withSymbol := fmt.Sprintf("%s/ws/%s", uri, symbol)
	if option == "" {
		return withSymbol, nil
	}
	return fmt.Sprintf("%s@%s", withSymbol, option), nil
}
