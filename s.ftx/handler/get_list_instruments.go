package handler

import (
	"context"

	"swallowtail/libraries/gerrors"
	"swallowtail/s.ftx/client"
	"swallowtail/s.ftx/marshaling"
	ftxproto "swallowtail/s.ftx/proto"
)

// ListFTXInstruments ...
func ListFTXInstruments(
	ctx context.Context, in *ftxproto.ListFTXInstrumentsRequest,
) (*ftxproto.ListFTXInstrumentsResponse, error) {
	rsp, err := client.ListInstruments(ctx, &client.ListInstrumentsRequest{})
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_list_ftx_instruments", nil)
	}

	protos := marshaling.InstrumentsDTOToProtos(rsp.Instruments)

	return &ftxproto.ListFTXInstrumentsResponse{
		Instruments: protos,
	}, nil
}
