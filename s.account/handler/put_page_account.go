package handler

import (
	"context"

	"github.com/monzo/terrors"

	"swallowtail/s.account/dao"
	"swallowtail/s.account/pager"
	accountproto "swallowtail/s.account/proto"
)

func (s *AccountService) PageAccount(
	ctx context.Context, in *accountproto.PageAccountRequest,
) (*accountproto.PageAccountResponse, error) {
	errParams := map[string]string{
		"account_id": in.UserId,
	}

	account, err := dao.ReadAccountByUserID(ctx, in.UserId)
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
