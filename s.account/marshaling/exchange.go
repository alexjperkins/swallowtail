package marshaling

import (
	"swallowtail/libraries/encryption"
	"swallowtail/libraries/util"
	"swallowtail/s.account/domain"
	accountproto "swallowtail/s.account/proto"

	"github.com/monzo/terrors"
)

// ExchangeProtoToDomain marshals the respective proto to the domain.
func ExchangeProtoToDomain(in *accountproto.Exchange) (*domain.Exchange, error) {
	// TODO: we need a proper passphrase here.
	encryptedAPIKey, err := encryption.EncryptWithAES([]byte(in.ApiKey), "passphrase")
	if err != nil {
		return nil, terrors.Augment(err, "Failed to marshal proto to domain; encryption of api key failed", nil)
	}

	encryptedSecretKey, err := encryption.EncryptWithAES([]byte(in.SecretKey), "passphrase")
	if err != nil {
		return nil, terrors.Augment(err, "Failed to marshal proto to domain; encryption of api key failed", nil)
	}

	return &domain.Exchange{
		Exchange:  in.ExchangeType.String(),
		APIKey:    encryptedAPIKey,
		SecretKey: encryptedSecretKey,
		IsActive:  in.IsActive,
	}, nil
}

// ExchangeDomainToProto marshals an exchange domain to the respective proto.
// All keys are masked by default.
func ExchangeDomainToProto(in *domain.Exchange) (*accountproto.Exchange, error) {
	exchangeType, err := convertExchangeTypeToProto(in.Exchange)
	if err != nil {
		return nil, err
	}

	decryptedAPIKey, err := encryption.DecryptWithAES([]byte(in.APIKey), "passphrase")
	if err != nil {
		return nil, terrors.Augment(err, "Failed to marshal domain to proto; decryption of api key failed", nil)
	}

	decryptedSecretKey, err := encryption.DecryptWithAES([]byte(in.SecretKey), "passphrase")
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

// ExchangeDomainToProtos ...
func ExchangeDomainToProtos(ins []*domain.Exchange) ([]*accountproto.Exchange, error) {
	protos := []*accountproto.Exchange{}
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
