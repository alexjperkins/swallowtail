package handler

import (
	"context"
	"swallowtail/s.account/dao"
	"swallowtail/s.account/marshaling"
	accountproto "swallowtail/s.account/proto"

	"github.com/monzo/slog"
	"github.com/monzo/terrors"
)

// UpdateAccount ...
func (s *AccountService) UpdateAccount(
	ctx context.Context, in *accountproto.UpdateAccountRequest,
) (*accountproto.UpdateAccountResponse, error) {
	switch {
	case in.UserId == "":
		return nil, terrors.PreconditionFailed("missing-param", "Missing parameter; user id cannot be empty", nil)
	}

	errParams := map[string]string{
		"user_id": in.UserId,
	}

	mutation := marshaling.UpdateAccountProtoToDomain(in)
	account, err := dao.UpdateAccount(ctx, mutation)

	switch {
	case terrors.Is(err, terrors.ErrNotFound):
		return nil, terrors.Augment(err, "Failed to update account; account with that user id doesn't exist", errParams)
	case err != nil:
		return nil, terrors.Augment(err, "Failed to update account", errParams)
	}

	errParams["updated"] = account.Updated.String()
	slog.Info(ctx, "Updated account", errParams)

	return &accountproto.UpdateAccountResponse{
		Account: marshaling.AccountDomainToProto(account),
	}, nil
}
