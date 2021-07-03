package handler

import (
	"context"
	"swallowtail/libraries/util"
	"swallowtail/s.account/dao"
	"swallowtail/s.account/domain"
	accountproto "swallowtail/s.account/proto"

	"github.com/monzo/slog"
	"github.com/monzo/terrors"
)

// CreateAccount ...
func (a *AccountService) CreateAccount(
	ctx context.Context, in *accountproto.CreateAccountRequest,
) (*accountproto.CreateAccountResponse, error) {
	errParams := map[string]string{
		"user_id":  in.UserId,
		"username": in.Username,
	}

	account := &domain.Account{
		UserID:            in.UserId,
		Username:          in.Username,
		Password:          util.Sha256Hash(in.Password),
		Email:             in.Email,
		HighPriorityPager: in.HighPriorityPager.String(),
		LowPriorityPager:  in.LowPriorityPager.String(),
	}

	if err := dao.CreateAccount(ctx, account); err != nil {
		return nil, terrors.Augment(err, "Failed to create account", errParams)
	}

	slog.Info(ctx, "Created new account", errParams)

	return &accountproto.CreateAccountResponse{}, nil
}
