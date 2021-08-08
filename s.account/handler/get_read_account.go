package handler

import (
	"context"
	"swallowtail/libraries/gerrors"
	"swallowtail/s.account/dao"
	"swallowtail/s.account/marshaling"
	accountproto "swallowtail/s.account/proto"
)

// ReadAccount reads an account via the user ID, which is the discord ID.
func (a *AccountService) ReadAccount(
	ctx context.Context, in *accountproto.ReadAccountRequest,
) (*accountproto.ReadAccountResponse, error) {
	errParams := map[string]string{
		"user_id": in.UserId,
	}

	account, err := dao.ReadAccountByUserID(ctx, in.UserId)
	switch {
	case gerrors.Is(err, gerrors.ErrNotFound, "account-not-found"):
		return nil, gerrors.Augment(err, "This user doesn't have an existing account", errParams)
	case err != nil:
		return nil, gerrors.Augment(err, "Failed to read account", errParams)
	}

	return &accountproto.ReadAccountResponse{
		Account: marshaling.AccountDomainToProto(account),
	}, nil
}
