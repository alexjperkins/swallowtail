package handler

import (
	"context"
	"swallowtail/s.account/dao"
	"swallowtail/s.account/marshaling"
	accountproto "swallowtail/s.account/proto"

	"github.com/monzo/terrors"
)

// ListAccounts returns a list of all given accounts.
func (s *AccountService) ListAccounts(
	ctx context.Context, in *accountproto.ListAccountsRequest,
) (*accountproto.ListAccountsResponse, error) {
	accounts, err := dao.ListAccounts(ctx)
	if err != nil {
		return nil, terrors.Augment(err, "Failed to list accounts", nil)

	}

	var protoAccounts []*accountproto.Account
	for _, account := range accounts {
		protoAccounts = append(protoAccounts, marshaling.AccountDomainToProto(account))
	}

	return &accountproto.ListAccountsResponse{
		Accounts: protoAccounts,
	}, nil
}
