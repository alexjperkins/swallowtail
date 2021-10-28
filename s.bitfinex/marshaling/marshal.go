package marshaling

import (
	"swallowtail/s.bitfinex/dto"
	bitfinexproto "swallowtail/s.bitfinex/proto"
)

func convertOperative(in int) bool {
	switch in {
	case 0:
		// 0 is maintainance.
		return false
	default:
		// 1 is operational.
		return true
	}
}

// GetStatusDTOToProto ...
func GetStatusDTOToProto(in *dto.GetStatusResponse) *bitfinexproto.GetBitfinexStatusResponse {
	return &bitfinexproto.GetBitfinexStatusResponse{
		Operational:     convertOperative(in.Operative),
		ServerLatencyMs: int64(in.ServerLatency),
	}
}
