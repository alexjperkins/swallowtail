package marshaling

import (
	"swallowtail/libraries/gerrors"
	accountproto "swallowtail/s.account/proto"
	tradeengineproto "swallowtail/s.trade-engine/proto"
)

// AccountExchangeToVenueCredentials ...
func AccountExchangeToVenueCredentials(exchange *accountproto.Exchange) (*tradeengineproto.VenueCredentials, error) {
	var venue tradeengineproto.VENUE
	switch exchange.ExchangeType {
	case accountproto.ExchangeType_BINANCE:
		venue = tradeengineproto.VENUE_BINANCE
	case accountproto.ExchangeType_DERIBIT:
		venue = tradeengineproto.VENUE_DERIBIT
	case accountproto.ExchangeType_FTX:
		venue = tradeengineproto.VENUE_FTX
	case accountproto.ExchangeType_BITFINEX:
		venue = tradeengineproto.VENUE_BITFINEX
	default:
		return nil, gerrors.FailedPrecondition("failed_to_marshal_account_exchange_to_venue_credentials.translation.exchange", map[string]string{
			"exchange_type": exchange.ExchangeType.String(),
		})
	}

	return &tradeengineproto.VenueCredentials{
		ApiKey:     exchange.ExchangeId,
		SecretKey:  exchange.SecretKey,
		Subaccount: exchange.SubAccount,
		Venue:      venue,
	}, nil
}
