package handler

import bookmarkerproto "swallowtail/s.bookmarker/proto"

// BookmarkerService defines the gRPC service for the bookmarker service.
type BookmarkerService struct {
	*bookmarkerproto.UnimplementedBookmarkerServer
}
