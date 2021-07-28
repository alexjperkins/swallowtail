package marshaling

import (
	"swallowtail/s.bybt/client"
	bybtproto "swallowtail/s.bybt/proto"
)

// GetExchangeFundingRatesProtoToDTO ...
func GetExchangeFundingRatesProtoToDTO(in *client.GetExchangeFundingRatesByAssetResponse) *bybtproto.GetExchangeFundingRatesResponse {
	return &bybtproto.GetExchangeFundingRatesResponse{}
}
