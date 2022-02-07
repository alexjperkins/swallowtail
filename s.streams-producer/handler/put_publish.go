package handler

import (
	"context"

	"swallowtail/libraries/gerrors"
	streamsproducerproto "swallowtail/s.streams-producer/proto"
)

// Publish ...
func (s *StreamsProducerService) Publish(
	ctx context.Context, in *streamsproducerproto.PublishRequest,
) (*streamsproducerproto.PublishResponse, error) {
	return nil, gerrors.Unimplemented("publish.unimplemented", nil)
}
