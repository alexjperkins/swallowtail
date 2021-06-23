package handler

import (
	"github.com/monzo/terrors"
	"github.com/monzo/typhon"

	"swallowtail/s.account/dao"
	"swallowtail/s.account/domain"
	"swallowtail/s.account/pager"
	accountproto "swallowtail/s.account/proto"
)

func PUTPageAccount(req typhon.Request) typhon.Response {
	body := accountproto.PageAccountRequest{}
	if err := req.Decode(&body); err != nil {
		return typhon.Response{Error: err}
	}

	errParams := map[string]string{
		"account_id": body.Id,
	}

	account, err := dao.ReadAccountByID(req, body.Id)
	if err != nil {
		return typhon.Response{Error: terrors.Augment(err, "Failed to read account", errParams)}
	}

	var pagerType string
	switch body.Priority {
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
		return typhon.Response{Error: terrors.Augment(err, "Invalid pager type set to account", errParams)}
	}

	identifier, err := getIdentifierFromAccount(account, pagerType)
	if err != nil {
		return typhon.Response{Error: terrors.Augment(err, "Cannot page user; missing identifier on account", errParams)}
	}

	if err := pager.Page(req, identifier, body.Content); err != nil {
		return typhon.Response{Error: terrors.Augment(err, "Failed to page user", errParams)}
	}

	return req.Response(&accountproto.PageAccountResponse{})
}

func getIdentifierFromAccount(account *domain.Account, pagerType string) (string, error) {
	// TODO::
	return "", nil
}
