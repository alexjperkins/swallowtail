package handler

import (
	"context"
	"time"

	"github.com/monzo/slog"
	"github.com/monzo/terrors"

	"swallowtail/libraries/gerrors"
	"swallowtail/s.account/dao"
	"swallowtail/s.account/domain"
	accountproto "swallowtail/s.account/proto"
)

// CreateAccount ...
func (a *AccountService) CreateAccount(
	ctx context.Context, in *accountproto.CreateAccountRequest,
) (*accountproto.CreateAccountResponse, error) {
	switch {
	case in.UserId == "":
		return nil, gerrors.FailedPrecondition("missing_param.user_id", nil)
	case in.Username == "":
		return nil, gerrors.FailedPrecondition("missing_param.username", nil)
	}

	errParams := map[string]string{
		"user_id":  in.UserId,
		"username": in.Username,
	}

	account, err := dao.ReadAccountByUserID(ctx, in.UserId)
	switch {
	case gerrors.Is(err, gerrors.ErrNotFound, "account_not_found"):
		// This is fine; we don't already have an account - so let's create one.
	case err != nil:
		return nil, gerrors.Augment(err, "Failed to read account by user id; couldn't check if account already exists", errParams)
	case account != nil:
		// We've read out an already existing account, let's return an error.
		errParams["account_created"] = account.Created.String()
		return nil, gerrors.AlreadyExists("account-already-exists", errParams)
	}

	account = &domain.Account{
		UserID:            in.UserId,
		Username:          in.Username,
		Email:             in.Email,
		HighPriorityPager: in.HighPriorityPager.String(),
		LowPriorityPager:  in.LowPriorityPager.String(),
	}

	if err := dao.CreateAccount(ctx, account); err != nil {
		return nil, terrors.Augment(err, "Failed to create account", errParams)
	}

	slog.Info(ctx, "Created new account", errParams)

	if err := notifyPulseChannel(ctx, in.UserId, in.Username, time.Now().UTC().Truncate(time.Second)); err != nil {
		// Best effort
		slog.Error(ctx, "Failed to notify pulse channel for %s, %s: Error: ", in.UserId, err)
	}

	return &accountproto.CreateAccountResponse{}, nil
}
