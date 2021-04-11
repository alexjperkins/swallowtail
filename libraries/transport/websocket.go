package transport

import (
	"context"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
	"github.com/monzo/slog"
)

type WsConfig struct {
	Endpoint string
	BufSize  int
}

type Websocket struct {
	cfg          *WsConfig
	conn         *websocket.Conn
	recevierDone chan struct{}
}

type WsMessage struct {
	Type    int
	Raw     []byte
	Created time.Time
	Sender  string
}

func NewWebsocket(ctx context.Context, cfg *WsConfig) *Websocket {
	c, _, err := websocket.DefaultDialer.DialContext(ctx, cfg.Endpoint, nil)
	if err != nil {
		panic("Shit! No websocket connection")
	}
	slog.Info(ctx, fmt.Sprintf("creating ws -> %s", cfg.Endpoint))
	return &Websocket{
		cfg:          cfg,
		conn:         c,
		recevierDone: make(chan struct{}, 1),
	}
}

func (ws *Websocket) Send(ctx context.Context, msg *WsMessage, timeout time.Duration) {
	sent := make(chan error)
	defer close(sent)
	go func() {
		err := ws.conn.WriteMessage(msg.Type, msg.Raw)
		if err != nil {
			return
		}
		sent <- nil
	}()
	select {
	case <-sent:
		return
	case <-time.After(timeout):
		return
	}
}

func (ws *Websocket) BlockingSend(msg *WsMessage) error {
	return ws.conn.WriteMessage(msg.Type, msg.Raw)
}

func (ws *Websocket) BlockingSendJSON(msg interface{}) error {
	return ws.conn.WriteJSON(msg)
}

func (ws *Websocket) SendJSON(ctx context.Context, msg interface{}, timeout time.Duration) {
	sent := make(chan error)
	defer close(sent)
	go func() {
		err := ws.conn.WriteJSON(msg)
		if err != nil {
			return
		}
		sent <- nil
	}()
	select {
	case <-sent:
		return
	case <-time.After(timeout):
		return
	}
}

func (ws *Websocket) Receiver(ctx context.Context) (chan *WsMessage, chan error) {
	c := make(chan *WsMessage, ws.cfg.BufSize)
	errC := make(chan error, ws.cfg.BufSize)
	slog.Info(ctx, "Starting websocket receiver...")
	go func() {
		defer func() {
			close(c)
			close(errC)
			slog.Info(nil, "Websocket receiver stopped.")
		}()
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			t, msg, err := ws.conn.ReadMessage()
			if err != nil {
				errC <- err
			}
			wsMsg := &WsMessage{
				Type:    t,
				Raw:     msg,
				Sender:  ws.cfg.Endpoint,
				Created: time.Now(),
			}

			select {
			case c <- wsMsg:
			default:
				// best effort for now
			}
		}
	}()
	return c, errC
}

func (ws *Websocket) StopReceiver() {
	defer slog.Info(nil, "Web receiver stopping...")
	ws.recevierDone <- struct{}{}
}
