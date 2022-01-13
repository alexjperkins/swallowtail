package marshaling

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"swallowtail/s.account/domain"
	accountproto "swallowtail/s.account/proto"
)

// AccountDomainToProto marshals an account domain object into the account proto definition.
func AccountDomainToProto(account *domain.Account) *accountproto.Account {
	return &accountproto.Account{
		UserId:             account.UserID,
		Username:           account.Username,
		Email:              account.Email,
		IsFuturesMember:    account.IsFuturesMember,
		IsAdmin:            account.IsAdmin,
		Created:            timestamppb.New(account.Created),
		LastUpdated:        timestamppb.New(account.Updated),
		PrimaryVenue:       account.PrimaryVenue,
		DefaultDcaStrategy: account.DefaultDCAStrategy,
	}
}

// UpdateAccountProtoToDomain marshals a `UpdateAccountRequest` proto message to the domain.
func UpdateAccountProtoToDomain(in *accountproto.UpdateAccountRequest) *domain.Account {
	return &domain.Account{
		UserID:             in.UserId,
		Username:           in.Username,
		Email:              in.Email,
		PhoneNumber:        in.PhoneNumber,
		HighPriorityPager:  in.HighPriorityPager.String(),
		LowPriorityPager:   in.LowPriorityPager.String(),
		IsFuturesMember:    in.IsFutures,
		IsAdmin:            in.IsAdmin,
		DefaultDCAStrategy: in.DefaultDcaStrategy,
	}
}
