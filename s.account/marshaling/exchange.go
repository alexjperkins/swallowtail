package marshaling

import (
	"swallowtail/libraries/util"
	"swallowtail/s.account/domain"
	accountproto "swallowtail/s.account/proto"

	"github.com/monzo/terrors"
)

// ExchangeProtoToDomain marshals the respective proto to the domain.
func ExchangeProtoToDomain(in *accountproto.Exchange) *domain.Exchange {
	return &domain.Exchange{
		Exchange:  in.Exchange.String(),
		APIKey:    in.ApiKey,
		SecretKey: in.SecretKey,
		UserID:    in.UserId,
		IsActive:  in.IsActive,
	}
}

// ExchangeDomainToProto marshals an exchange domain to the respective proto.
// All keys are masked by default.
func ExchangeDomainToProto(in *domain.Exchange) (*accountproto.Exchange, error) {
	exchangeType, err := convertExchangeTypeToProto(in.Exchange)
	if err != nil {
		return nil, err
	}

	return &accountproto.Exchange{
		ExchangeId: in.ExchangeID,
		UserId:     in.UserID,
		ApiKey:     util.MaskKey(in.APIKey, 4),
		SecretKey:  util.MaskKey(in.SecretKey, 4),
		Exchange:   exchangeType,
		IsActive:   in.IsActive,
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
