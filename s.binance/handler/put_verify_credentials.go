package handler

import (
	"context"
	"swallowtail/libraries/gerrors"
	"swallowtail/s.binance/client"
	"swallowtail/s.binance/marshaling"
	binanceproto "swallowtail/s.binance/proto"

	"github.com/monzo/slog"
)

func (s *BinanceService) VerifyCredentials(
	ctx context.Context, in *binanceproto.VerifyCredentialsRequest,
) (*binanceproto.VerifyCredentialsResponse, error) {
	switch {
	case in.UserId == "":
		return nil, gerrors.BadParam("missing_param.user_id", nil)
	case in.GetCredentials() == nil:
		return nil, gerrors.BadParam("missing_param.credentials", nil)
	case in.Credentials.ApiKey == "":
		return nil, gerrors.BadParam("missing_param.credentials.api_key", nil)
	case in.Credentials.SecretKey == "":
		return nil, gerrors.BadParam("missing_param.credentials.secret_key", nil)
	}

	// TODO; check the user actually exists

	dtoCredentials := marshaling.CredentialsProtoToDTO(in.Credentials)

	rsp, err := client.VerifyCredentials(ctx, dtoCredentials)
	if err != nil {
		slog.Error(ctx, "%+v: %v", rsp, err)
		return nil, gerrors.Augment(err, "failed_to_verify_credentials", nil)
	}

	slog.Warn(ctx, "%+v: %v", rsp, err)

	proto := marshaling.VerifyRequestDTOToProto(rsp)

	slog.Info(ctx, "%s: verified credentials %v", in.UserId, proto.Success)

	return proto, nil
}
