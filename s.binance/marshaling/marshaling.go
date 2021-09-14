package marshaling

import (
	"strconv"
	"strings"
	"swallowtail/libraries/gerrors"
	"swallowtail/s.binance/client"
	binanceproto "swallowtail/s.binance/proto"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// CredentialsProtoToDTO ...
func CredentialsProtoToDTO(in *binanceproto.Credentials) *client.Credentials {
	return &client.Credentials{
		APIKey:    in.ApiKey,
		SecretKey: in.SecretKey,
	}
}

// VerifyRequestDTOToProto ...
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

// PerpetualFuturesAccountBalanceDTOToProto ...
func PerpetualFuturesAccountBalanceDTOToProto(in *client.PerpetualFuturesAccountBalance) (*binanceproto.ReadPerpetualFuturesAccountResponse, error) {
	balance, err := strconv.ParseFloat(in.Balance, 64)
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_parse_float.balance", nil)
	}

	availableBalance, err := strconv.ParseFloat(in.AvailableBalance, 64)
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_parse_float.available_balance", nil)
	}

	return &binanceproto.ReadPerpetualFuturesAccountResponse{
		Asset:            in.Asset,
		Balance:          float32(balance),
		AvailableBalance: float32(availableBalance),
		LastUpdated:      timestamppb.New(time.Unix(int64(in.LastUpdated/1_000), 0)),
	}, nil
}

func isSuccess(rsp *client.VerifyCredentialsResponse) (bool, string) {
	reasons := []string{}

	if !rsp.EnableReading {
		reasons = append(reasons, "Please enable the ability to read account")
	}

	if !rsp.EnableFutures {
		reasons = append(reasons, "Please enable futures access")
	}

	if rsp.EnableWithdrawals {
		reasons = append(reasons, "You have withdrawals enabled, please turn them off")
	}

	if !rsp.IPRestrict {
		reasons = append(reasons, "You have no ip restrictions; please consider adding them")
	}

	if !rsp.EnableSpotAndMarginTrading {
		reasons = append(reasons, "Please enable spot access")
	}

	return rsp.EnableReading && rsp.EnableFutures && rsp.EnableSpotAndMarginTrading, strings.Join(reasons, ",")
}
