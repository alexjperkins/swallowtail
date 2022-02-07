package streams

import (
	streamsconsumerproto "swallowtail/s.streams-consumer/proto"
	"sync"
)

type rpcSubscription struct {
	sync.Mutex

	topic   string
	group   string
	command *streamsconsumerproto.Command
	streamsconsumerproto.StreamsconsumerClient

	handler Handler
}
