package handler

import (
	"context"
	"swallowtail/libraries/gerrors"
	accountproto "swallowtail/s.account/proto"
	discordproto "swallowtail/s.discord/proto"
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

func readUserRoles(ctx context.Context, userID string) ([]*discordproto.Role, error) {
	rsp, err := (&discordproto.ReadUserRolesRequest{
		UserId: userID,
	}).Send(ctx).Response()
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_read_user_discord_roles", map[string]string{
			"user_id": userID,
		})
	}

	return rsp.GetRoles(), nil
}

func isUserRegistered(ctx context.Context, userID string) (bool, error) {
	_, err := (&accountproto.ReadAccountRequest{
		UserId: userID,
	}).Send(ctx).Response()
	switch {
	case gerrors.Is(err, gerrors.ErrNotFound, "account_not_found"):
		return false, nil
	case err != nil:
		return false, gerrors.Augment(err, "failed_to_check_if_user_register", nil)
	}

	return true, nil
}
