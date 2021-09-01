package handler

import (
	"context"
	"swallowtail/libraries/gerrors"
	accountproto "swallowtail/s.account/proto"
)

// setUserAsFuturesMember
//
// Sets the user as a futures member.
func setUserAsFuturesMember(ctx context.Context, userID string) error {
	_, err := (&accountproto.UpdateAccountRequest{
		ActorId:   accountproto.ActorSystemPayments,
		UserId:    userID,
		IsFutures: true,
	}).Send(ctx).Response()
	if err != nil {
		return gerrors.Augment(err, "failed_to_set_user_as_futures_member", map[string]string{
			"user_id": userID,
		})
	}

	return nil
}
