package marshaling

import (
	accountproto "swallowtail/s.account/proto"
	tradeengineproto "swallowtail/s.trade-engine/proto"
)

// VenueAccountToVenueCredentials ...
func VenueAccountToVenueCredentials(venueAccount *accountproto.VenueAccount) *tradeengineproto.VenueCredentials {
	return &tradeengineproto.VenueCredentials{
		ApiKey:     venueAccount.ApiKey,
		SecretKey:  venueAccount.SecretKey,
		Subaccount: venueAccount.SubAccount,
		Venue:      venueAccount.Venue,
	}
}
