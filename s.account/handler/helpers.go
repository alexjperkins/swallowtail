package handler

import (
	"swallowtail/libraries/gerrors"
	"swallowtail/s.account/domain"
	accountproto "swallowtail/s.account/proto"
)

func getIdentifierFromAccount(account *domain.Account, pagerType string) (string, error) {
	errParams := map[string]string{
		"user_id": account.UserID,
	}

	switch pagerType {
	case accountproto.PagerType_DISCORD.String():
		return account.UserID, nil

	case accountproto.PagerType_EMAIL.String():
		if account.Email == "" {
			return "", gerrors.FailedPrecondition("failed_to_get_identifier_from_account.email", errParams)
		}
		return account.Email, nil
	case accountproto.PagerType_PHONE.String():
		if account.PhoneNumber == "" {
			return "", gerrors.FailedPrecondition("failed_to_get_identifier_from_account.phone_number", errParams)
		}
		return account.PhoneNumber, nil
	case accountproto.PagerType_SMS.String():
		if account.PhoneNumber == "" {
			return "", gerrors.FailedPrecondition("failed_to_get_identifier_from_account.sms", errParams)
		}
		return account.PhoneNumber, nil
	}

	errParams["pager_type"] = pagerType
	return "", gerrors.FailedPrecondition("failed_to_get_identifier_from_account.unknown_pager_type", errParams)
}

func isValidActorID(actorID string) bool {
	switch actorID {
	case accountproto.ActorSystemPayments, accountproto.ActorManual:
		return true
	default:
		return false
	}
}

func isValidActorUnmaskedRequest(actorID string, isRequestingUnmaskedCredentials bool) bool {
	if !isRequestingUnmaskedCredentials {
		return true
	}

	if actorID != accountproto.ActorSystemTradeEngine {
		return false
	}

	return true
}
