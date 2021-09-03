package handler

import (
	"context"
	"swallowtail/libraries/gerrors"
	accountproto "swallowtail/s.account/proto"
	ftxproto "swallowtail/s.ftx/proto"
)

func isValidActor(ctx context.Context, actorID string) (bool, error) {
	if actorID == ftxproto.FTXDepositAccountActorPaymentsSystem {
		return true, nil
	}

	// Check the actor is authorized to make this request.
	account, err := (&accountproto.ReadAccountRequest{
		UserId: actorID,
	}).Send(ctx).Response()
	if err != nil {
		return false, gerrors.Augment(err, "failed_to_list_account_deposits.failed_to_read_account_of_actor", nil)
	}

	return account.GetAccount().IsAdmin, nil
}
