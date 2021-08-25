package marshaling

import (
	"strings"
	"swallowtail/s.binance/client"
	binanceproto "swallowtail/s.binance/proto"
)

// CredentialsProtoToDTO ...
func CredentialsProtoToDTO(in *binanceproto.Credentials) *client.Credentials {
	return &client.Credentials{
		APIKey:    in.ApiKey,
		SecretKey: in.SecretKey,
	}
}

func VerifyRequestDTOToProto(in *client.VerifyCredentialsResponse) *binanceproto.VerifyCredentialsResponse {
	isSuccess, reason := isSuccess(in)

	return &binanceproto.VerifyCredentialsResponse{
		Success:         isSuccess,
		ReadEnabled:     in.EnableReading,
		FuturesEnabled:  in.EnableFutures,
		WithdrawEnabled: in.EnableWithdrawals,
		SpotEnabled:     in.EnableSpotAndMarginTrading,
		OptionsEnabled:  in.EnableVanillaOptions,
		IpRestrictions:  in.IPRestrict,
		Reason:          reason,
	}
}

func isSuccess(rsp *client.VerifyCredentialsResponse) (bool, string) {
	reasons := []string{}
	if !rsp.EnableReading {
		reasons = append(reasons, "please enabled ability to read account")
	}
	if !rsp.EnableFutures {
		reasons = append(reasons, "please enable futures access")
	}
	if rsp.EnableWithdrawals {
		reasons = append(reasons, "withdrawals enabled, please turn off")
	}
	if !rsp.IPRestrict {
		reasons = append(reasons, "no ip restrictions, please consider adding")
	}
	if !rsp.EnableSpotAndMarginTrading {
		reasons = append(reasons, "please enable spot access")
	}

	return rsp.EnableReading && rsp.EnableFutures && rsp.EnableSpotAndMarginTrading, strings.Join(reasons, ",")
}
