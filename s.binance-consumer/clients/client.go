package clients

import (
	"context"
	"fmt"
	"strings"
	"swallowtail/libraries/multiplexing"
	"swallowtail/libraries/transport"
	"swallowtail/s.binance/domain"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/monzo/slog"
)

var (
	streamURI string = "wss://stream.binance.com:9443"

	// timeouts
	heartbeatPeriod  = 3 * time.Minute
	heartbeatTimeout = 5 * time.Second // Best effort

	// Default Pongs
	heartbeatPong = &domain.BinanceStreamPong{}

	// Constructors per mapping TODO move to domain
	constructortMapping = map[string]func() interface{}{
		// Trades
		"trade": func() interface{} {
			return &domain.BinanceTradeEvent{}
		},
		// Candlesticks
		"kline": func() interface{} {
			return &domain.BinanceKlineEvent{}
		},
	}
)

// StreamClient is a client for interacting with Binance
type StreamClient struct {
	ws    transport.StreamingJSONTransport
	wsCfg *transport.WsConfig
	wsMtx sync.Mutex // For reconnecting

	instrument string
	option     string

	multiplexers []*multiplexing.Multiplex
	errCh        chan error
	done         chan struct{}
}

func NewStreamingClient(ctx context.Context, instrument string, multiplexers []*multiplexing.Multiplex) *StreamClient {
	if ok, err := validateInstrument(instrument); !ok {
		panic(fmt.Sprintf("Invalid instrument: %s", err.Error()))
	}

	option, err := parseOptionFromInstrument(instrument)
	if err != nil {
		panic(err.Error())
	}

	endpoint, _ := buildStreamEndpoint(streamURI, instrument)

	slog.Info(ctx, "Starting streaming client -> %s", endpoint)
	bufSize := 16
	cfg := &transport.WsConfig{
		Endpoint: endpoint,
		BufSize:  bufSize,
	}
	c := &StreamClient{
		ws:           transport.StreamingJSONTransport(ctx, cfg),
		wsCfg:        cfg,
		wsMtx:        sync.Mutex{},
		multiplexers: multiplexers,
		instrument:   instrument,
		option:       option,
		done:         make(chan struct{}, 1),
	}
	errCh := c.Start(ctx, bufSize)
	c.errCh = errCh
	return c
}

func (c *StreamClient) Start(ctx context.Context, bufSize int) chan error {
	// Pull constructor method to build event before marshalling.
	eventConstructor, ok := constructortMapping[c.option]
	if !ok {
		panic(fmt.Sprintf("Cannot handle events for instrument: %s", c.instrument))
	}

	// Start reciever websocket receiver
	ch, errCh := c.ws.Receiver(ctx)

	eventCh := make(chan interface{}, bufSize)

	// < 24 hour ticker, since Binance cancels after 24; thus we must reconnect.
	t := time.NewTicker(time.Hour * 23)

	// Multiplex messages from websocket stream to consumers
	go func() {
		defer slog.Info(nil, "Client stopped.")
		for {
			select {
			case rmsg := <-ch:
				e, err := domain.WsMsgToBinanceEvent(rmsg, eventConstructor)
				if err != nil {
					errCh <- err
				}
				eventCh <- e
			case err := <-errCh:
				// Just print to console for now before we have a logger DI.
				slog.Info(ctx, "Message failed: %s", err.Error())
			case <-c.done:
				for _, m := range c.multiplexers {
					m.Stop()
				}
				c.ws.StopReceiver()
				return
			case <-ctx.Done():
				for _, m := range c.multiplexers {
					m.Stop()
				}
				c.ws.StopReceiver()
				return
			case <-t.C:
				slog.Info(ctx, "Reconnecting websocket before binance closes")
				c.wsMtx.Lock()
				c.ws = transport.NewWebsocket(ctx, c.wsCfg)
				// Arbitary wait to allow for reconnection. We probably want a much better way of doing this here.
				// We also want to pause other goroutines here as well. This is bad practise.
				time.Sleep(time.Second * 2)
				c.wsMtx.Unlock()
			}
		}
	}()
	for _, m := range c.multiplexers {
		m.Start(ctx, eventCh)
	}
	// Start heartbeat server
	go c.heartbeat(ctx)
	return errCh
}

func (c *StreamClient) Stop() {
	slog.Info(nil, "Stopping client...")
	c.done <- struct{}{}
	close(c.done)
}

func (c *StreamClient) heartbeat(ctx context.Context) {
	t := time.NewTicker(heartbeatPeriod)
	defer func() {
		t.Stop()
		slog.Info(ctx, "Stopping heartbeat server.")
	}()
	for {
		select {
		case <-t.C:
			slog.Info(ctx, "Sending heartbeat pong")
			c.ws.Send(
				nil,
				&transport.WsMessage{
					Type:    websocket.PongMessage,
					Raw:     []byte("--heartbeat--"),
					Created: time.Now(),
				},
				heartbeatTimeout,
			)
		case <-ctx.Done():
			return
		case <-c.done:
			return
		}
	}
}

func buildStreamEndpoint(uri, instrument string) (string, error) {
	return fmt.Sprintf("%s/ws/%s", uri, instrument), nil
}

func validateInstrument(instrument string) (bool, error) {
	// TODO
	return true, nil
}

func parseOptionFromInstrument(instrument string) (string, error) {
	a := strings.Split(instrument, "@")
	if len(a) != 2 {
		return "", fmt.Errorf("Cannot parse option from: %s", instrument)
	}
	a = strings.Split(a[1], "_")
	return a[0], nil
}
