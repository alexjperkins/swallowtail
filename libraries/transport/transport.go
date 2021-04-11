package transport

import (
	"context"
	"time"
)

type StreamingTransport interface {
	Send(context.Context, *WsMessage, time.Duration)
	BlockingSend(*WsMessage) error
	Receiver(context.Context) (chan *WsMessage, chan error)
	StopReceiver()
}

type StreamingJSONTransport interface {
	StreamingTransport
	SendJSON(context.Context, interface{}, time.Duration)
	BlockingSendJSON(interface{}) error
}

type Transport interface {
	auth(context.Context) error
	Send(context.Context, interface{}) error
}
