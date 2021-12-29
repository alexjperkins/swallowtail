package handler

import (
	"swallowtail/libraries/gerrors"
	"swallowtail/s.account/domain"
	accountproto "swallowtail/s.account/proto"
	tradeengineproto "swallowtail/s.trade-engine/proto"
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
	case accountproto.ActorSystemPayments, accountproto.ActorManual, accountproto.ActorSystemTradeEngine:
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

func validateVenueAccount(venueAccount *accountproto.VenueAccount) error {
	if venueAccount == nil {
		return gerrors.BadParam("missing_param.venue_account", nil)
	}

	switch {
	case venueAccount.ApiKey == "":
		return gerrors.BadParam("missing_param.api_key", nil)
	case venueAccount.SecretKey == "":
		return gerrors.BadParam("missing_param.secret_key", nil)
	}

	switch venueAccount.Venue {
	case tradeengineproto.VENUE_FTX:
		if venueAccount.SubAccount == "" {
			return gerrors.FailedPrecondition("subaccount_required_for_ftx", nil)
		}
	case tradeengineproto.VENUE_UNREQUIRED:
		return gerrors.BadParam("missing_param.venue_account.venue", nil)
	}

	return nil
}
