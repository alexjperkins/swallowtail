package handler

import (
	"context"
	"swallowtail/libraries/gerrors"
	"swallowtail/s.ftx/client"
	"swallowtail/s.ftx/marshaling"
	ftxproto "swallowtail/s.ftx/proto"
)

// ListAccountDeposits ...
func (s *FTXService) ListAccountDeposits(
	ctx context.Context, in *ftxproto.ListAccountDepositsRequest,
) (*ftxproto.ListAccountDepositsResponse, error) {
	switch {
	case in.ActorId == "":
		return nil, gerrors.BadParam("missing_param.actor_id", nil)
	}

	errParams := map[string]string{
		"actor_id": in.ActorId,
	}

	// Check the actor is authorized to make this request.
	ok, err := isValidActor(ctx, in.ActorId)
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_list_account_deposits", errParams)
	}

	if !ok {
		return nil, gerrors.Unauthenticated("failed_to_list_account_deposits.unauthorized_actor", errParams)
	}

	rsp, err := client.ListAccountDeposits(ctx, &client.ListAccountDepositsRequest{}, &client.PaginationFilter{
		// We need millisecond resolution
		Start: in.Start.Seconds / 1_000_000,
		End:   in.Start.Seconds / 1_000_000,
	})

	protoDeposits := marshaling.DepositsDTOToProto(rsp.Deposits)

	return &ftxproto.ListAccountDepositsResponse{
		Deposits: protoDeposits,
	}, gerrors.Unimplemented("failed_to_list_account_deposits.unimplemented", nil)
}
