package transport

import (
	"context"
	"time"

	"github.com/gorilla/websocket"
)

type WsConfig struct {
	Endpoint string
}

type Websocket struct {
	cfg  *WsConfig
	conn *websocket.Conn
}

type WsMessage struct {
	Type     int
	Raw      []byte
	Received time.Time
	Sender   string
}

func NewWebsocket(cfg *WsConfig) *Websocket {
	return &Websocket{
		cfg:  cfg,
		conn: &websocket.Conn{},
	}
}

func (ws *Websocket) Send(ctx context.Context, msg *WsMessage) error {
	return nil
}
func (ws *Websocket) Receiver(ctx context.Context) (chan *WsMessage, chan error) {
	c := make(chan *WsMessage)
	errC := make(chan error)
	go func() {
		defer func() {
			close(c)
			close(errC)
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
				Type:     t,
				Raw:      msg,
				Sender:   ws.cfg.Endpoint,
				Received: time.Now(),
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
