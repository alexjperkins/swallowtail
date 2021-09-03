package handler

import (
	"context"
	"swallowtail/libraries/gerrors"
	accountproto "swallowtail/s.account/proto"
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
	account, err := (&accountproto.ReadAccountRequest{
		UserId: in.ActorId,
	}).Send(ctx).Response()
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_list_account_deposits.failed_to_read_account_of_actor", errParams)
	}

	if !account.GetAccount().IsAdmin {
		return nil, gerrors.Unauthenticated("failed_to_list_account_deposits.actor_unauthorized", errParams)
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
