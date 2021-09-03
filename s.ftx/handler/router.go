package handler

import (
	ftxproto "swallowtail/s.ftx/proto"
)

// FTXService defines the gRPC service for the FTX.
type FTXService struct {
	*ftxproto.UnimplementedFtxServer
}
