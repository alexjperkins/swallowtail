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
	switch {
	case in.UserId == "":
		return nil, terrors.PreconditionFailed("missing_param.user_id", "Missing Param; user_id", nil)
	case in.Password == "":
		return nil, terrors.PreconditionFailed("missing_param.username", "Missing Param; username", nil)
	case in.Username == "":
		return nil, terrors.PreconditionFailed("missing_param.password", "Missing Param; password", nil)
	}

	errParams := map[string]string{
		"user_id":  in.UserId,
		"username": in.Username,
	}

	account, err := dao.ReadAccountByUserID(ctx, in.UserId)
	switch {
	case terrors.Is(err, terrors.ErrNotFound):
		// This is fine; we don't already have an account - so let's create one.
	case err != nil:
		return nil, terrors.Augment(err, "Failed to read account by user id; couldn't check if account already exists", errParams)
	case account != nil:
		// We've read out an already existing account, let's return an error.
		errParams["account_created"] = account.Created.String()
		return nil, terrors.PreconditionFailed("account-already-exists", "Cannot create an account with this user id; account already exists", errParams)
	}

	// TODO: SQL doesn't like this
	_, err = util.Sha256Hash(in.Password)
	if err != nil {
		return nil, terrors.Augment(err, "Failed to create an account; error hashing password", errParams)
	}

	account = &domain.Account{
		UserID:            in.UserId,
		Username:          in.Username,
		Password:          in.Password,
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
