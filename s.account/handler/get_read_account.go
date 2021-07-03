package handler

import (
	"context"
	"swallowtail/s.account/dao"
	"swallowtail/s.account/marshaling"
	accountproto "swallowtail/s.account/proto"

	"github.com/monzo/terrors"
)

func (a *AccountService) ReadAccount(
	ctx context.Context, in *accountproto.ReadAccountRequest,
) (*accountproto.ReadAccountResponse, error) {
	errParams := map[string]string{
		"user_id": in.UserId,
	}

	account, err := dao.ReadAccountByUserID(ctx, in.UserId)
	if err != nil {
		return nil, terrors.Augment(err, "Failed to read account by user id", errParams)
	}

	return &accountproto.ReadAccountResponse{
		Account: marshaling.AccountDomainToProto(account),
	}, nil
}
