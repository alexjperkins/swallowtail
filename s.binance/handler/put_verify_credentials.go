package handler

import (
	"context"
	"swallowtail/libraries/gerrors"
	accountproto "swallowtail/s.account/proto"
	"swallowtail/s.binance/client"
	"swallowtail/s.binance/marshaling"
	binanceproto "swallowtail/s.binance/proto"

	"github.com/monzo/slog"
)

// VerifyCredentials ...
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

	errParams := map[string]string{
		"user_id": in.UserId,
	}

	_, err := (&accountproto.ReadAccountRequest{
		UserId: in.UserId,
	}).Send(ctx).Response()
	switch {
	case gerrors.Is(err, gerrors.ErrNotFound, "failed_to_read_account.account_not_exist"):
	case err != nil:
		return nil, gerrors.Augment(err, "failed_to_verify_credentials.failed_to_read_accout", errParams)
	}

	dtoCredentials := marshaling.CredentialsProtoToDTO(in.Credentials)

	rsp, err := client.VerifyCredentials(ctx, dtoCredentials)
	if err != nil {
		slog.Error(ctx, "%+v: %v", rsp, err)
		return nil, gerrors.Augment(err, "failed_to_verify_credentials", nil)
	}

	proto := marshaling.VerifyRequestDTOToProto(rsp)

	slog.Info(ctx, "%s: verified credentials %v", in.UserId, proto.Success)

	return proto, nil
}
