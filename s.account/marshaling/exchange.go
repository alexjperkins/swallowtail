package marshaling

import (
	"github.com/monzo/terrors"

	"swallowtail/libraries/encryption"
	"swallowtail/libraries/gerrors"
	"swallowtail/libraries/util"
	"swallowtail/s.account/domain"
	accountproto "swallowtail/s.account/proto"
)

// ExchangeProtoToDomain marshals the respective proto to the domain.
func ExchangeProtoToDomain(userID string, exchange *accountproto.Exchange) (*domain.Exchange, error) {
	// TODO: we need a proper passphrase here.
	encryptedAPIKey, err := encryption.EncryptWithAES([]byte(exchange.ApiKey), "passphrase")
	if err != nil {
		return nil, gerrors.Augment(err, "failed-to-marshal-proto-to-domain.bad-api-key", nil)
	}

	encryptedSecretKey, err := encryption.EncryptWithAES([]byte(exchange.SecretKey), "passphrase")
	if err != nil {
		return nil, gerrors.Augment(err, "failed-to-marshal-proto-to-domain.bad-secret-key", nil)
	}

	return &domain.Exchange{
		ExchangeType: exchange.ExchangeType.String(),
		APIKey:       encryptedAPIKey,
		SecretKey:    encryptedSecretKey,
		IsActive:     exchange.IsActive,
		UserID:       userID,
	}, nil
}

// ExchangeDomainToProtos ...
func ExchangeDomainsToProtos(ins []*domain.Exchange) ([]*accountproto.Exchange, error) {
	protos := make([]*accountproto.Exchange, 0, len(ins))
	for _, in := range ins {
		proto, err := ExchangeDomainToProto(in)
		if err != nil {
			// TODO; better handling of this. Multi-error/
			return nil, err
		}
		protos = append(protos, proto)
	}
	return protos, nil
}

// ExchangeDomainToProto marshals an exchange domain to the respective proto.
// All keys are masked by default.
func ExchangeDomainToProto(in *domain.Exchange) (*accountproto.Exchange, error) {
	exchangeType, err := convertExchangeTypeToProto(in.ExchangeType)
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

	return &accountproto.Exchange{
		ExchangeId:   in.ExchangeID,
		ApiKey:       util.MaskKey(decryptedAPIKey, 4),
		SecretKey:    util.MaskKey(decryptedSecretKey, 4),
		ExchangeType: exchangeType,
		IsActive:     in.IsActive,
	}, nil
}

// ExchangeDomainToProtosUnmasked ...
func ExchangeDomainsToProtosUnmasked(ins []*domain.Exchange) ([]*accountproto.Exchange, error) {
	protos := make([]*accountproto.Exchange, 0, len(ins))

	for _, in := range ins {
		proto, err := ExchangeDomainToProtoUnmasked(in)
		if err != nil {
			return nil, err
		}

		protos = append(protos, proto)
	}

	return protos, nil
}

// ExchangeDomainToProtoUnmasked ...
// NOTE: only use this on internal endpoints; we cannot allow keys to be leaked.
func ExchangeDomainToProtoUnmasked(in *domain.Exchange) (*accountproto.Exchange, error) {
	exchangeType, err := convertExchangeTypeToProto(in.ExchangeType)
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

	return &accountproto.Exchange{
		ExchangeId:   in.ExchangeID,
		ApiKey:       decryptedAPIKey,
		SecretKey:    decryptedSecretKey,
		ExchangeType: exchangeType,
		IsActive:     in.IsActive,
	}, nil
}

func convertExchangeTypeToProto(t string) (accountproto.ExchangeType, error) {
	switch t {
	case accountproto.ExchangeType_BINANCE.String():
		return accountproto.ExchangeType_BINANCE, nil
	case accountproto.ExchangeType_FTX.String():
		return accountproto.ExchangeType_FTX, nil
	default:
		return 0, terrors.PreconditionFailed("unsupported-exchange-type", "Bad exchange type", map[string]string{
			"exchange_type": t,
		})
	}
}
