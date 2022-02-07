package handler

import (
	"swallowtail/libraries/gerrors"

	streamsconsumerproto "swallowtail/s.streams-consumer/proto"
)

func (s StreamsConsumerService) Subscribe(
	command *streamsconsumerproto.Command, stream streamsconsumerproto.Streamsconsumer_SubscribeServer,
) error {
	return gerrors.Unimplemented("subscribe.unimplemented", nil)
}
