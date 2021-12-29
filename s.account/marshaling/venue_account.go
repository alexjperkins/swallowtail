package marshaling

import (
	"github.com/monzo/terrors"

	"swallowtail/libraries/encryption"
	"swallowtail/libraries/gerrors"
	"swallowtail/libraries/util"
	"swallowtail/s.account/domain"
	accountproto "swallowtail/s.account/proto"
	tradeengineproto "swallowtail/s.trade-engine/proto"
)

// VenueAccountProtoToDomain marshals the respective proto to the domain.
func VenueAccountProtoToDomain(userID string, venueAccount *accountproto.VenueAccount) (*domain.VenueAccount, error) {
	// TODO: we need a proper passphrase here.
	encryptedAPIKey, err := encryption.EncryptWithAES([]byte(venueAccount.ApiKey), "passphrase")
	if err != nil {
		return nil, gerrors.Augment(err, "failed-to-marshal-proto-to-domain.bad-api-key", nil)
	}

	encryptedSecretKey, err := encryption.EncryptWithAES([]byte(venueAccount.SecretKey), "passphrase")
	if err != nil {
		return nil, gerrors.Augment(err, "failed-to-marshal-proto-to-domain.bad-secret-key", nil)
	}

	return &domain.VenueAccount{
		VenueID:   venueAccount.Venue.String(),
		APIKey:    encryptedAPIKey,
		SecretKey: encryptedSecretKey,
		IsActive:  venueAccount.IsActive,
		UserID:    userID,
	}, nil
}

// VenueAccountDomainToProtos ...
func VenueAccountDomainsToProtos(ins []*domain.VenueAccount) ([]*accountproto.VenueAccount, error) {
	protos := make([]*accountproto.VenueAccount, 0, len(ins))
	for _, in := range ins {
		proto, err := VenueAccountDomainToProto(in)
		if err != nil {
			// TODO; better handling of this. Multi-error/
			return nil, err
		}
		protos = append(protos, proto)
	}
	return protos, nil
}

// VenueAccountDomainToProto marshals an exchange domain to the respective proto.
// All keys are masked by default.
func VenueAccountDomainToProto(in *domain.VenueAccount) (*accountproto.VenueAccount, error) {
	venue, err := convertVenueIDToProto(in.VenueID)
	if err != nil {
		return nil, err
	}

	decryptedAPIKey, err := encryption.DecryptWithAES(in.APIKey, "passphrase")
	if err != nil {
		return nil, terrors.Augment(err, "Failed to marshal domain to proto; decryption of api key failed", nil)
	}

	decryptedSecretKey, err := encryption.DecryptWithAES(in.SecretKey, "passphrase")
	if err != nil {
		return nil, terrors.Augment(err, "Failed to marshal domain to proto; decryption of api key failed", nil)
	}

	return &accountproto.VenueAccount{
		VenueAccountId: in.VenueAccountID,
		ApiKey:         util.MaskKey(decryptedAPIKey, 4),
		SecretKey:      util.MaskKey(decryptedSecretKey, 4),
		Venue:          venue,
		IsActive:       in.IsActive,
	}, nil
}

// VenueAccountDomainToProtosUnmasked ...
func VenueAccountDomainsToProtosUnmasked(ins []*domain.VenueAccount) ([]*accountproto.VenueAccount, error) {
	protos := make([]*accountproto.VenueAccount, 0, len(ins))

	for _, in := range ins {
		proto, err := VenueAccountDomainToProtoUnmasked(in)
		if err != nil {
			return nil, err
		}

		protos = append(protos, proto)
	}

	return protos, nil
}

// VenueAccountDomainToProtoUnmasked ...
// NOTE: only use this on internal endpoints; we cannot allow keys to be leaked.
func VenueAccountDomainToProtoUnmasked(in *domain.VenueAccount) (*accountproto.VenueAccount, error) {
	venue, err := convertVenueIDToProto(in.VenueID)
	if err != nil {
		return nil, err
	}

	decryptedAPIKey, err := encryption.DecryptWithAES(in.APIKey, "passphrase")
	if err != nil {
		return nil, terrors.Augment(err, "Failed to marshal domain to proto; decryption of api key failed", nil)
	}

	decryptedSecretKey, err := encryption.DecryptWithAES(in.SecretKey, "passphrase")
	if err != nil {
		return nil, terrors.Augment(err, "Failed to marshal domain to proto; decryption of api key failed", nil)
	}

	return &accountproto.VenueAccount{
		VenueAccountId: in.VenueAccountID,
		ApiKey:         decryptedAPIKey,
		SecretKey:      decryptedSecretKey,
		Venue:          venue,
		IsActive:       in.IsActive,
	}, nil
}

func convertVenueIDToProto(venueID string) (tradeengineproto.VENUE, error) {
	switch venueID {
	case tradeengineproto.VENUE_BINANCE.String():
		return tradeengineproto.VENUE_BINANCE, nil
	case tradeengineproto.VENUE_BITFINEX.String():
		return tradeengineproto.VENUE_BITFINEX, nil
	case tradeengineproto.VENUE_DERIBIT.String():
		return tradeengineproto.VENUE_DERIBIT, nil
	case tradeengineproto.VENUE_FTX.String():
		return tradeengineproto.VENUE_FTX, nil
	default:
		return 0, gerrors.Unimplemented("unsupported_venue", map[string]string{
			"venue_id": venueID,
		})
	}
}
