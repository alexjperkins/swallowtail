package handler

import (
	"context"

	"github.com/monzo/terrors"

	"swallowtail/s.account/dao"
	"swallowtail/s.account/domain"
	"swallowtail/s.account/pager"
	accountproto "swallowtail/s.account/proto"
)

func (s *AccountService) PUTPageAccount(
	ctx context.Context, in *accountproto.PageAccountRequest,
) (*accountproto.PageAccountResponse, error) {
	errParams := map[string]string{
		"account_id": in.Id,
	}

	account, err := dao.ReadAccountByID(ctx, in.Id)
	if err != nil {
		return nil, terrors.Augment(err, "Failed to read account", errParams)
	}

	var pagerType string
	switch in.Priority {
	case accountproto.PagerPriority_HIGH:
		pagerType = account.HighPriorityPager
	case accountproto.PagerPriority_LOW:
		pagerType = account.LowPriorityPager
	default:
		pagerType = account.LowPriorityPager
	}

	pager, err := pager.GetPagerByID(pagerType)
	if err != nil {
		errParams["pager_type"] = pagerType
		return nil, terrors.Augment(err, "Invalid pager type set to account", errParams)
	}

	identifier, err := getIdentifierFromAccount(account, pagerType)
	if err != nil {
		return nil, terrors.Augment(err, "Cannot page user; missing identifier on account", errParams)
	}

	if err := pager.Page(ctx, identifier, in.Content); err != nil {
		return nil, terrors.Augment(err, "Failed to page user", errParams)
	}

	return &accountproto.PageAccountResponse{}, nil
}

func getIdentifierFromAccount(account *domain.Account, pagerType string) (string, error) {
	errParams := map[string]string{
		"account_id": account.AccountID,
	}

	switch pagerType {
	case accountproto.PagerType_DISCORD.String():
		if account.DiscordID == "" {
			return "", terrors.PreconditionFailed("account-discord-id-missing", "Cannot page account via discord; discord id not set", errParams)
		}
		return account.DiscordID, nil
	case accountproto.PagerType_EMAIL.String():
		if account.Email == "" {
			return "", terrors.PreconditionFailed("account-email-missing", "Cannot page account via email; email not set", errParams)
		}
		return account.Email, nil
	case accountproto.PagerType_PHONE.String():
		if account.PhoneNumber == "" {
			return "", terrors.PreconditionFailed("account-phone-number-missing", "Cannot page account via phone number; phone number not set", errParams)
		}
		return account.PhoneNumber, nil
	case accountproto.PagerType_SMS.String():
		if account.PhoneNumber == "" {
			return "", terrors.PreconditionFailed("account-phone-number-missing", "Cannot page account via sms; sms not set", errParams)
		}
		return account.PhoneNumber, nil
	}

	errParams["pager_type"] = pagerType
	return "", terrors.PreconditionFailed("unknown-pager-type", "Cannot page account; unknown pager type", errParams)
}
