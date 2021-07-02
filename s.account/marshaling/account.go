package marshaling

import (
	"swallowtail/s.account/domain"
	accountproto "swallowtail/s.account/proto"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// AccountDomainToProto marshals an account domain object into the account proto definition.
func AccountDomainToProto(account *domain.Account) *accountproto.Account {
	return &accountproto.Account{
		UserId:          account.UserID,
		Email:           account.Email,
		IsFuturesMember: account.IsFuturesMember,
		IsAdmin:         account.IsAdmin,
		Created:         timestamppb.New(account.Created),
		LastUpdated:     timestamppb.New(account.Updated),
	}
}
