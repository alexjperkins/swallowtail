package handler

import (
	"context"

	bookmarkerproto "swallowtail/s.bookmarker/proto"
)

// TmpPing is tmp endpoint for validate k8s deployments.
func (s *BookmarkerService) TmpPing(
	ctx context.Context, in *bookmarkerproto.TmpPingRequest,
) (*bookmarkerproto.TmpPingResponse, error) {
	return &bookmarkerproto.TmpPingResponse{
		Message: in.Message,
	}, nil
}
