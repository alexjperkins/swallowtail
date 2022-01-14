package handler

import (
	"context"

	"swallowtail/libraries/gerrors"
	"swallowtail/s.ftx/client"
	"swallowtail/s.ftx/marshaling"
	ftxproto "swallowtail/s.ftx/proto"
)

// ListAccountBalances ...
func (s *FTXService) ListAccountBalances(
	ctx context.Context, in *ftxproto.ListAccountBalancesRequest,
) (*ftxproto.ListAccountBalancesResponse, error) {
	// Basic validation.
	switch {
	case in.Credentials == nil:
		return nil, gerrors.BadParam("missing_param.credentials", nil)
	}

	// Validate credential.
	if err := validateCredentials(in.Credentials); err != nil {
		return nil, gerrors.Augment(err, "failed_to_list_account_balances", nil)
	}

	errParams := map[string]string{
		"subaccount": in.GetCredentials().Subaccount,
	}

	// Marshal credentials to DTO.
	domainCredentials := marshaling.VenueCredentialsProtoToFTXCredentials(in.GetCredentials())

	// List account balances.
	rsp, err := client.ListAccountBalances(ctx, domainCredentials)
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_list_account_balances", errParams)
	}

	// Marshal account balance to protos.
	protoAccountBalances := marshaling.AccountBalancesDTOToProtos(rsp.AccountBalances)

	return &ftxproto.ListAccountBalancesResponse{
		AccountBalances: protoAccountBalances,
	}, nil
}
