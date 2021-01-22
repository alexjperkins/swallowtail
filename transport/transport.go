package transport

import "context"

type StreamingTransport interface {
	Send(context.Context, *WsMessage) error
	Receiver(context.Context) (chan *WsMessage, chan error)
}

type Transport interface {
	auth(context.Context) error
	Send(context.Context, interface{}) error
}
