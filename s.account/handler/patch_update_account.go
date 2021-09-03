package handler

import (
	"context"
	"swallowtail/libraries/gerrors"
	"swallowtail/s.account/dao"
	"swallowtail/s.account/marshaling"
	accountproto "swallowtail/s.account/proto"

	"github.com/monzo/slog"
)

// UpdateAccount ...
func (s *AccountService) UpdateAccount(
	ctx context.Context, in *accountproto.UpdateAccountRequest,
) (*accountproto.UpdateAccountResponse, error) {
	switch {
	case in.UserId == "":
		return nil, gerrors.BadParam("missing_param.user_id", nil)
	case (in.IsAdmin || in.IsFutures) && !isValidActorID(in.ActorId):
		// Here if the user is setting a users futures or admin level; then we require a certain actor id.
		return nil, gerrors.BadParam("missing_param.actor_id", map[string]string{
			"actor_id": in.ActorId,
		})
	}

	errParams := map[string]string{
		"user_id": in.UserId,
	}

	mutation := marshaling.UpdateAccountProtoToDomain(in)

	account, err := dao.UpdateAccount(ctx, mutation)
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_update_account", errParams)
	}

	slog.Info(ctx, "Updated account %s: %v", account.UserID, account.Updated)

	return &accountproto.UpdateAccountResponse{
		Account: marshaling.AccountDomainToProto(account),
	}, nil
}
